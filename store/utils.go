package store

import (
	"encoding/hex"
	"errors"
)

func verifyHexString(input string) error {
	decoded, err := hex.DecodeString(input)
	if err != nil {
		return errors.New("wrong secret hash: not a valid hexadecimal string")
	}
	if len(decoded) != 32 {
		return errors.New("wrong secret hash: length should be 32 bytes (64 characters)")
	}

	return nil
}
