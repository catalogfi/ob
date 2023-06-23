package user

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spruceid/siwe-go"
	"github.com/susruth/wbtc-garden/orderbook/model"
	"github.com/susruth/wbtc-garden/orderbook/rest"
)

type Auth interface {
	rest.Auth
}

type auth struct {
}

type Claims struct {
	UserWallet string `json:"userWallet"`
	jwt.StandardClaims
}

func NewAuth() Auth {
	return &auth{}
}

func (a *auth) Verfiy(req model.VerifySiwe) (*jwt.Token, error) {
	parsedMessage, err := siwe.ParseMessage(req.Message)
	if err != nil {
		return nil, fmt.Errorf("Error parsing message: %w ", err)
	}

	valid, err := parsedMessage.ValidNow()
	if err != nil {
		return nil, fmt.Errorf("Error validating message: %w ", err)
	}

	if !valid {
		return nil, fmt.Errorf("Validating expired Token")
	}

	publicAddr, err := parsedMessage.VerifyEIP191(req.Signature)
	if err != nil {
		return nil, fmt.Errorf("Error verifying message: %w ", err)
	}
	fromAddress := crypto.PubkeyToAddress(*publicAddr)

	claims := &Claims{
		UserWallet: fromAddress.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token, nil

}
