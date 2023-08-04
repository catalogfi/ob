package store

import (
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/susruth/wbtc-garden/blockchain"
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/price"
	"github.com/susruth/wbtc-garden/rest"
	"github.com/susruth/wbtc-garden/watcher"
	"gorm.io/gorm"
)

type Store interface {
	rest.Store
	watcher.Store
	price.Store
}

type store struct {
	mu    *sync.RWMutex
	cache map[string]float64

	db *gorm.DB
}

func New(dialector gorm.Dialector, opts ...gorm.Option) (Store, error) {
	db, err := gorm.Open(dialector, opts...)
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&model.Order{}, &model.AtomicSwap{}); err != nil {
		return nil, err
	}
	return &store{mu: new(sync.RWMutex), cache: make(map[string]float64), db: db}, nil
}

func (s *store) GetValueLocked(user string, chain model.Chain) (*big.Int, error) {
	var initAmounts, followAmounts []model.LockedAmount
	if err := s.db.Table("atomic_swaps").
		Select("asset as asset,SUM(amount::bigint) as amount").
		Joins("JOIN orders ON orders.initiator_atomic_swap_id = atomic_swaps.id").
		Where("orders.maker = ? AND (orders.status = ? OR orders.status = ? OR orders.status = ?) AND atomic_swaps.chain = ?", user, model.InitiatorAtomicSwapInitiated, model.FollowerAtomicSwapInitiated, model.FollowerAtomicSwapRefunded, chain).
		Group("asset").
		Find(&initAmounts).Error; err != nil {
		return big.NewInt(0), err
	}
	if err := s.db.Table("atomic_swaps").
		Select("asset as asset,SUM(amount::bigint) as amount").
		Joins("JOIN orders ON orders.follower_atomic_swap_id = atomic_swaps.id").
		Where("orders.taker = ? AND (orders.status = ? OR orders.status = ? OR orders.status = ?) AND atomic_swaps.chain = ?", user, model.InitiatorAtomicSwapRedeemed, model.FollowerAtomicSwapInitiated, model.InitiatorAtomicSwapRefunded, chain).
		Group("asset").
		Find(&followAmounts).Error; err != nil {
		return big.NewInt(0), err
	}
	combinedArray := model.CombineAndAddAmount(initAmounts, followAmounts)

	if len(combinedArray) == 0 {
		return big.NewInt(0), nil
	}
	sum := big.NewInt(0)
	for _, tokenAmount := range combinedArray {

		if tokenAmount.Amount.Valid {
			priceInUSD, err := s.Price("bitcoin", "ethereum")
			if err != nil {
				return big.NewInt(0), err
			}
			sum = new(big.Int).Add(price.GetPrice(model.Asset(tokenAmount.Asset), chain, big.NewInt(tokenAmount.Amount.Int64), big.NewInt(int64(priceInUSD))), sum)
			fmt.Println(sum)
		}
	}
	//flooring the sum
	return sum, nil
}

func (s *store) CreateOrder(creator, sendAddress, receiveAddress, orderPair, sendAmount, receiveAmount, secretHash string, userBtcWalletAddress string, urls map[model.Chain]string) (uint, error) {

	sendChain, receiveChain, sendAsset, receiveAsset, err := model.ParseOrderPair(orderPair)
	if err != nil {
		return 0, err
	}

	if err := model.VerifyHexString(secretHash); err != nil {
		return 0, err
	}

	priceByOracle, err := s.Price("bitcoin", "ethereum")
	if err != nil {
		return 0, err
	}

	initiatorLockValue, err := s.GetValueLocked(creator, sendChain)
	if err != nil {
		return 0, err
	}

	initiatorMinConfirmations := GetMinConfirmations(initiatorLockValue, sendChain)

	initiatorAtomicSwap := model.AtomicSwap{
		InitiatorAddress:     sendAddress,
		Chain:                sendChain,
		Asset:                sendAsset,
		Amount:               sendAmount,
		PriceByOracle:        priceByOracle,
		MinimumConfirmations: initiatorMinConfirmations, // TODO: add custom confirmation by users
	}

	followerAtomicSwap := model.AtomicSwap{
		RedeemerAddress: receiveAddress,
		Chain:           receiveChain,
		Asset:           receiveAsset,
		Amount:          receiveAmount,
	}

	orders, err := s.FilterOrders(creator, "", "", "", "", 0, 0, 0, 0, 0, false)
	if err != nil {
		return 0, err
	}

	sendAmt, ok := new(big.Int).SetString(sendAmount, 10)
	if !ok {
		return 0, fmt.Errorf("invalid send amount: %s", sendAmount)
	}

	receiveAmt, ok := new(big.Int).SetString(receiveAmount, 10)
	if !ok {
		return 0, fmt.Errorf("invalid receive amount: %s", receiveAmount)
	}

	// validate orderpair
	fromChain, toChain, _, _, err := model.ParseOrderPair(orderPair)
	if err != nil {
		return 0, err
	}
	if _, err := blockchain.CalculateExpiry(fromChain, true, urls); err != nil {

		return 0, err
	}
	if _, err := blockchain.CalculateExpiry(toChain, false, urls); err != nil {

		return 0, err
	}

	// ignoring accuracy
	price, _ := new(big.Float).Quo(new(big.Float).SetInt(sendAmt), new(big.Float).SetInt(receiveAmt)).Float64()

	order := model.Order{
		Maker:                 creator,
		OrderPair:             orderPair,
		InitiatorAtomicSwapID: initiatorAtomicSwap.ID,
		FollowerAtomicSwapID:  followerAtomicSwap.ID,
		InitiatorAtomicSwap:   &initiatorAtomicSwap,
		FollowerAtomicSwap:    &followerAtomicSwap,
		Price:                 price,
		SecretHash:            secretHash,
		Status:                model.OrderCreated,
		SecretNonce:           uint64(len(orders)) + 1,
		UserBtcWalletAddress:  userBtcWalletAddress,
	}

	if tx := s.db.Create(&initiatorAtomicSwap); tx.Error != nil {
		return 0, tx.Error
	}

	if tx := s.db.Create(&followerAtomicSwap); tx.Error != nil {
		return 0, tx.Error
	}
	if tx := s.db.Create(&order); tx.Error != nil {
		return 0, tx.Error
	}

	return order.ID, nil
}

