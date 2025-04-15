package client

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sikozonpc/ecom/types"
)

type Client struct {
	ethClient *ethclient.Client
}

func NewEthClient(url string) (*Client, error) {
	newClient, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	blockNumber, err := newClient.BlockNumber(nil)
	if err != nil {
		log.Fatalf("Failed to get block number: %v", err)
	}

	fmt.Printf("Current block number: %d\n", blockNumber)

	return &Client{
		ethClient: newClient,
	}, err
}

func (Client) Send(iotData types.IotData) {
	// 1. Typ urządzenia
	if !IsValidDeviceType(iotData.DeviceType) {
		log.Println("Invalid device type")
		return
	}

	// 2. Czy urządzenie jest zarejestrowane

	// 3. Hash danych
	dataBytes, err := json.Marshal(iotData.Data)
	if err != nil {
		log.Println("Failed to marshal IoT data:", err)
		return
	}
	hash := sha256.Sum256(dataBytes)

	// 4. Dekoduj klucz publiczny PEM
	block, _ := pem.Decode([]byte(iotData.PublicKey))
	if block == nil {
		log.Println("Failed to parse PEM public key")
		return
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Println("Failed to parse public key:", err)
		return
	}

	pubKey, ok := pubInterface.(*ecdsa.PublicKey)
	if !ok {
		log.Println("Not ECDSA public key")
		return
	}

	// 5. Dekoduj podpis
	sigBytes, err := base64.StdEncoding.DecodeString(iotData.Signature)
	if err != nil {
		log.Println("Failed to decode signature:", err)
		return
	}

	var sig struct{ R, S *big.Int }
	_, err = asn1.Unmarshal(sigBytes, &sig)
	if err != nil {
		log.Println("Failed to unmarshal signature:", err)
		return
	}

	// 6. Weryfikacja podpisu
	if !ecdsa.Verify(pubKey, hash[:], sig.R, sig.S) {
		log.Println("Signature verification failed")
		return
	}

	// 7. Wysyłanie do blockchaina
	// TODO: smart kontrakt!!!
}

func IsValidDeviceType(deviceType types.DeviceType) bool {
	switch deviceType {
	case types.WeatherStation, types.AirQualitySensor, types.TemperatureSensor:
		return true
	default:
		return false
	}
}
