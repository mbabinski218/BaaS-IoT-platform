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
	"github.com/mbabinski218/BaaS-IoT-platform/types"
	"github.com/mbabinski218/BaaS-IoT-platform/utils"
)

type Client struct {
	ethClient        *ethclient.Client
	dataHashRegistry *dataHashRegistry.DataHashRegistry
	auth             *bind.TransactOpts
	address          common.Address
}

func NewEthClient(url string, privateKeyHex string, contractAddress string) (*Client, error) {
	if configs.Envs.BlockchainMode == types.BCNone {
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
	if configs.Envs.BlockchainMode == types.BCNone {
		return
	}

	c.auth.GasLimit = gasLimit
}

func (c *Client) Send(dataId uuid.UUID, hash [32]byte, deviceId uuid.UUID) (time.Duration, time.Duration, time.Duration, error) {
	if configs.Envs.BlockchainMode == types.BCNone {
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
	start := time.Now()

	var success bool = true
	var err error = nil

	switch configs.Envs.BlockchainMode {
	case types.BCFullCheck:
		success, err = verifyHashFullCheck(c, dataId, hash)
	case types.BCLightCheck:
		success, err = verifyHashLightCheck(c, dataId, hash)
	}

	duration := time.Since(start)

	return success, duration, err
}

func verifyHashFullCheck(c *Client, dataId uuid.UUID, hash [32]byte) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	success, err := c.dataHashRegistry.VerifyHash(&bind.CallOpts{Context: ctx}, dataId, hash)
	if err != nil {
		return false, fmt.Errorf("failed to verify hash (full): %w", err)
	}

	return success, nil
}

func verifyHashLightCheck(c *Client, dataId uuid.UUID, hash [32]byte) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	events, err := c.dataHashRegistry.FilterHashStored(&bind.FilterOpts{Context: ctx}, [][16]byte{dataId})
	if err != nil {
		return false, fmt.Errorf("failed to get events: %w", err)
	}

	for events.Next() {
		event := events.Event
		if event.Id == dataId && event.DataHash == hash {
			return true, nil
		}
	}

	return false, fmt.Errorf("failed to verify hash (light)")
}

func (c *Client) VerifyHashes(docs []types.DocData) (bool, time.Duration, error) {
	start := time.Now()

	var success bool = true
	var err error = nil

	switch configs.Envs.BlockchainMode {
	case types.BCFullCheck:
		success, err = executeBlockchainFullCheck(c, docs)
	case types.BCLightCheck:
		success, err = executeBlockchainLightCheck(c, docs)
	case types.BCBatchCheck:
		success, err = executeBlockchainBatchCheck(c, docs)
	}

	duration := time.Since(start)
	return success, duration, err
}

func executeBlockchainFullCheck(c *Client, docs []types.DocData) (bool, error) {
	result := true

	for _, doc := range docs {
		hash, err := utils.CalculateHash(doc.Data)
		if err != nil {
			log.Printf("failed to calculate hash for doc with id: %v, hash: %x\n", doc.Id, hash)
			result = false
			continue
		}

		success, _, err := c.VerifyHash(doc.Id, hash)
		if err != nil {
			log.Printf("failed to verify hash for doc with id: %v, hash: %x\n", doc.Id, hash)
			result = false
			continue
		}
		if !success {
			log.Printf("Blockchain check failed for doc with id: %v, hash: %x\n", doc.Id, hash)
			result = false
		}
	}

	return result, nil
}

type result struct {
	Success bool
	Hash    [32]byte
}

func executeBlockchainLightCheck(c *Client, docs []types.DocData) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results := make(map[[16]byte]result, len(docs))
	ids := make([][16]byte, len(docs))

	for _, doc := range docs {
		id := doc.Id
		hash, err := utils.CalculateHash(doc.Data)
		if err != nil {
			return false, fmt.Errorf("failed to calculate hash for doc with id: %v, error: %w", doc.Id, err)
		}

		ids = append(ids, id)

		results[id] = result{
			Success: false,
			Hash:    hash,
		}
	}

	events, err := c.dataHashRegistry.FilterHashStored(&bind.FilterOpts{Context: ctx}, ids)
	if err != nil {
		return false, fmt.Errorf("failed to get events: %w", err)
	}

	for events.Next() {
		event := events.Event
		if res, ok := results[event.Id]; ok {
			if res.Hash == event.DataHash {
				res.Success = true
				results[event.Id] = res
			}
		}
	}

	for id, res := range results {
		if !res.Success {
			return false, fmt.Errorf("blockchain check failed for doc with id: %v, hash: %x", id, res.Hash)
		}
	}

	return true, nil
}

func executeBlockchainBatchCheck(c *Client, docs []types.DocData) (bool, error) {
	return false, nil
}
