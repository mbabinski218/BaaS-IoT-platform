package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os/exec"
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	typesEth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/mbabinski218/BaaS-IoT-platform/blockchain/smartContracts"
	"github.com/mbabinski218/BaaS-IoT-platform/configs"
	"github.com/mbabinski218/BaaS-IoT-platform/types"
	"github.com/mbabinski218/BaaS-IoT-platform/utils"
)

type Client struct {
	ethClient        *ethclient.Client
	dataHashRegistry *smartContracts.DataHashRegistry
	dataHashAddress  common.Address
	batchRegistry    *smartContracts.BatchRegistry
	batchAddress     common.Address
	auth             *bind.TransactOpts
	nonceManager     *NonceManager
	BatchStartTime   time.Time
}

func NewEthClient(url, privateKeyHex, dataHashContractAddress, batchContractAddress string) (*Client, error) {
	if configs.Envs.BlockchainMode == types.BCNone {
		log.Println("Blockchain is disabled")
		return nil, nil
	}

	switch configs.Envs.BlockchainMode {
	case types.BCFullCheck:
		log.Println("Blockchain mode: Full")
	case types.BCLightCheck:
		log.Println("Blockchain mode: Light")
	case types.BCBatchCheck:
		log.Println("Blockchain mode: Batch")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.BlockchainContextTimeout)*time.Second)
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

	if configs.Envs.BlockchainGasLimit != 0 {
		auth.GasLimit = uint64(configs.Envs.BlockchainGasLimit)
		log.Println("Blockchain gas limit set to:", auth.GasLimit)
	}

	if configs.Envs.BlockchainGasTipCap != 0 {
		auth.GasTipCap = big.NewInt(configs.Envs.BlockchainGasTipCap)
		log.Println("Blockchain gas tip cap set to:", auth.GasTipCap)
	}

	if configs.Envs.BlockchainGasFeeCap != 0 {
		auth.GasFeeCap = big.NewInt(configs.Envs.BlockchainGasFeeCap)
		log.Println("Blockchain gas fee cap set to:", auth.GasFeeCap)
	}

	var dataHashContract *smartContracts.DataHashRegistry
	var dataHashContractAddr common.Address

	if dataHashContractAddress == "" {
		dataHashContractAddr, _, dataHashContract, err = smartContracts.DeployDataHashRegistry(auth, client)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy dataHashRegistry contract: %w", err)
		}
		log.Println("Blockchain deployed new dataHashRegistry contract at:", dataHashContractAddr.Hex())
	} else {
		dataHashContractAddr = common.HexToAddress(dataHashContractAddress)
		dataHashContract, err = smartContracts.NewDataHashRegistry(dataHashContractAddr, client)
		if err != nil {
			return nil, fmt.Errorf("failed to bind dataHashRegistry contract: %w", err)
		}
		log.Println("Blockchain using existing dataHashRegistry contract at:", dataHashContractAddr.Hex())
	}

	var batchContract *smartContracts.BatchRegistry
	var batchContractAddr common.Address

	if configs.Envs.BlockchainMode == types.BCBatchCheck {
		if batchContractAddress == "" {
			batchContractAddr, _, batchContract, err = smartContracts.DeployBatchRegistry(auth, client)
			if err != nil {
				return nil, fmt.Errorf("failed to deploy batchRegistry contract: %w", err)
			}
			log.Println("Blockchain deployed new batchRegistry contract at:", batchContractAddr.Hex())
		} else {
			batchContractAddr = common.HexToAddress(batchContractAddress)
			batchContract, err = smartContracts.NewBatchRegistry(batchContractAddr, client)
			if err != nil {
				return nil, fmt.Errorf("failed to bind batchRegistry contract: %w", err)
			}
			log.Println("Blockchain using existing batchRegistry contract at:", batchContractAddr.Hex())
		}
	}

	nm := NewNonceManager()
	if err := nm.Init(client, auth.From, ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize nonce manager: %w", err)
	}

	log.Println("Ethereum client created successfully")
	return &Client{
		ethClient:        client,
		dataHashRegistry: dataHashContract,
		dataHashAddress:  dataHashContractAddr,
		batchRegistry:    batchContract,
		batchAddress:     batchContractAddr,
		auth:             auth,
		nonceManager:     nm,
		BatchStartTime:   time.Time{},
	}, nil
}

func (c *Client) createOpt(ctx context.Context) bind.TransactOpts {
	return bind.TransactOpts{
		Context:   ctx,
		From:      c.auth.From,
		Signer:    c.auth.Signer,
		GasLimit:  c.auth.GasLimit,
		GasTipCap: c.auth.GasTipCap,
		GasFeeCap: c.auth.GasFeeCap,
		Nonce:     c.nonceManager.Next(),
	}
}