func (s *store) FillOrder(orderID uint, filler, sendAddress, receiveAddress string, urls map[model.Chain]string) error {
	order := &model.Order{}
	if tx := s.db.First(order, orderID); tx.Error != nil {
		return tx.Error
	}
	if order.Status != model.OrderCreated {
		return fmt.Errorf("order already filled, current status: %v", order.Status)
	}

	fromChain, toChain, _, _, err := model.ParseOrderPair(order.OrderPair)
	if err != nil {
		return fmt.Errorf("constraint violation: invalid order pair: %v", err)
	}
	initiateAtomicSwapTimelock, err := blockchain.CalculateExpiry(fromChain, true, urls)
	if err != nil {
		return fmt.Errorf("constraint violation: invalid order pair: %v", err)
	}
	followerAtomicSwapTimelock, err := blockchain.CalculateExpiry(toChain, false, urls)
	if err != nil {
		return fmt.Errorf("constraint violation: invalid order pair: %v", err)
	}

	initiateAtomicSwap := &model.AtomicSwap{}
	if tx := s.db.First(initiateAtomicSwap, order.InitiatorAtomicSwapID); tx.Error != nil {
		return tx.Error
	}
	followerAtomicSwap := &model.AtomicSwap{}
	if tx := s.db.First(followerAtomicSwap, order.FollowerAtomicSwapID); tx.Error != nil {
		return tx.Error
	}

	priceByOracle, err := s.Price("bitcoin", "ethereum")
	if err != nil {
		return err
	}
	followerLockedValue, err := s.GetValueLocked(filler, toChain)
	if err != nil {
		return err
	}

	initiatorMinConfirmations := GetMinConfirmations(followerLockedValue, toChain)

	initiateAtomicSwap.RedeemerAddress = receiveAddress
	followerAtomicSwap.InitiatorAddress = sendAddress
	initiateAtomicSwap.Timelock = initiateAtomicSwapTimelock
	followerAtomicSwap.Timelock = followerAtomicSwapTimelock
	followerAtomicSwap.PriceByOracle = priceByOracle
	followerAtomicSwap.MinimumConfirmations = initiatorMinConfirmations
	order.Taker = filler
	order.Status = model.OrderFilled
	if tx := s.db.Save(order); tx.Error != nil {
		return tx.Error
	}
	if tx := s.db.Save(initiateAtomicSwap); tx.Error != nil {
		return tx.Error
	}
	if tx := s.db.Save(followerAtomicSwap); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s *store) CancelOrder(creator string, orderID uint) error {
	order := &model.Order{}
	if tx := s.db.First(order, orderID); tx.Error != nil {
		return tx.Error
	}
	if order.Maker != creator {
		return fmt.Errorf("order can be cancelled only by its creator")
	}
	if order.Status != model.OrderCreated {
		return fmt.Errorf("order can be cancelled only if it is not filled, current status: %v", order.Status)
	}
	if tx := s.db.Delete(order); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s *store) FilterOrders(maker, taker, orderPair, secretHash, sort string, status model.Status, minPrice, maxPrice float64, page, perPage int, verbose bool) ([]model.Order, error) {
	orders := []model.Order{}
	tx := s.db
	if orderPair != "" {
		tx = tx.Where("order_pair = ?", orderPair)
	}
	if minPrice != 0 {
		tx = tx.Where("price >= ?", minPrice)
	}
	if maxPrice != 0 {
		tx = tx.Where("price <= ?", maxPrice)
	}
	if status != model.Unknown {
		tx = tx.Where("status = ?", status)
	}
	if maker != "" {
		tx = tx.Where("maker = ?", maker)
	}
	if taker != "" {
		tx = tx.Where("taker = ?", taker)
	}
	if secretHash != "" {
		tx = tx.Where("secret_hash = ?", secretHash)
	}

	// sort
	orderByList := strings.Split(sort, ",")
	orderByQuery := ""
	for _, orderBy := range orderByList {
		if orderBy == "" {
			continue
		}
		if orderByQuery != "" {
			orderByQuery += ", "
		}
		if orderBy[0] == '-' {
			orderByQuery += orderBy[1:] + " DESC"
		} else {
			orderByQuery += orderBy
		}
	}
	if orderByQuery != "" {
		tx = tx.Order(orderByQuery)
	}

	// pagination
	if page != 0 && perPage != 0 {
		tx = tx.Offset((page - 1) * perPage).Limit(perPage)
	}

	if tx = tx.Find(&orders); tx.Error != nil {
		return nil, tx.Error
	}

	if verbose {
		for i := range orders {
			if err := s.fillSwapDetails(&orders[i]); err != nil {
				return nil, err
			}
		}
	}
	return orders, nil
}

