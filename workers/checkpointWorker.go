package workers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
	"time"

	"github.com/mbabinski218/BaaS-IoT-platform/blockchain"
	"github.com/mbabinski218/BaaS-IoT-platform/configs"
	"github.com/mbabinski218/BaaS-IoT-platform/database"
	"github.com/mbabinski218/BaaS-IoT-platform/services"
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
	}

	err = cw.Test(blockNumber)
	if err != nil {
		log.Println("Test failed:", err)
	}

	if err = cw.blockchain.StartMining(); err != nil {
		log.Println("Error starting mining:", err)
	}

	duration := time.Since(start)
	fmt.Println("Checkpoint duration:", duration)
}

func (cw *CheckpointWorker) Test(blockNumber uint64) error {
	numberOfRepeats := configs.Envs.BlockchainCheckpointCallRepeats
	handler := services.NewHandler(cw.blockchain, cw.database)
	apiURL := fmt.Sprintf("http://%s/Get", configs.Envs.PublicHost)

	// Create or open Excel file
	fileName := fmt.Sprintf("getByTime_checkpoint_results_%s.xlsx", time.Now().Format("20060102_150405"))
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
	fromStart := cw.startTime.Add(time.Duration(configs.Envs.BlockchainBatchInterval*5) * time.Second).Truncate(time.Second)
	toStart := fromStart.Add(time.Duration(configs.Envs.BlockchainBatchInterval) * time.Second).Truncate(time.Second)
	fromCenter := cw.startTime.Add(time.Duration(configs.Envs.BlockchainBatchInterval*int64(blockNumber/2)) * time.Second).Truncate(time.Second)
	toCenter := fromCenter.Add(time.Duration(configs.Envs.BlockchainBatchInterval) * time.Second).Truncate(time.Second)
	fromEnd := cw.startTime.Add(time.Duration(configs.Envs.BlockchainBatchInterval*int64(blockNumber-5)) * time.Second).Truncate(time.Second)
	toEnd := fromEnd.Add(time.Duration(configs.Envs.BlockchainBatchInterval) * time.Second).Truncate(time.Second)

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
		var respWriter http.ResponseWriter

		req := &http.Request{
			Method: "GET",
			URL: &url.URL{
				Path:     apiURL,
				RawQuery: fmt.Sprintf("from=%s&to=%s", url.QueryEscape(fromStartStr), url.QueryEscape(toStartStr)),
			},
		}
		resp := handler.HandleGetFromTo(nil, req)
		if resp == nil {
			return fmt.Errorf("1.API request failed: resp is nil")
		}

		// Write data
		f.SetCellValue("Results", fmt.Sprintf("C%d", rowNum), resp[types.MongoDuration])
		f.SetCellValue("Results", fmt.Sprintf("D%d", rowNum), resp[types.BlockchainDuration])
		f.SetCellValue("Results", fmt.Sprintf("E%d", rowNum), resp[types.TotalDuration])
		f.SetCellValue("Results", fmt.Sprintf("F%d", rowNum), resp[types.Missed])

		// API
		req = &http.Request{
			Method: "GET",
			URL: &url.URL{
				Path:     apiURL,
				RawQuery: fmt.Sprintf("from=%s&to=%s", url.QueryEscape(fromCenterStr), url.QueryEscape(toCenterStr)),
			},
		}
		resp = handler.HandleGetFromTo(nil, req)
		if resp == nil {
			return fmt.Errorf("2.API request failed: resp is nil")
		}

		// Write data
		f.SetCellValue("Results", fmt.Sprintf("G%d", rowNum), resp[types.MongoDuration])
		f.SetCellValue("Results", fmt.Sprintf("H%d", rowNum), resp[types.BlockchainDuration])
		f.SetCellValue("Results", fmt.Sprintf("I%d", rowNum), resp[types.TotalDuration])
		f.SetCellValue("Results", fmt.Sprintf("J%d", rowNum), resp[types.Missed])

		// API
		req = &http.Request{
			Method: "GET",
			URL: &url.URL{
				Path:     apiURL,
				RawQuery: fmt.Sprintf("from=%s&to=%s", url.QueryEscape(fromEndStr), url.QueryEscape(toEndStr)),
			},
		}
		resp = handler.HandleGetFromTo(respWriter, req)
		if resp == nil {
			return fmt.Errorf("3.API request failed: resp is nil")
		}

		// Write data
		f.SetCellValue("Results", fmt.Sprintf("K%d", rowNum), resp[types.MongoDuration])
		f.SetCellValue("Results", fmt.Sprintf("L%d", rowNum), resp[types.BlockchainDuration])
		f.SetCellValue("Results", fmt.Sprintf("M%d", rowNum), resp[types.TotalDuration])
		f.SetCellValue("Results", fmt.Sprintf("N%d", rowNum), resp[types.Missed])

		rowNum++
	}

	// Save file
	if err := f.SaveAs(fileName); err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	// Create or open Excel file
	fileName = fmt.Sprintf("getById_checkpoint_results_%s.xlsx", time.Now().Format("20060102_150405"))
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
		return fmt.Errorf("API request failed: resp is nil")
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
		var respWriter http.ResponseWriter

		req := &http.Request{
			Method: "GET",
			URL: &url.URL{
				Path: fmt.Sprintf(apiURL, "/%s", firstDocId.String()),
			},
		}
		resp := handler.HandleGet(respWriter, req)
		if resp == nil {
			return fmt.Errorf("4.API request failed: resp is nil")
		}

		// Write data
		f.SetCellValue("Results", fmt.Sprintf("C%d", rowNum), resp[types.MongoDuration])
		f.SetCellValue("Results", fmt.Sprintf("D%d", rowNum), resp[types.BlockchainDuration])
		f.SetCellValue("Results", fmt.Sprintf("E%d", rowNum), resp[types.TotalDuration])

		// API
		req = &http.Request{
			Method: "GET",
			URL: &url.URL{
				Path: fmt.Sprintf(apiURL, "/%s", centerDocId.String()),
			},
		}
		resp = handler.HandleGet(respWriter, req)
		if resp == nil {
			return fmt.Errorf("5.API request failed: resp is nil")
		}

		// Write data
		f.SetCellValue("Results", fmt.Sprintf("F%d", rowNum), resp[types.MongoDuration])
		f.SetCellValue("Results", fmt.Sprintf("G%d", rowNum), resp[types.BlockchainDuration])
		f.SetCellValue("Results", fmt.Sprintf("H%d", rowNum), resp[types.TotalDuration])

		// API
		req = &http.Request{
			Method: "GET",
			URL: &url.URL{
				Path: fmt.Sprintf(apiURL, "/%s", lastDocId.String()),
			},
		}
		resp = handler.HandleGet(respWriter, req)
		if resp == nil {
			return fmt.Errorf("6.API request failed: resp is nil")
		}

		// Write data
		f.SetCellValue("Results", fmt.Sprintf("I%d", rowNum), resp[types.MongoDuration])
		f.SetCellValue("Results", fmt.Sprintf("J%d", rowNum), resp[types.BlockchainDuration])
		f.SetCellValue("Results", fmt.Sprintf("K%d", rowNum), resp[types.TotalDuration])

		rowNum++
	}

	// Save file
	if err := f.SaveAs(fileName); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}
