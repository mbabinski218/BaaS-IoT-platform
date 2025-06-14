package main

import (
	"log"

	"github.com/mbabinski218/BaaS-IoT-platform/api"
	"github.com/mbabinski218/BaaS-IoT-platform/blockchain"
	"github.com/mbabinski218/BaaS-IoT-platform/configs"
	"github.com/mbabinski218/BaaS-IoT-platform/database"
	"github.com/mbabinski218/BaaS-IoT-platform/types"
	"github.com/mbabinski218/BaaS-IoT-platform/worker"
)

func main() {
	// Database initialization
	databaseClient, err := database.Connect(configs.Envs.MongoDbUri, configs.Envs.MongoDbName, configs.Envs.MongoDbCollectionName)
	if err != nil {
		log.Fatal(err)
	}

	// Blockchain client initialization
	ethClient, err := blockchain.NewEthClient(configs.Envs.BlockchainUrl, configs.Envs.BlockchainPrivateKey, configs.Envs.BlockchainContractAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Workers initialization
	workers := []worker.Worker{}

	if configs.Envs.AuditEnabled {
		workers = append(workers, worker.NewAuditWorker(configs.Envs.AuditTimeout, configs.Envs.AuditSize, databaseClient, ethClient))
	}

	if configs.Envs.BlockchainMode == types.BCBatchCheck {
		workers = append(workers, worker.NewBatchWorker(configs.Envs.BlockchainBatchInterval, databaseClient, ethClient))
	}

	for _, w := range workers {
		go w.Start()
	}

	// Start the API server
	server := api.NewAPIServer(configs.Envs.PublicHost, ethClient, databaseClient)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
