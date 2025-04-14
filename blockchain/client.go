package client

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

func verifySignature(data []byte, signatureHex string, pubKeyHex string) bool {
	sig, err := hex.DecodeString(strings.TrimPrefix(signatureHex, "0x"))
	if err != nil {
		log.Println("Invalid signature:", err)
		return false
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	pubKeyBytes, err := hex.DecodeString(strings.TrimPrefix(pubKeyHex, "0x"))
	if err != nil {
		log.Println("Invalid public key:", err)
		return false
	}
	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		log.Println("Failed to unmarshal public key:", err)
		return false
	}

	hash := crypto.Keccak256(data)
	return crypto.VerifySignature(crypto.FromECDSAPub(pubKey), hash, sig[:64])
}

func sendToSmartContract(encrypted []byte, sender common.Address, client *ethclient.Client, privateKey *ecdsa.PrivateKey) error {
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111)) // Sepolia
	if err != nil {
		return err
	}

	// TODO: odpowiedni kontrakt (ABI, adres kontraktu)
	address := common.HexToAddress("0xYourContractAddress")
	instance, err := NewYourContract(address, client)
	if err != nil {
		return err
	}

	tx, err := instance.StoreEncryptedData(auth, sender, encrypted)
	if err != nil {
		return err
	}

	log.Println("TX sent:", tx.Hash().Hex())
	return nil
}

func (Client) Send(data types.IotData) {
	if !verifySignature(data) {
		log.Fatal("Signature verification failed")
	}

	fmt.Println("--------------Podpis poprawny---------------")

	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal("Failed to connect to Ethereum node:", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA("YOUR_BACKEND_PRIVATE_KEY")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)
	nonce, _ := client.PendingNonceAt(context.Background(), fromAddress)
	chainID, _ := client.ChainID(context.Background())
	gasPrice, _ := client.SuggestGasPrice(context.Background())

	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	contractAddress := common.HexToAddress("...")
	contractAbi, err := abi.JSON(strings.NewReader(`[{"inputs":[{"internalType":"string","name":"_data","type":"string"}],"name":"storeData","outputs":[],"stateMutability":"nonpayable","type":"function"}]`))
	if err != nil {
		log.Fatal(err)
	}

	inputData, err := contractAbi.Pack("storeData", data)
	if err != nil {
		log.Fatal(err)
	}
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), auth.GasLimit, gasPrice, inputData)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transakcja wysłana:", signedTx.Hash().Hex())
}
