package models

import (
	"crypto/ecdsa"
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
)

// User represents a user of the exchange
type User struct {
	ID         uint64
	PrivateKey *ecdsa.PrivateKey
	Address    string
}

// NewUser creates a new user with the given private key and ID
func NewUser(privKey string, userID uint64) (*User, error) {
	if privKey == "" {
		return nil, errors.New("private key cannot be empty")
	}

	pk, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to convert public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return &User{
		ID:         userID,
		PrivateKey: pk,
		Address:    address,
	}, nil
}

// GetAddress returns the Ethereum address of the user
func (u *User) GetAddress() string {
	return u.Address
}
