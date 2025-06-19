package workers

import (
	"fmt"
	"log"
	"time"

	"github.com/mbabinski218/BaaS-IoT-platform/blockchain"
	"github.com/mbabinski218/BaaS-IoT-platform/database"
	"github.com/mbabinski218/BaaS-IoT-platform/utils"
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

	for _, docData := range auditData {
		hash, err := utils.CalculateHash(docData.Data)
		if err != nil {
			fmt.Println("Failed to calculate hash for data Id", docData.Id, ":", err)
			continue
		}

		success, _, err := aw.blockchain.VerifyHash(docData.Id, hash, time.Time{}, nil)
		if err != nil {
			fmt.Println("Failed to update audit status for data Id", docData.Id, ":", err)
			continue
		}
		if !success {
			fmt.Println("Failed to verify hash for data with Id", docData.Id)
			continue
		}
	}

	duration := time.Since(start)
	fmt.Printf("-------- Audit completed successfully (for: %v docs)--------\n", len(auditData))
	fmt.Println("Audit duration:", duration)
}
