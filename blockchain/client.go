package blockchain

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	dataHashRegistry "github.com/mbabinski218/BaaS-IoT-platform/blockchain/smartContracts"
	"github.com/mbabinski218/BaaS-IoT-platform/configs"
)

type Client struct {
	ethClient        *ethclient.Client
	dataHashRegistry *dataHashRegistry.DataHashRegistry
	auth             *bind.TransactOpts
	address          common.Address
}

func NewEthClient(url string, privateKeyHex string, contractAddress string) (*Client, error) {
	if !configs.Envs.BlockchainEnabled {
		log.Println("Blockchain client is disabled")
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	auth.GasLimit = 3000000

	var contract *dataHashRegistry.DataHashRegistry
	var contractAddr common.Address

	if contractAddress == "" {
		contractAddr, _, contract, err = dataHashRegistry.DeployDataHashRegistry(auth, client)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy contract: %w", err)
		}
		log.Println("Deployed new contract at:", contractAddr.Hex())
	} else {
		contractAddr = common.HexToAddress(contractAddress)
		contract, err = dataHashRegistry.NewDataHashRegistry(contractAddr, client)
		if err != nil {
			return nil, fmt.Errorf("failed to bind to contract: %w", err)
		}
		log.Println("Using existing contract at:", contractAddr.Hex())
	}

	log.Println("Ethereum client created successfully")
	return &Client{
		ethClient:        client,
		dataHashRegistry: contract,
		auth:             auth,
		address:          contractAddr,
	}, nil
}

func (c *Client) SetGasLimit(gasLimit uint64) {
	if !configs.Envs.BlockchainEnabled {
		return
	}

	c.auth.GasLimit = gasLimit
}

func (c *Client) Send(dataId uuid.UUID, hash [32]byte, deviceId uuid.UUID) (time.Duration, time.Duration, time.Duration, error) {
	if !configs.Envs.BlockchainEnabled {
		return 0, 0, 0, nil
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	opt := bind.TransactOpts{
		Context:  ctx,
		From:     c.auth.From,
		Signer:   c.auth.Signer,
		GasLimit: c.auth.GasLimit,
		GasPrice: c.auth.GasPrice,
	}

	sendStart := time.Now()
	transaction, err := c.dataHashRegistry.StoreHash(&opt, dataId, hash, deviceId)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("transaction send error: %w", err)
	}
	sendDuration := time.Since(sendStart)

	mineStart := time.Now()
	receipt, err := bind.WaitMined(ctx, c.ethClient, transaction)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("wait mined error: %w", err)
	}
	if receipt.Status != 1 {
		return 0, 0, 0, fmt.Errorf("transaction failed: %s", transaction.Hash().Hex())
	}
	mineDuration := time.Since(mineStart)

	log.Println("Transaction hash:", transaction.Hash().Hex())

	duration := time.Since(start)
	return duration, sendDuration, mineDuration, nil
}

func (c *Client) VerifyHash(dataId uuid.UUID, hash [32]byte) (bool, time.Duration, error) {
	if !configs.Envs.BlockchainEnabled {
		return true, 0, nil
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	exists, err := c.dataHashRegistry.VerifyHash(&bind.CallOpts{Context: ctx}, dataId, hash)
	if err != nil {
		return false, 0, fmt.Errorf("failed to verify hash: %w", err)
	}

	duration := time.Since(start)

	return exists, duration, nil
}
