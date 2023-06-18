package store

import (
	"fmt"

	"github.com/susruth/wbtc-garden/executor"
	"github.com/susruth/wbtc-garden/model"
	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

type Store interface {
	SubStore(chain string) executor.Store
}

func New(dialector gorm.Dialector, opts ...gorm.Option) (Store, error) {
	db, err := gorm.Open(dialector, opts...)
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&model.Transaction{}); err != nil {
		return nil, err
	}
	return &store{db: db}, nil
}

type subStore struct {
	db    *gorm.DB
	chain string
}

func (s *store) SubStore(chain string) executor.Store {
	return &subStore{db: s.db, chain: chain}
}

func (s *subStore) Transactions(address string) ([]model.Transaction, error) {
	txs := []model.Transaction{}
	if res := s.db.Find(&txs, "from_address = ? OR to_address = ? AND chain = ?", address, address, s.chain); res.Error != nil {
		return nil, fmt.Errorf("no such orders for the given address (%s): %v", address, res.Error)
	}
	return txs, nil
}

func (s *subStore) PendingTransactions() ([]model.Transaction, error) {
	txs := []model.Transaction{}
	if res := s.db.Find(&txs, "status < 5 AND chain = ?", s.chain); res.Error != nil {
		return nil, fmt.Errorf("no such orders for the given address: %v", res.Error)
	}
	return txs, nil
}

func (s *subStore) PutTransaction(tx model.Transaction) error {
	tx.Chain = s.chain
	return s.db.Create(&tx).Error
}

func (s *subStore) UpdateTransaction(tx model.Transaction) error {
	tx.Chain = s.chain
	return s.db.Save(&tx).Error
}
