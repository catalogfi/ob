package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model

	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	Secret      string `json:"secret"`
	SecretHash  string `json:"secretHash"`
	Amount      uint64 `json:"amount"`
	WBTCExpiry  int64  `json:"wbtcExpiry"`

	InitiatorInitiateTxHash string `json:"initiatorInitiateTxHash"`
	FollowerInitiateTxHash  string `json:"followerInitiateTxHash"`
	InitiatorRedeemTxHash   string `json:"initiatorRedeemTxHash"`
	FollowerRedeemTxHash    string `json:"followerRedeemTxHash"`
	FollowerRefundTxHash    string `json:"followerRefundTxHash"`

	Fee    uint64 `json:"fee"`
	Status uint8  `json:"status"`
}

type Account struct {
	BtcAddress       string  `json:"btcAddress"`
	WbtcAddress      string  `json:"wbtcAddress"`
	WbtcTokenAddress string  `json:"wbtcTokenAddress"`
	BtcBalance       string  `json:"btcBalance"`
	WbtcBalance      string  `json:"wbtcBalance"`
	Fee              float64 `json:"feeInBips"`
}
