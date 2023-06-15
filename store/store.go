package store

import (
	"fmt"

	"github.com/susruth/wbtc-garden-server/executor"
	"github.com/susruth/wbtc-garden-server/model"
	"github.com/susruth/wbtc-garden-server/rest"
	"gorm.io/gorm"
)

type store struct {
	db *gorm.DB
}

type Store interface {
	rest.Store
	executor.Store
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

func (s *store) Transactions(address string) ([]model.Transaction, error) {
	txs := []model.Transaction{}
	if res := s.db.Find(&txs, "address = ?", address); res.Error != nil {
		return nil, fmt.Errorf("no such orders for the given address (%s): %v", address, res.Error)
	}
	return txs, nil
}

func (s *store) PendingTransactions() ([]model.Transaction, error) {
	txs := []model.Transaction{}
	if res := s.db.Find(&txs, "status < 5"); res.Error != nil {
		return nil, fmt.Errorf("no such orders for the given address: %v", res.Error)
	}
	return txs, nil
}

func (s *store) PutTransaction(tx model.Transaction) error {
	return s.db.Create(&tx).Error
}

func (s *store) UpdateTransaction(tx model.Transaction) error {
	return s.db.Save(&tx).Error
}
