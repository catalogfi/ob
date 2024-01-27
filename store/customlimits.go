package store

import "time"

type CustomLimits struct {
	Amount int64
	Expiry time.Time
}

var GrantedLimits = map[string]CustomLimits{
	"0xe103abfa0f867e53cef1ad2cb0dcbc193b385a93": {
		Amount: 300000000,
		Expiry: time.Date(2024, time.January, 28, 0, 0, 0, 0, time.UTC),
	},
	"0x3cb762058f019c3abcd5e4a07957ee996ee319bd": {
		Amount: 300000000,
		Expiry: time.Date(2024, time.January, 28, 0, 0, 0, 0, time.UTC),
	},
}
