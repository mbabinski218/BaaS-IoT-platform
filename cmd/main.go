package main

import (
	"log"

	"github.com/sikozonpc/ecom/api"
	"github.com/sikozonpc/ecom/blockchain"
	"github.com/sikozonpc/ecom/configs"
	"github.com/sikozonpc/ecom/database"
	"github.com/sikozonpc/ecom/worker"
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

	workers := []worker.Worker{
		worker.NewAuditWorker(configs.Envs.AuditTimeout, configs.Envs.AuditSize, databaseClient, ethClient),
	}
	for _, w := range workers {
		go w.Start()
	}

	server := api.NewAPIServer(configs.Envs.PublicHost, ethClient, databaseClient)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
