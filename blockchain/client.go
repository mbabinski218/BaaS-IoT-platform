package blockchain

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	dataHashRegistry "github.com/sikozonpc/ecom/blockchain/smartContracts"
)

type Client struct {
	ethClient        *ethclient.Client
	dataHashRegistry *dataHashRegistry.DataHashRegistry
	auth             *bind.TransactOpts
	address          common.Address
}

func NewEthClient(url string, privateKeyHex string, contractAddress string) (*Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	var contract *dataHashRegistry.DataHashRegistry
	var contractAddr common.Address

	if contractAddress == "" {
		contractAddr, _, contract, err = dataHashRegistry.DeployDataHashRegistry(auth, client)
		if err != nil {
			return nil, fmt.Errorf("failed to deploy contract: %w", err)
		}
		fmt.Println("Deployed new contract at:", contractAddr.Hex())
	} else {
		contractAddr = common.HexToAddress(contractAddress)
		contract, err = dataHashRegistry.NewDataHashRegistry(contractAddr, client)
		if err != nil {
			return nil, fmt.Errorf("failed to bind to contract: %w", err)
		}
		fmt.Println("Using existing contract at:", contractAddr.Hex())
	}

	return &Client{
		ethClient:        client,
		dataHashRegistry: contract,
		auth:             auth,
		address:          contractAddr,
	}, nil
}

func (c *Client) SetGasLimit(gasLimit uint64) {
	c.auth.GasLimit = gasLimit
}

func (c *Client) Send(dataId uuid.UUID, hash [32]byte, deviceId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opt := bind.TransactOpts{
		Context:  ctx,
		From:     c.auth.From,
		Signer:   c.auth.Signer,
		GasLimit: c.auth.GasLimit,
		GasPrice: c.auth.GasPrice,
	}

	transaction, err := c.dataHashRegistry.StoreHash(&opt, dataId, hash, deviceId)
	if err != nil {
		return fmt.Errorf("transaction send error: %w", err)
	}

	receipt, err := bind.WaitMined(ctx, c.ethClient, transaction)
	if err != nil {
		return fmt.Errorf("wait mined error: %w", err)
	}
	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed: %s", transaction.Hash().Hex())
	}

	fmt.Println("Transaction hash:", transaction.Hash().Hex())
	return nil
}

// func (Client) Send(iotData types.IotData) {
// 	// 1. Typ urządzenia
// 	if !IsValidDeviceType(iotData.DeviceType) {
// 		log.Println("Invalid device type")
// 		return
// 	}

// 	// 2. Czy urządzenie jest zarejestrowane

// 	// 3. Hash danych
// 	dataBytes, err := json.Marshal(iotData.Data)
// 	if err != nil {
// 		log.Println("Failed to marshal IoT data:", err)
// 		return
// 	}
// 	hash := sha256.Sum256(dataBytes)

// 	// 4. Dekoduj klucz publiczny PEM
// 	block, _ := pem.Decode([]byte(iotData.PublicKey))
// 	if block == nil {
// 		log.Println("Failed to parse PEM public key")
// 		return
// 	}

// 	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
// 	if err != nil {
// 		log.Println("Failed to parse public key:", err)
// 		return
// 	}

// 	pubKey, ok := pubInterface.(*ecdsa.PublicKey)
// 	if !ok {
// 		log.Println("Not ECDSA public key")
// 		return
// 	}

// 	// 5. Dekoduj podpis
// 	sigBytes, err := base64.StdEncoding.DecodeString(iotData.Signature)
// 	if err != nil {
// 		log.Println("Failed to decode signature:", err)
// 		return
// 	}

// 	var sig struct{ R, S *big.Int }
// 	_, err = asn1.Unmarshal(sigBytes, &sig)
// 	if err != nil {
// 		log.Println("Failed to unmarshal signature:", err)
// 		return
// 	}

// 	// 6. Weryfikacja podpisu
// 	if !ecdsa.Verify(pubKey, hash[:], sig.R, sig.S) {
// 		log.Println("Signature verification failed")
// 		return
// 	}

// 	// 7. Wysyłanie do blockchaina
// 	// TODO: smart kontrakt!!!
// }

// func IsValidDeviceType(deviceType types.DeviceType) bool {
// 	switch deviceType {
// 	case types.WeatherStation, types.AirQualitySensor, types.TemperatureSensor, types.ParkingSensor:
// 		return true
// 	default:
// 		return false
// 	}
// }
