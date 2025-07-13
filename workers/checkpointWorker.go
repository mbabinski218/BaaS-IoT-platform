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
	numberOfRepeats := configs.Envs.BlockchainCheckpointCallRepeats
	apiURL := "http://" + configs.Envs.PublicHost + "/api/v1/get"

	// Create or open Excel file
	fileName := "getByTime_checkpoint_results.xlsx"
	var f *excelize.File
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		f = excelize.NewFile()
		f.SetSheetName("Sheet1", "Results")
		f.SetCellValue("Results", "A1", "Timestamp")
		f.SetCellValue("Results", "B1", "Blocks")
		f.SetCellValue("Results", "C1", "First - Db duration")
		f.SetCellValue("Results", "D1", "First - Blockchain duration")
		f.SetCellValue("Results", "E1", "First - Total duration")
		f.SetCellValue("Results", "F1", "First - Missed")
		f.SetCellValue("Results", "G1", "Center - Db duration")
		f.SetCellValue("Results", "H1", "Center - Blockchain duration")
		f.SetCellValue("Results", "I1", "Center - Total duration")
		f.SetCellValue("Results", "J1", "Center - Missed")
		f.SetCellValue("Results", "K1", "Last - Db duration")
		f.SetCellValue("Results", "L1", "Last - Blockchain duration")
		f.SetCellValue("Results", "M1", "Last - Total duration")
		f.SetCellValue("Results", "N1", "Last - Missed")
	} else {
		return fmt.Errorf("file %s already exists", fileName)
	}

	// Find next empty row
	rows, err := f.GetRows("Results")
	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}
	rowNum := len(rows) + 1

	// Get by time
	fromStart := cw.startTime.Add(time.Duration(configs.Envs.BlockchainBatchInterval) * time.Second)
	toStart := fromStart.Add(time.Duration(configs.Envs.BlockchainBatchInterval) * time.Second)
	fromCenter := cw.startTime.Add(time.Duration(configs.Envs.BlockchainBatchInterval*int64(blockNumber/2)) * time.Second)
	toCenter := fromCenter.Add(time.Duration(configs.Envs.BlockchainBatchInterval) * time.Second)
	fromEnd := cw.startTime.Add(time.Duration(configs.Envs.BlockchainBatchInterval*int64(blockNumber)) * time.Second)
	toEnd := fromEnd.Add(time.Duration(configs.Envs.BlockchainBatchInterval) * time.Second)

	fromStartStr := fromStart.Format(types.TimeLayout)
	toStartStr := toStart.Format(types.TimeLayout)
	fromCenterStr := fromCenter.Format(types.TimeLayout)
	toCenterStr := toCenter.Format(types.TimeLayout)
	fromEndStr := fromEnd.Format(types.TimeLayout)
	toEndStr := toEnd.Format(types.TimeLayout)

	for range numberOfRepeats {
		// Write data
		f.SetCellValue("Results", fmt.Sprintf("A%d", rowNum), time.Now().Format(types.TimeLayout))
		f.SetCellValue("Results", fmt.Sprintf("B%d", rowNum), blockNumber)

		// API
		var data map[string]any

		resp, err := http.Get(apiURL + fmt.Sprintf("?from=%s&to=%s", fromStartStr, toStartStr))
		if err != nil {
			return fmt.Errorf("API request failed: %w", err)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}

		// Write data
		f.SetCellValue("Results", fmt.Sprintf("C%d", rowNum), data[types.MongoDuration])
		f.SetCellValue("Results", fmt.Sprintf("D%d", rowNum), data[types.BlockchainDuration])
		f.SetCellValue("Results", fmt.Sprintf("E%d", rowNum), data[types.TotalDuration])
		f.SetCellValue("Results", fmt.Sprintf("F%d", rowNum), data[types.Missed])

		// API
		resp, err = http.Get(apiURL + fmt.Sprintf("?from=%s&to=%s", fromCenterStr, toCenterStr))
		if err != nil {
			return fmt.Errorf("API request failed: %w", err)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}

		// Write data
		f.SetCellValue("Results", fmt.Sprintf("G%d", rowNum), data[types.MongoDuration])
		f.SetCellValue("Results", fmt.Sprintf("H%d", rowNum), data[types.BlockchainDuration])
		f.SetCellValue("Results", fmt.Sprintf("I%d", rowNum), data[types.TotalDuration])
		f.SetCellValue("Results", fmt.Sprintf("J%d", rowNum), data[types.Missed])

		// API
		resp, err = http.Get(apiURL + fmt.Sprintf("?from=%s&to=%s", fromEndStr, toEndStr))
		if err != nil {
			return fmt.Errorf("API request failed: %w", err)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}

		// Write data
		f.SetCellValue("Results", fmt.Sprintf("K%d", rowNum), data[types.MongoDuration])
		f.SetCellValue("Results", fmt.Sprintf("L%d", rowNum), data[types.BlockchainDuration])
		f.SetCellValue("Results", fmt.Sprintf("M%d", rowNum), data[types.TotalDuration])
		f.SetCellValue("Results", fmt.Sprintf("N%d", rowNum), data[types.Missed])

		rowNum++
	}

	// Save file
	if err := f.SaveAs(fileName); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	// Create or open Excel file
	fileName = "getById_checkpoint_results.xlsx"
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
	rows, err = f.GetRows("Results")
	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}
	rowNum = len(rows) + 1

	// Get by id
	firstDocId, err := cw.database.GetFirstDocumentId()
	if err != nil {
		return fmt.Errorf("failed to get first document ID: %w", err)
	}
	centerDocId, err := cw.database.GetCenterDocumentId()
	if err != nil {
		return fmt.Errorf("failed to get center document ID: %w", err)
	}
	lastDocId, err := cw.database.GetLastDocumentId()
	if err != nil {
		return fmt.Errorf("failed to get last document ID: %w", err)
	}

	for range numberOfRepeats {
		// Write data
		f.SetCellValue("Results", fmt.Sprintf("A%d", rowNum), time.Now().Format(types.TimeLayout))
		f.SetCellValue("Results", fmt.Sprintf("B%d", rowNum), blockNumber)

		// API
		var data map[string]any

		resp, err := http.Get(apiURL + fmt.Sprintf("/%s", firstDocId.String()))
		if err != nil {
			return fmt.Errorf("API request failed: %w", err)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}

		// Write data
		f.SetCellValue("Results", fmt.Sprintf("C%d", rowNum), data[types.MongoDuration])
		f.SetCellValue("Results", fmt.Sprintf("D%d", rowNum), data[types.BlockchainDuration])
		f.SetCellValue("Results", fmt.Sprintf("E%d", rowNum), data[types.TotalDuration])

		// API
		resp, err = http.Get(apiURL + fmt.Sprintf("/%s", centerDocId.String()))
		if err != nil {
			return fmt.Errorf("API request failed: %w", err)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}

		// Write data
		f.SetCellValue("Results", fmt.Sprintf("F%d", rowNum), data[types.MongoDuration])
		f.SetCellValue("Results", fmt.Sprintf("G%d", rowNum), data[types.BlockchainDuration])
		f.SetCellValue("Results", fmt.Sprintf("H%d", rowNum), data[types.TotalDuration])

		// API
		resp, err = http.Get(apiURL + fmt.Sprintf("/%s", lastDocId.String()))
		if err != nil {
			return fmt.Errorf("API request failed: %w", err)
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}

		// Write data
		f.SetCellValue("Results", fmt.Sprintf("I%d", rowNum), data[types.MongoDuration])
		f.SetCellValue("Results", fmt.Sprintf("J%d", rowNum), data[types.BlockchainDuration])
		f.SetCellValue("Results", fmt.Sprintf("K%d", rowNum), data[types.TotalDuration])

		rowNum++
	}

	// Save file
	if err := f.SaveAs(fileName); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}
