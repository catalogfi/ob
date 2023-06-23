package bot

import "gorm.io/gorm"

type Status uint

const (
	Filled Status = iota
	InitiatorInitiated
	FollowerInitiated
	FollowerRedeemed
	InitiatorRedeemed
)

type Order struct {
	gorm.Model

	SecretHash string
	Status     Status
}

type Store interface {
	PutSecret(secretHash, secret string) error
	Secret(secretHash string) (string, error)
	PutStatus(secretHash string, status Status) error
	Status(secretHash string) (Status, error)
}

type store struct {
	db *gorm.DB
}

func NewStore(dialector gorm.Dialector, opts ...gorm.Option) (Store, error) {
	db, err := gorm.Open(dialector, opts...)
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Order{}); err != nil {
		return nil, err
	}
	return &store{db: db}, nil
}

func (s *store) PutSecret(secretHash, secret string) error {
	order := Order{
		SecretHash: secretHash,
		Status:     0,
	}
	if tx := s.db.Create(&order); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s *store) Secret(secretHash string) (string, error) {
	var order Order
	if tx := s.db.Where("secret_hash = ?", secretHash).First(&order); tx.Error != nil {
		return "", tx.Error
	}
	return order.SecretHash, nil
}

func (s *store) PutStatus(secretHash string, status Status) error {
	var order Order
	if tx := s.db.Where("secret_hash = ?", secretHash).First(&order); tx.Error != nil {
		return tx.Error
	}
	order.Status = status
	if tx := s.db.Save(&order); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s *store) Status(secretHash string) (Status, error) {
	var order Order
	if tx := s.db.Where("secret_hash = ?", secretHash).First(&order); tx.Error != nil {
		return 0, tx.Error
	}
	return order.Status, nil
}
