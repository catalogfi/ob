package swapper

import (
	"errors"
	"fmt"
	"time"
)

type InitiatorSwap interface {
	Initiate() (string, error)
	WaitForRedeem() ([]byte, string, error)
	IsRedeemed() (bool, []byte, string, error)
	Refund() (string, error)
	Expired() (bool, error)
}

type RedeemerSwap interface {
	Redeem(secret []byte) (string, error)
	IsInitiated() (bool, string, error)
	WaitForInitiate() (string, error)
}

var ErrInitiateTimeout = errors.New("initiate timeout")
var ErrRedeemTimeout = errors.New("redeem timeout")

func ExecuteAtomicSwapFirst(initiator InitiatorSwap, redeemer RedeemerSwap, secret []byte) error {
	fmt.Println(1)
	if _, err := initiator.Initiate(); err != nil {
		fmt.Println("Initiation Failed", err)
		return err
	}
	fmt.Println(2)
	time.Sleep(1 * time.Second)
	if _, err := redeemer.WaitForInitiate(); err != nil {
		if err == ErrInitiateTimeout {
			_, err = initiator.Refund()
			if err != nil {
				return err
			}
		}
		return err
	}
	time.Sleep(1 * time.Second)
	fmt.Println(3)
	if _, err := redeemer.Redeem(secret); err != nil {
		return err
	}
	fmt.Println(4)
	return nil
}

func ExecuteAtomicSwapSecond(initiator InitiatorSwap, redeemer RedeemerSwap) error {
	time.Sleep(1 * time.Second)
	if _, err := redeemer.WaitForInitiate(); err != nil {
		return err
	}
	if _, err := initiator.Initiate(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	secret, _, err := initiator.WaitForRedeem()
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	if secret != nil {
		if _, err := redeemer.Redeem(secret); err != nil {
			return err
		}
	}
	return nil
}
