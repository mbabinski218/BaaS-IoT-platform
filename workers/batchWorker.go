package workers

import (
	"fmt"
	"log"
	"time"

	"github.com/mbabinski218/BaaS-IoT-platform/blockchain"
	"github.com/mbabinski218/BaaS-IoT-platform/database"
	"github.com/mbabinski218/BaaS-IoT-platform/types"
	"github.com/mbabinski218/BaaS-IoT-platform/utils"
)

type BatchWorker struct {
	Interval   int64
	database   *database.Client
	blockchain *blockchain.Client
	startTime  *time.Time
}

func NewBatchWorker(interval int64, db *database.Client, bc *blockchain.Client, st *time.Time) *BatchWorker {
	return &BatchWorker{
		Interval:   interval,
		database:   db,
		blockchain: bc,
		startTime:  st,
	}
}

func (bw *BatchWorker) Start() {
	startTime, _ := time.ParseInLocation(types.TimeLayout, types.BlockchainBatchStartTime, time.UTC)
	elapsed := time.Since(startTime)

	interval := time.Duration(bw.Interval) * time.Second

	ticks := elapsed / interval
	nextTick := startTime.Add((ticks + 1) * interval)

	time.Sleep(time.Until(nextTick))
	bw.startTime = &nextTick

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	bw.blockchain.BatchStartTime = startTime

	log.Println("Batch worker started at:", startTime, "with interval:", bw.Interval, "seconds")

	for range ticker.C {
		bw.performBatch()
	}
}

func (bw *BatchWorker) performBatch() {
	start := time.Now()

	now := time.Now().UTC().Add(2 * time.Hour).Truncate(time.Second)
	lastTick := now.Add(-time.Duration(bw.Interval) * time.Second)
	docs, _, err := bw.database.GetFromTo(lastTick, now)
	if err != nil {
		log.Println("Failed to get batch data:", err)
		return
	}

	if len(docs) == 0 {
		log.Println("No data to process in batch")
		return
	}

	root, audit, err := utils.CreateMerkleRoot(docs)
	if err != nil {
		log.Println("Failed to create Merkle root:", err)
		return
	}

	if err := bw.blockchain.StoreRoot(lastTick, root); err != nil {
		log.Println("Failed to store root in blockchain:", err)
		return
	}

	for _, doc := range docs {
		if err := bw.database.UpdateProof(doc.Id, audit[doc.Id]); err != nil {
			log.Println("Failed to store Merkle audit info for doc with id:", doc.Id, "err:", err)
		}
	}

	duration := time.Since(start)
	fmt.Printf("-------- Batch completed successfully (for: %v docs)--------\n", len(docs))
	fmt.Println("Batch duration:", duration)
}