func (s *store) GetActiveOrders() ([]model.Order, error) {
	orders := []model.Order{}
	if tx := s.db.Where("status > ? AND status < ?", model.OrderCreated, model.OrderExecuted).Find(&orders); tx.Error != nil {
		return nil, tx.Error
	}
	for i := range orders {
		if err := s.fillSwapDetails(&orders[i]); err != nil {
			return nil, err
		}
	}
	return orders, nil
}

func (s *store) GetOrder(orderID uint) (*model.Order, error) {
	order := &model.Order{
		InitiatorAtomicSwap: &model.AtomicSwap{},
		FollowerAtomicSwap:  &model.AtomicSwap{},
	}
	if tx := s.db.First(order, orderID); tx.Error != nil {
		return nil, tx.Error
	}
	return order, s.fillSwapDetails(order)
}

func (s *store) UpdateOrder(order *model.Order) error {
	if tx := s.db.Save(order.FollowerAtomicSwap); tx.Error != nil {
		return tx.Error
	}
	if tx := s.db.Save(order.InitiatorAtomicSwap); tx.Error != nil {
		return tx.Error
	}
	if tx := s.db.Save(order); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s *store) fillSwapDetails(order *model.Order) error {
	order.FollowerAtomicSwap = &model.AtomicSwap{}
	order.InitiatorAtomicSwap = &model.AtomicSwap{}
	if tx := s.db.First(order.InitiatorAtomicSwap, order.InitiatorAtomicSwapID); tx.Error != nil {
		return tx.Error
	}
	if tx := s.db.First(order.FollowerAtomicSwap, order.FollowerAtomicSwapID); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s *store) SetPrice(fromChain string, toChain string, price float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[fromChain+toChain] = price
	return nil
}
func (s *store) Price(fromChain string, toChain string) (float64, error) {
	// s.mu.RLock()
	// defer s.mu.RUnlock()

	price, ok := s.cache[fromChain+toChain]
	if !ok {
		return 0, fmt.Errorf("price not found, please try later")
	}
	return price, nil
}

func GetMinConfirmations(value *big.Int, chain model.Chain) uint64 {
	if chain.IsBTC() {
		switch {
		case value.Cmp(big.NewInt(10000)) < 1:
			return 1

		case value.Cmp(big.NewInt(100000)) < 1:
			return 2

		case value.Cmp(big.NewInt(1000000)) < 1:
			return 4

		case value.Cmp(big.NewInt(10000000)) < 1:
			return 6

		case value.Cmp(big.NewInt(100000000)) < 1:
			return 8

		default:
			return 12
		}
	} else if chain.IsEVM() {
		switch {
		case value.Cmp(big.NewInt(10000)) < 1:
			return 6

		case value.Cmp(big.NewInt(100000)) < 1:
			return 12

		case value.Cmp(big.NewInt(1000000)) < 1:
			return 18

		case value.Cmp(big.NewInt(10000000)) < 1:
			return 24

		case value.Cmp(big.NewInt(100000000)) < 1:
			return 30

		default:
			return 100
		}
	}
	return 0
}

// func secretHashAlreadyExists(orderPair string, secretHash string) (bool, error) {
// 	fromChain, toChain, fromAsset, ToAsset, err := model.ParseOrderPair(orderPair)
// 	if err != nil {
// 		return false, err
// 	}
// 	if model.Chain(fromChain).IsEVM() {
// 		queryChainForsecretHash(string(fromAsset), secretHash)
// 	} else if model.Chain(toChain).IsEVM() {
// 		queryChainForsecretHash(string(ToAsset), secretHash)
// 	}
// 	// TODO:
// 	return false, nil
// }

// func queryChainForsecretHash(asset string, secretHash string) (string, error) {
// 	// TODO: need to do a change in smart contract to complete this function
// 	return "", nil
// }
