package ethclient

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client wraps an Ethereum client with additional functionality
type Client struct {
	*ethclient.Client
	ChainID *big.Int
}

// New creates a new Ethereum client
func New(url string) (*Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	// Default to local Ganache chain ID
	chainID := big.NewInt(1337)

	return &Client{
		Client:  client,
		ChainID: chainID,
	}, nil
}

// GetBalance returns the balance of the given address
func (c *Client) GetBalance(address string) (*big.Int, error) {
	return c.BalanceAt(context.Background(), common.HexToAddress(address), nil)
}

// TransferETH transfers ETH from one account to another
func (c *Client) TransferETH(priKey *ecdsa.PrivateKey, to common.Address, amount *big.Int) error {
	publicKey := priKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := c.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	gasLimit := uint64(21000) // in units
	gasPrice, err := c.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	var data []byte
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(c.ChainID), priKey)
	if err != nil {
		return err
	}

	return c.SendTransaction(context.Background(), signedTx)
}