func (c *Client) Send(dataId uuid.UUID, hash [32]byte, deviceId uuid.UUID) (time.Duration, time.Duration, time.Duration, error) {
	if configs.Envs.BlockchainMode == types.BCNone || configs.Envs.BlockchainMode == types.BCBatchCheck {
		return 0, 0, 0, nil
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.BlockchainContextTimeout)*2*time.Second)
	defer cancel()

	opt := c.createOpt(ctx)

	sendStart := time.Now()
	transaction, err := c.dataHashRegistry.StoreHash(&opt, dataId, hash, deviceId)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("transaction send error: %w", err)
	}
	sendDuration := time.Since(sendStart)

	var mineDuration time.Duration
	if !configs.Envs.BlockchainAsyncMode {
		mineStart := time.Now()
		receipt, err := bind.WaitMined(ctx, c.ethClient, transaction)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("wait mined error: %w", err)
		}
		if receipt.Status != 1 {
			return 0, 0, 0, fmt.Errorf("transaction failed: %s", transaction.Hash().Hex())
		}
		mineDuration = time.Since(mineStart)

		log.Println("Transaction hash:", transaction.Hash().Hex())
	} else {
		go c.mine(transaction)
	}

	duration := time.Since(start)
	return duration, sendDuration, mineDuration, nil
}

func (c *Client) mine(transaction *typesEth.Transaction) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.BlockchainContextTimeout)*2*time.Second)
	defer cancel()

	receipt, err := bind.WaitMined(ctx, c.ethClient, transaction)
	if err != nil {
		log.Printf("wait mined error: %v\n", err)
		return
	}
	if receipt.Status != 1 {
		log.Printf("transaction failed: %s\n", transaction.Hash().Hex())
		return
	}
	log.Printf("Transaction mined successfully: %s\n", transaction.Hash().Hex())
}

func (c *Client) StoreRoot(startTimestamp time.Time, root [32]byte) error {
	if configs.Envs.BlockchainMode != types.BCBatchCheck {
		return fmt.Errorf("StoreRoot is only available in BCBatchCheck mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.BlockchainContextTimeout)*2*time.Second)
	defer cancel()

	opt := c.createOpt(ctx)

	transaction, err := c.batchRegistry.StoreRoot(&opt, big.NewInt(startTimestamp.Unix()), root)
	if err != nil {
		return fmt.Errorf("transaction store root error: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, c.ethClient, transaction)
	if err != nil {
		return fmt.Errorf("wait mined error: %w", err)
	}
	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed: %s", transaction.Hash().Hex())
	}

	log.Println("Transaction hash:", transaction.Hash().Hex())

	return nil
}

func (c *Client) VerifyHash(dataId uuid.UUID, hash [32]byte, timestamp time.Time, proof [][32]byte) (bool, time.Duration, error) {
	start := time.Now()

	var success bool = true
	var err error = nil

	switch configs.Envs.BlockchainMode {
	case types.BCFullCheck:
		success, err = verifyHashFullCheck(c, dataId, hash)
	case types.BCLightCheck:
		success, err = verifyHashLightCheck(c, dataId, hash)
	case types.BCBatchCheck:
		success, err = verifyHashBatchCheck(c, timestamp, hash, proof)
	default:
		return false, 0, fmt.Errorf("unknown blockchain mode")
	}

	duration := time.Since(start)

	return success, duration, err
}

func verifyHashFullCheck(c *Client, dataId uuid.UUID, hash [32]byte) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.BlockchainContextTimeout)*time.Second)
	defer cancel()

	success, err := c.dataHashRegistry.VerifyHash(&bind.CallOpts{Context: ctx}, dataId, hash)
	if err != nil {
		return false, fmt.Errorf("failed to verify hash (full): %w", err)
	}

	return success, nil
}

func verifyHashLightCheck(c *Client, dataId uuid.UUID, hash [32]byte) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.BlockchainContextTimeout)*time.Second)
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

func verifyHashBatchCheck(c *Client, timestamp time.Time, hash [32]byte, proof [][32]byte) (bool, error) {
	if len(proof) == 0 {
		return true, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.BlockchainContextTimeout)*3*time.Second)
	defer cancel()

	interval := time.Duration(configs.Envs.BlockchainBatchInterval) * time.Second
	elapsed := timestamp.Sub(c.BatchStartTime)
	ticks := int64(elapsed / interval)
	batchTime := c.BatchStartTime.Add(time.Duration(ticks) * interval)

	success, err := c.batchRegistry.VerifyProof(&bind.CallOpts{Context: ctx}, big.NewInt(batchTime.Unix()), hash, proof)
	if err != nil {
		return false, fmt.Errorf("failed to verify hash (batch): %w", err)
	}

	return success, nil
}

