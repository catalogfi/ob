package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model

	FromAddress string `json:"fromAddress"`
	FromExpiry  int64  `json:"fromExpiry"`
	ToAddress   string `json:"toAddress"`
	ToExpiry    int64  `json:"toExpiry"`
	Secret      string `json:"secret"`
	SecretHash  string `json:"secretHash"`
	Amount      uint64 `json:"amount"`

	InitiatorInitiateTxHash string `json:"initiatorInitiateTxHash"`
	FollowerInitiateTxHash  string `json:"followerInitiateTxHash"`
	InitiatorRedeemTxHash   string `json:"initiatorRedeemTxHash"`
	FollowerRedeemTxHash    string `json:"followerRedeemTxHash"`
	FollowerRefundTxHash    string `json:"followerRefundTxHash"`

	Fee    uint64 `json:"fee"`
	Status uint8  `json:"status"`
}

type Account struct {
	BtcAddress  string  `json:"btcAddress"`
	WbtcAddress string  `json:"wbtcAddress"`
	BtcBalance  string  `json:"btcBalance"`
	WbtcBalance string  `json:"wbtcBalance"`
	Fee         float64 `json:"feeInBips"`
}
