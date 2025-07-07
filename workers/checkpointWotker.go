package workers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/mbabinski218/BaaS-IoT-platform/blockchain"
	"github.com/mbabinski218/BaaS-IoT-platform/configs"
	"github.com/mbabinski218/BaaS-IoT-platform/database"
	"github.com/xuri/excelize/v2"
)

type CheckpointWorker struct {
	Interval   uint64
	database   *database.Client
	blockchain *blockchain.Client
}

func NewCheckpointWorker(db *database.Client, bc *blockchain.Client) *CheckpointWorker {
	return &CheckpointWorker{
		database:   db,
		blockchain: bc,
	}
}

func (cw *CheckpointWorker) Start() {
	interval := configs.Envs.BlockchainSecondsPerBlock

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	log.Println("Checkpoint worker started with interval:", interval, "seconds")

	for range ticker.C {
		cw.performBatch()
	}
}

func (cw *CheckpointWorker) performBatch() {
	start := time.Now()

	blockNumber, err := cw.blockchain.GetBlockNumber()
	if err != nil {
		log.Fatalln("Error getting block number:", err)
		return
	}

	if !slices.Contains(configs.Envs.BlockchainCheckpoints, blockNumber) {
		return
	}

	if err = cw.blockchain.StopMining(); err != nil {
		log.Println("Error stopping mining:", err)
		return
	}

	err = cw.test()
	if err != nil {
		log.Println("Test failed:", err)
		return
	}

	if err = cw.blockchain.StartMining(); err != nil {
		log.Println("Error starting mining:", err)
		return
	}

	duration := time.Since(start)
	fmt.Println("Checkpoint duration:", duration)
}

func (cw *CheckpointWorker) test() error {
	apiURL := "http://" + configs.Envs.PublicHost + "/api/v1/get"

	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Optionally, parse the response if needed
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode API response: %w", err)
	}

	// Create or open Excel file
	date := time.Now()
	fileName := date.Format("2006-01-02") + "_checkpoint_results.xlsx"
	var f *excelize.File
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		f = excelize.NewFile()
		f.SetSheetName("Sheet1", "Results")
		// Write header
		f.SetCellValue("Results", "A1", "Timestamp")
		f.SetCellValue("Results", "B1", "Duration(ms)")
		f.SetCellValue("Results", "C1", "StatusCode")
	} else {
		f, err = excelize.OpenFile(fileName)
		if err != nil {
			return fmt.Errorf("failed to open Excel file: %w", err)
		}
	}

	// Find next empty row
	rows, err := f.GetRows("Results")
	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}
	rowNum := len(rows) + 1

	// Write data
	f.SetCellValue("Results", fmt.Sprintf("A%d", rowNum), time.Now().Format(time.RFC3339))
	f.SetCellValue("Results", fmt.Sprintf("C%d", rowNum), resp.StatusCode)

	// Save file
	if err := f.SaveAs(fileName); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}