func (c *Client) VerifyHashes(docs []types.DocData, fromTimestamp time.Time, toTimestamp time.Time) (bool, time.Duration, error) {
	start := time.Now()

	var success bool = true
	var err error = nil

	switch configs.Envs.BlockchainMode {
	case types.BCFullCheck:
		success, err = executeBlockchainFullCheck(c, docs)
	case types.BCLightCheck:
		success, err = executeBlockchainLightCheck(c, docs)
	case types.BCBatchCheck:
		success, err = executeBlockchainBatchCheck(c, docs, fromTimestamp, toTimestamp)
	default:
		return false, 0, fmt.Errorf("unknown blockchain mode")
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

		success, _, err := c.VerifyHash(doc.Id, hash, time.Time{}, nil)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.BlockchainContextTimeout)*time.Second)
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

func executeBlockchainBatchCheck(c *Client, docs []types.DocData, from, to time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.BlockchainContextTimeout)*time.Second)
	defer cancel()

	err := verifyTimestamps(c.BatchStartTime, from, to)
	if err != nil {
		return false, fmt.Errorf("invalid batch timestamps: %w", err)
	}

	interval := time.Duration(configs.Envs.BlockchainBatchInterval) * time.Second
	slots := make(map[time.Time][]types.DocData)

	for _, doc := range docs {
		timestampStr, ok := doc.Data["timestamp"].(string)
		if !ok {
			continue
		}

		timestamp, err := time.Parse(types.TimeLayout, timestampStr)
		if err != nil {
			continue
		}

		if timestamp.Before(from) || !timestamp.Before(to) {
			continue
		}

		offset := timestamp.Sub(from)
		slotIndex := int(offset / interval)
		slotStart := from.Add(time.Duration(slotIndex) * interval)

		slots[slotStart] = append(slots[slotStart], doc)
	}

	roots := make(map[time.Time][32]byte)
	for slotStart, slotDocs := range slots {
		if len(slotDocs) == 0 {
			continue
		}
		root, _, err := utils.CreateMerkleRoot(slotDocs)
		if err != nil {
			return false, fmt.Errorf("failed to create Merkle root for slot starting at %v: %w", slotStart, err)
		}
		roots[slotStart] = root
	}

	for slotStart, root := range roots {
		success, err := c.batchRegistry.VerifyRoot(&bind.CallOpts{Context: ctx}, big.NewInt(slotStart.Unix()), root)
		if err != nil {
			return false, fmt.Errorf("failed to verify root for slot starting at %v: %w", slotStart, err)
		}
		if !success {
			return false, fmt.Errorf("blockchain check failed for slot starting at %v with root %x", slotStart, root)
		}
	}

	return true, nil
}

func verifyTimestamps(batchStartTime, from, to time.Time) error {
	if from.Before(batchStartTime) || to.Before(batchStartTime) {
		return fmt.Errorf("timestamps must be after the batch start time: %v", batchStartTime)
	}

	fromDiff := from.Sub(batchStartTime)
	toDiff := to.Sub(batchStartTime)

	interval := time.Duration(configs.Envs.BlockchainBatchInterval) * time.Second
	if fromDiff%interval != 0 || toDiff%interval != 0 {
		return fmt.Errorf("timestamps must be aligned with the batch interval: %v and start time: %v", interval, batchStartTime)
	}

	return nil
}

func (c *Client) GetBlockNumber() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.BlockchainContextTimeout)*time.Second)
	defer cancel()

	blockNumber, err := c.ethClient.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get block number: %w", err)
	}
	return blockNumber, nil
}

func (c *Client) StopMining() error {
	cmd := exec.Command("ssh", "-p "+configs.Envs.BlockchainServerPort, configs.Envs.BlockchainServerIP, "docker", "pause", configs.Envs.BlockchainValidators)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop container: %v - output: %s", err, output)
	}
	return nil
}

func (c *Client) StartMining() error {
	cmd := exec.Command("ssh", "-p "+configs.Envs.BlockchainServerPort, configs.Envs.BlockchainServerIP, "docker", "unpause", configs.Envs.BlockchainValidators)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start container: %v - output: %s", err, output)
	}

	return nil
}

func (c *Client) IsMining() bool {
	cmd := exec.Command("ssh", "-p "+configs.Envs.BlockchainServerPort, configs.Envs.BlockchainServerIP, "docker", "inspect", "--format={{.State.Running}}", configs.Envs.BlockchainValidators)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	lines := []string{}
	for _, line := range string(output) {
		if line == '\r' || line == '\n' {
			continue
		}
		lines = append(lines, string(line))
	}
	return !slices.Contains(lines, "false")
}
