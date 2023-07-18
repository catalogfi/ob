package store

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/susruth/wbtc-garden/blockchain"
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/price"
	"github.com/susruth/wbtc-garden/rest"
	"gorm.io/gorm"
)

type Store interface {
	rest.Store
	// update an order, used to update the status
	UpdateOrder(order *model.Order) error
	// get all active orders
	GetActiveOrders() ([]model.Order, error)
}

type store struct {
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
	return &store{db: db}, nil
}
func (s *store) GetValueLocked(user string, chain model.Chain) (int64, error) {
	var initAmounts, followAmounts []model.LockedAmount
	if err := s.db.Table("atomic_swaps").
		Select("asset as asset,SUM(amount) as amount").
		Joins("JOIN orders ON orders.initiator_atomic_swap_id = atomic_swaps.id").
		Where("orders.maker = ? AND (orders.status = ? OR orders.status = ?) AND atomic_swaps.chain = ?", user, model.InitiatorAtomicSwapInitiated, model.FollowerAtomicSwapInitiated, chain).
		Group("asset").
		Find(&initAmounts).Error; err != nil {
		return 0, err
	}
	if err := s.db.Table("atomic_swaps").
		Select("asset as asset,SUM(amount) as amount").
		Joins("JOIN orders ON orders.follower_atomic_swap_id = atomic_swaps.id").
		Where("orders.taker = ? AND (orders.status = ? OR orders.status = ?) AND atomic_swaps.chain = ?", user, model.InitiatorAtomicSwapRedeemed, model.FollowerAtomicSwapInitiated, chain).
		Group("asset").
		Find(&followAmounts).Error; err != nil {
		return 0, err
	}
	combinedArray := model.CombineAndAddAmount(initAmounts, followAmounts)

	var sum float64 = 0
	for _, tokenAmount := range combinedArray {
		if tokenAmount.Amount.Valid {
			priceInUSD, err := price.GetPriceInUSD(tokenAmount.Asset, chain)
			if err != nil {
				return 0, err
			}
			sum += price.GetPrice(tokenAmount.Asset, chain, float64(tokenAmount.Amount.Int64), priceInUSD)
		}
	}
	//flooring the sum
	return int64(sum), nil
}

func (s *store) CreateOrder(creator, sendAddress, recieveAddress, orderPair, sendAmount, recieveAmount, secretHash string, urls map[model.Chain]string) (uint, error) {
	sendChain, recieveChain, sendAsset, recieveAsset, err := model.ParseOrderPair(orderPair)
	if err != nil {
		return 0, err
	}

	initiatorAtomicSwap := model.AtomicSwap{
		InitiatorAddress: sendAddress,
		Chain:            sendChain,
		Asset:            sendAsset,
		Amount:           sendAmount,
	}

	followerAtomicSwap := model.AtomicSwap{
		RedeemerAddress: recieveAddress,
		Chain:           recieveChain,
		Asset:           recieveAsset,
		Amount:          recieveAmount,
	}

	fmt.Println(creator, "address")

	orders, err := s.FilterOrders(creator, "", "", "", "", 0, 0, 0, 0, 0, false)
	if err != nil {
		return 0, err
	}

	sendAmt, ok := new(big.Int).SetString(sendAmount, 10)
	if !ok {
		return 0, fmt.Errorf("invalid send amount: %s", sendAmount)
	}

	recieveAmt, ok := new(big.Int).SetString(recieveAmount, 10)
	if !ok {
		return 0, fmt.Errorf("invalid recieve amount: %s", recieveAmount)
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
	price, _ := new(big.Float).Quo(new(big.Float).SetInt(sendAmt), new(big.Float).SetInt(recieveAmt)).Float64()

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

func (s *store) FillOrder(orderID uint, filler, sendAddress, recieveAddress string, urls map[model.Chain]string) error {
	order := &model.Order{}
	if tx := s.db.First(order, orderID); tx.Error != nil {
		return tx.Error
	}
	if order.Status != model.OrderCreated {
		return fmt.Errorf("order already filled, current status: %v", order.Status)
	}

	fromChain, toChain, _, _, err := model.ParseOrderPair(order.OrderPair)
	if err != nil {
		panic(fmt.Errorf("constraint violation: invalid order pair: %v", err))
	}
	initiateAtomicSwapTimelock, err := blockchain.CalculateExpiry(fromChain, true, urls)
	if err != nil {
		panic(fmt.Errorf("constraint violation: invalid order pair: %v", err))
	}
	followerAtomicSwapTimelock, err := blockchain.CalculateExpiry(toChain, false, urls)
	if err != nil {
		panic(fmt.Errorf("constraint violation: invalid order pair: %v", err))
	}

	initiateAtomicSwap := &model.AtomicSwap{}
	if tx := s.db.First(initiateAtomicSwap, order.InitiatorAtomicSwapID); tx.Error != nil {
		return tx.Error
	}
	followerAtomicSwap := &model.AtomicSwap{}
	if tx := s.db.First(followerAtomicSwap, order.FollowerAtomicSwapID); tx.Error != nil {
		return tx.Error
	}
	initiateAtomicSwap.RedeemerAddress = recieveAddress
	followerAtomicSwap.InitiatorAddress = sendAddress
	initiateAtomicSwap.Timelock = initiateAtomicSwapTimelock
	followerAtomicSwap.Timelock = followerAtomicSwapTimelock
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
