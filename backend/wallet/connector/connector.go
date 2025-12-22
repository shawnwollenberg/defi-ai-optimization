package connector

import (
	"context"
	"errors"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// WalletConnector handles wallet connections and transactions
type WalletConnector struct {
	ethClient  *ethclient.Client
	baseClient *ethclient.Client
}

// Transaction represents a built transaction
type Transaction struct {
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	GasLimit uint64 `json:"gas_limit"`
	GasPrice string `json:"gas_price"`
	ChainID  int64  `json:"chain_id"`
}

// NewWalletConnector creates a new wallet connector
func NewWalletConnector() *WalletConnector {
	ethRPC := os.Getenv("ETH_RPC_URL")
	baseRPC := os.Getenv("BASE_RPC_URL")
	
	var ethClient, baseClient *ethclient.Client
	var err error
	
	if ethRPC != "" {
		ethClient, err = ethclient.Dial(ethRPC)
		if err != nil {
			// Log error but don't fail - client will be nil
		}
	}
	
	if baseRPC != "" {
		baseClient, err = ethclient.Dial(baseRPC)
		if err != nil {
			// Log error but don't fail
		}
	}
	
	return &WalletConnector{
		ethClient:  ethClient,
		baseClient: baseClient,
	}
}

// GetClient returns the appropriate client for a chain
func (wc *WalletConnector) GetClient(chain string) (*ethclient.Client, error) {
	if chain == "base" {
		if wc.baseClient == nil {
			return nil, errors.New("Base client not initialized")
		}
		return wc.baseClient, nil
	}
	
	if wc.ethClient == nil {
		return nil, errors.New("Ethereum client not initialized")
	}
	return wc.ethClient, nil
}

// BuildTransaction builds a transaction for a given chain
func (wc *WalletConnector) BuildTransaction(chain, to, value, data string, params map[string]interface{}) (*Transaction, error) {
	client, err := wc.GetClient(chain)
	if err != nil {
		return nil, err
	}
	
	toAddress := common.HexToAddress(to)
	
	// Parse value
	valueBig := big.NewInt(0)
	if value != "" {
		valueBig, _ = valueBig.SetString(value, 10)
	}
	
	// Get chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}
	
	// Estimate gas
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: common.FromHex(data),
	})
	if err != nil {
		// Use default if estimation fails
		gasLimit = 21000
	}
	
	// Get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	
	tx := &Transaction{
		To:       to,
		Value:    valueBig.String(),
		Data:     data,
		GasLimit: gasLimit,
		GasPrice: gasPrice.String(),
		ChainID:  chainID.Int64(),
	}
	
	return tx, nil
}

// SignMessage signs a message (placeholder - actual signing happens client-side)
func (wc *WalletConnector) SignMessage(message string) (string, error) {
	// Message signing is typically done client-side with wallet
	// This is a placeholder for server-side operations if needed
	return "", errors.New("message signing must be done client-side")
}

// SendTransaction sends a signed transaction
func (wc *WalletConnector) SendTransaction(chain string, signedTx *types.Transaction) error {
	client, err := wc.GetClient(chain)
	if err != nil {
		return err
	}
	
	return client.SendTransaction(context.Background(), signedTx)
}

// GetTransactionReceipt gets a transaction receipt
func (wc *WalletConnector) GetTransactionReceipt(chain, txHash string) (*types.Receipt, error) {
	client, err := wc.GetClient(chain)
	if err != nil {
		return nil, err
	}
	
	hash := common.HexToHash(txHash)
	return client.TransactionReceipt(context.Background(), hash)
}

