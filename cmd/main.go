package main

import (
	"log"

	"github.com/sikozonpc/ecom/api"
	"github.com/sikozonpc/ecom/blockchain"
	"github.com/sikozonpc/ecom/configs"
	"github.com/sikozonpc/ecom/database"
)

func main() {
	databaseClient, err := database.Connect(configs.Envs.MongoDbUri, configs.Envs.MongoDbName, configs.Envs.MongoDbCollectionName)
	if err != nil {
		log.Fatal(err)
	}

	ethClient, err := blockchain.NewEthClient(configs.Envs.BlockchainUrl, configs.Envs.BlockchainPrivateKey, configs.Envs.BlockchainContractAddress)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(configs.Envs.PublicHost, ethClient, databaseClient)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
