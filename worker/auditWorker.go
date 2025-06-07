package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/sikozonpc/ecom/blockchain"
	"github.com/sikozonpc/ecom/database"
	"github.com/sikozonpc/ecom/utils"
)

type AuditWorker struct {
	Interval   int64
	Size       int64
	database   *database.Client
	blockchain *blockchain.Client
}

func NewAuditWorker(interval, size int64, db *database.Client, bc *blockchain.Client) *AuditWorker {
	return &AuditWorker{
		Interval:   interval,
		Size:       size,
		database:   db,
		blockchain: bc,
	}
}

func (aw *AuditWorker) Start() {
	ticker := time.NewTicker(time.Duration(aw.Interval) * time.Millisecond)
	defer ticker.Stop()

	log.Println("Audit worker started with interval:", aw.Interval, "ms and size:", aw.Size)

	for range ticker.C {
		aw.performAudit()
	}
}

func (aw *AuditWorker) performAudit() {
	start := time.Now()

	auditData, err := aw.database.GetAuditData(aw.Size)
	if err != nil {
		fmt.Println("Failed to get audit data:", err)
		return
	}

	for _, data := range auditData {
		hash, err := utils.CalculateHash(data)
		if err != nil {
			fmt.Println("Failed to calculate hash for data Id", data, ":", err)
			continue
		}

		success, _, err := aw.blockchain.VerifyHash(data.Id, hash)
		if err != nil {
			fmt.Println("Failed to update audit status for data Id", data.Id, ":", err)
			continue
		}
		if !success {
			fmt.Println("Failed to verify hash for data with Id", data.Id)
			continue
		}
	}

	duration := time.Since(start)
	fmt.Println("-------- Audit completed successfully --------")
	fmt.Println("Audit duration:", duration)
}
