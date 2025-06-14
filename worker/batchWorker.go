package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/mbabinski218/BaaS-IoT-platform/blockchain"
	"github.com/mbabinski218/BaaS-IoT-platform/database"
	"github.com/mbabinski218/BaaS-IoT-platform/types"
)

type BatchWorker struct {
	Interval   int64
	database   *database.Client
	blockchain *blockchain.Client
}

func NewBatchWorker(interval int64, db *database.Client, bc *blockchain.Client) *BatchWorker {
	return &BatchWorker{
		Interval:   interval,
		database:   db,
		blockchain: bc,
	}
}

func (bw *BatchWorker) Start() {
	startTime, _ := time.ParseInLocation(types.TimeLayout, "2025-01-01 00:00:00", time.UTC)
	elapsed := time.Since(startTime)

	interval := time.Duration(bw.Interval) * time.Minute

	ticks := elapsed / interval
	nextTick := startTime.Add((ticks + 1) * interval)

	time.Sleep(time.Until(nextTick))

	ticker := time.NewTicker(time.Duration(bw.Interval) * time.Minute)
	defer ticker.Stop()

	bw.blockchain.BatchStartTime = startTime

	log.Println("Batch worker started at:", startTime, "with interval:", bw.Interval, "minutes")

	for range ticker.C {
		bw.performBatch()
	}
}

func (bw *BatchWorker) performBatch() {
	start := time.Now()

	now := time.Now().Truncate(time.Minute)
	data, _, err := bw.database.GetFromTo(now, now.Add(time.Duration(bw.Interval)*time.Minute))
	if err != nil {
		log.Println("Failed to get batch data:", err)
		return
	}

	if len(data) == 0 {
		log.Println("No data to process in batch")
		return
	}

	duration := time.Since(start)
	fmt.Printf("-------- Batch completed successfully (for: %v docs)--------\n", len(data))
	fmt.Println("Batch duration:", duration)
}
