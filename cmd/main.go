package main

import (
	"log"
	"time"

	"github.com/mbabinski218/BaaS-IoT-platform/api"
	"github.com/mbabinski218/BaaS-IoT-platform/blockchain"
	"github.com/mbabinski218/BaaS-IoT-platform/configs"
	"github.com/mbabinski218/BaaS-IoT-platform/database"
	"github.com/mbabinski218/BaaS-IoT-platform/runners"
	"github.com/mbabinski218/BaaS-IoT-platform/types"
	"github.com/mbabinski218/BaaS-IoT-platform/workers"
)

func main() {
	// Database initialization
	databaseClient, err := database.Connect(configs.Envs.MongoDbUri, configs.Envs.MongoDbName, configs.Envs.MongoDbCollectionName)
	if err != nil {
		log.Fatal(err)
	}

	// Blockchain client initialization
	ethClient, err := blockchain.NewEthClient(configs.Envs.BlockchainUrl, configs.Envs.BlockchainPrivateKey, configs.Envs.BlockchainContractAddress, configs.Envs.BlockchainBatchContractAddress)
	if err != nil {
		log.Fatal(err)
	}

	// Iot simulator initialization
	if err := runners.RunIotSimulator(ethClient); err != nil {
		log.Fatal(err)
	}

	// Workers initialization
	startTime := time.Now()
	backgroundWorkers := []workers.Worker{}

	if configs.Envs.AuditEnabled {
		backgroundWorkers = append(backgroundWorkers, workers.NewAuditWorker(configs.Envs.AuditTimeout, configs.Envs.AuditSize, databaseClient, ethClient))
	}

	if configs.Envs.BlockchainMode == types.BCBatchCheck {
		backgroundWorkers = append(backgroundWorkers, workers.NewBatchWorker(configs.Envs.BlockchainBatchInterval, databaseClient, ethClient, &startTime))
	}

	if len(configs.Envs.BlockchainCheckpoints) > 0 {
		backgroundWorkers = append(backgroundWorkers, workers.NewCheckpointWorker(databaseClient, ethClient, &startTime))
	}

	for _, worker := range backgroundWorkers {
		go worker.Start()
	}

	// Start the API server
	server := api.NewAPIServer(configs.Envs.PublicHost, ethClient, databaseClient)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
