package workers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/mbabinski218/BaaS-IoT-platform/blockchain"
	"github.com/mbabinski218/BaaS-IoT-platform/configs"
	"github.com/mbabinski218/BaaS-IoT-platform/database"
	"github.com/mbabinski218/BaaS-IoT-platform/types"
	"github.com/xuri/excelize/v2"
)

type CheckpointWorker struct {
	Interval   uint64
	database   *database.Client
	blockchain *blockchain.Client
	startTime  *time.Time
}

func NewCheckpointWorker(db *database.Client, bc *blockchain.Client, st *time.Time) *CheckpointWorker {
	return &CheckpointWorker{
		database:   db,
		blockchain: bc,
		startTime:  st,
	}
}

func (cw *CheckpointWorker) Start() {
	interval := configs.Envs.BlockchainSecondsPerBlock

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	log.Println("Checkpoint worker started with interval:", interval, "seconds")

	for range ticker.C {
		cw.performTest()
	}
}

func (cw *CheckpointWorker) performTest() {
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

	err = cw.test(blockNumber)
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

func (cw *CheckpointWorker) test(blockNumber uint64) error {
	const numberOfRepeats = 10
	apiURL := "http://" + configs.Envs.PublicHost + "/api/v1/get"

	// Create or open Excel file
	date := time.Now()
	fileName := date.Format(types.ShortTimeLayout) + "_checkpoint_results.xlsx"
	var f *excelize.File
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		f = excelize.NewFile()
		f.SetSheetName("Sheet1", "Results")
		f.SetCellValue("Results", "A1", "Timestamp")
		f.SetCellValue("Results", "B1", "Blocks")
		f.SetCellValue("Results", "C1", "First - Db duration")
		f.SetCellValue("Results", "D1", "First - Blockchain duration")
		f.SetCellValue("Results", "E1", "First - Total duration")
		f.SetCellValue("Results", "F1", "Center - Db duration")
		f.SetCellValue("Results", "G1", "Center - Blockchain duration")
		f.SetCellValue("Results", "H1", "Center - Total duration")
		f.SetCellValue("Results", "I1", "Last - Db duration")
		f.SetCellValue("Results", "J1", "Last - Blockchain duration")
		f.SetCellValue("Results", "K1", "Last - Total duration")
	} else {
		return fmt.Errorf("file %s already exists", fileName)
	}

	// Find next empty row
	rows, err := f.GetRows("Results")
	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}
	rowNum := len(rows) + 1

	for range numberOfRepeats {
		// Api call for first, center, and last blocks

		// Write data
		f.SetCellValue("Results", fmt.Sprintf("A%d", rowNum), time.Now().Format(types.TimeLayout))
		f.SetCellValue("Results", fmt.Sprintf("B%d", rowNum), blockNumber)
	}

	// Api request
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Save file
	if err := f.SaveAs(fileName); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}
