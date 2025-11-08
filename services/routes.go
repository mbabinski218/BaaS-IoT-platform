package services

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	blockchain "github.com/mbabinski218/BaaS-IoT-platform/blockchain"
	"github.com/mbabinski218/BaaS-IoT-platform/configs"
	"github.com/mbabinski218/BaaS-IoT-platform/database"
	"github.com/mbabinski218/BaaS-IoT-platform/types"
	"github.com/mbabinski218/BaaS-IoT-platform/utils"
	"github.com/xuri/excelize/v2"
)

type Handler struct {
	blockchain *blockchain.Client
	database   *database.Client
	docCount   uint64
}

func NewHandler(bc *blockchain.Client, db *database.Client) *Handler {
	return &Handler{
		blockchain: bc,
		database:   db,
	}
}

func (h *Handler) DataRoutes(router *mux.Router) {
	router.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		_ = h.HandleSend(w, r)
	}).Methods("POST")

	router.HandleFunc("/get/{dataId}", func(w http.ResponseWriter, r *http.Request) {
		_ = h.HandleGet(w, r)
	}).Methods("GET")

	router.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		_ = h.HandleGetFromTo(w, r)
	}).Methods("GET")

	router.HandleFunc("/blocknumber", h.HandleGetBlockNumber).Methods("GET")
}

func (h *Handler) HandleSend(w http.ResponseWriter, r *http.Request) map[string]any {
	start := time.Now()

	var payload types.NewDataPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return nil
	}

	bsonDoc := utils.MapToBSON(payload)

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return nil
	}

	hash, err := utils.StringToBytes32(payload.Hash)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return nil
	}

	createdId, mongoDuration, err := h.database.Add(bsonDoc)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add data with id: %s to database: %w", payload.DataId, err))
		return nil
	}

	blockchainDuration, blockchainSendDuration, err := h.blockchain.Send(payload.DataId, hash, payload.DeviceId)
	if err != nil {
		h.database.Delete(createdId)

		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add data with id: %s hash to blockchain: %v", payload.DataId, err))
		return nil
	}

	utils.WriteJSON(w, http.StatusCreated, createdId)

	duration := time.Since(start)

	fmt.Println("-------- Data sent successfully --------")
	fmt.Println("MongoDB duration:", mongoDuration)
	fmt.Println("Blockchain duration:", blockchainDuration)
	fmt.Println("Blockchain send duration:", blockchainSendDuration)
	fmt.Println("Total duration:", duration)

	if configs.Envs.BlockchainLogSendToFile {
		LogSendToFile(createdId, mongoDuration, blockchainDuration, blockchainSendDuration)
	}

	result := make(map[string]any)
	result[types.MongoDuration] = mongoDuration.String()
	result[types.BlockchainDuration] = blockchainDuration.String()
	result[types.TotalDuration] = duration.String()
	utils.WriteJSON(w, http.StatusOK, result)

	h.docCount++
	fmt.Println("Document count:", h.docCount)
	if slices.Contains(configs.Envs.BlockchainMaxDocuments, h.docCount) {
		if err := h.Audit(h.docCount); err != nil {
			log.Printf("failed to create audit report: %v", err)
		}
	}

	return result

	// h.Audit(0)
	// return make(map[string]any)
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) map[string]any {
	start := time.Now()

	vars := mux.Vars(r)
	dataId := vars["dataId"]
	if dataId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("dataId is required"))
		return nil
	}

	uuid, err := uuid.Parse(dataId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid dataId format: %v", err))
		return nil
	}

	doc, proof, mongoDuration, err := h.database.Get(uuid)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get data from database: %v", err))
		return nil
	}

	if doc == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("data not found for id: %s", dataId))
		return nil
	}

	hash, err := utils.CalculateHash(doc)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to calculate hash: %v", err))
		return nil
	}

	var timestamp time.Time
	if proof != nil {
		timestamp, _ = time.Parse(types.TimeLayout, doc["timestamp"].(string))
	}

	success, blockchainDuration, err := h.blockchain.VerifyHash(uuid, hash, timestamp, proof)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("blockchain error: %v", err))
		return nil
	}
	if !success {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("blockchain - hash not found or invalid for dataId: %s", dataId))
		return nil
	}

	if configs.Envs.BlockchainMode == types.BCPeriodicBatchCheck && proof == nil {
		doc["verified"] = false
	}

	duration := time.Since(start)
	fmt.Println("-------- Data retrieved successfully --------")
	fmt.Println("MongoDB duration:", mongoDuration)
	fmt.Println("Blockchain duration:", blockchainDuration)
	fmt.Println("Total duration:", duration)

	result := make(map[string]any)
	result["data"] = doc
	result[types.MongoDuration] = mongoDuration.String()
	result[types.BlockchainDuration] = blockchainDuration.String()
	result[types.TotalDuration] = duration.String()
	utils.WriteJSON(w, http.StatusOK, result)
	return result
}

func (h *Handler) HandleGetFromTo(w http.ResponseWriter, r *http.Request) map[string]any {
	start := time.Now()

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("both 'from' and 'to' query parameters are required"))
		return nil
	}
	fromTimestamp, err := time.Parse(types.TimeLayout, from)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid 'from' timestamp format: %v", err))
		return nil
	}
	toTimestamp, err := time.Parse(types.TimeLayout, to)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid 'to' timestamp format: %v", err))
		return nil
	}

	if fromTimestamp.IsZero() || toTimestamp.IsZero() {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("timestamps cannot be zero"))
		return nil
	}
	if fromTimestamp.After(toTimestamp) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("'from' timestamp cannot be after 'to' timestamp"))
		return nil
	}
	if toTimestamp.Before(fromTimestamp) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("'to' timestamp cannot be before 'from' timestamp"))
		return nil
	}

	fixedFromTimestamp, fixedToTimestamp, err := utils.FixTimestamps(fromTimestamp, toTimestamp, time.Duration(configs.Envs.BlockchainBatchInterval))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to fix timestamps: %v", err))
		return nil
	}

	docs, mongoDuration, err := h.database.GetFromTo(fixedFromTimestamp, fixedToTimestamp)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get data from database: %v", err))
		return nil
	}

	if len(docs) == 0 {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("no data found between %s and %s", from, to))
		return nil
	}

	success, blockchainDuration, err := h.blockchain.VerifyHashes(docs, fixedFromTimestamp, fixedToTimestamp, true)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("blockchain error: %v", err))
		return nil
	}
	if !success {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("blockchain - failed to verify hashes for documents in the specified range"))
		return nil
	}

	var fixedDocs []types.DocData
	// If the blockchain verified the hashes, we can match with the original timestamps
	for _, doc := range docs {
		timestampStr, ok := doc.Data["timestamp"].(string)
		if !ok {
			continue
		}

		timestamp, err := time.Parse(types.TimeLayout, timestampStr)
		if err != nil {
			continue
		}

		if timestamp.Before(fromTimestamp) || timestamp.After(toTimestamp) {
			continue
		}

		fixedDocs = append(fixedDocs, doc)
	}

	duration := time.Since(start)
	fmt.Println("-------- Data retrieved successfully --------")
	fmt.Println("MongoDB duration:", mongoDuration)
	fmt.Println("Blockchain duration:", blockchainDuration)
	fmt.Println("Total duration:", duration)

	result := make(map[string]any)
	result["data"] = fixedDocs
	result[types.MongoDuration] = mongoDuration.String()
	result[types.BlockchainDuration] = blockchainDuration.String()
	result[types.TotalDuration] = duration.String()
	result[types.Missed] = len(docs) - len(fixedDocs)
	utils.WriteJSON(w, http.StatusOK, result)
	return result
}

func (h *Handler) HandleGetBlockNumber(w http.ResponseWriter, r *http.Request) {
	blockNumber, err := h.blockchain.GetBlockNumber()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("blockchain error: %v", err))
	} else {
		utils.WriteJSON(w, http.StatusOK, blockNumber)
	}
}

func (h *Handler) Audit(count uint64) error {
	if err := h.blockchain.StopMining(); err != nil {
		log.Println("Error stopping mining:", err)
	}
	utils.PauseSimulatorFile()

	apiURL := fmt.Sprintf("http://%s/get", configs.Envs.PublicHost)

	fileName := fmt.Sprintf("D:\\Studia_mgr\\Praca_magisterska\\Results\\audit_results_%s_%d.xlsx", configs.Envs.BlockchainMode.String(), count)
	var f *excelize.File
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		f = excelize.NewFile()
		f.SetSheetName("Sheet1", "Results")
		f.SetCellValue("Results", "A1", "Timestamp")
		f.SetCellValue("Results", "B1", "Db duration")
		f.SetCellValue("Results", "C1", "Blockchain duration")
		f.SetCellValue("Results", "D1", "Total duration")
		f.SetCellValue("Results", "E1", "Missed")
	} else {
		return fmt.Errorf("file %s already exists", fileName)
	}

	for range configs.Envs.BlockchainCheckpointCallRepeats {
		// Find next empty row
		rows, err := f.GetRows("Results")
		if err != nil {
			return fmt.Errorf("failed to get rows: %w", err)
		}
		rowNum := len(rows) + 1

		f.SetCellValue("Results", fmt.Sprintf("A%d", rowNum), time.Now().Format(types.TimeLayout))

		req := &http.Request{
			Method: "GET",
			URL: &url.URL{
				Path:     apiURL,
				RawQuery: "from=2025-01-01T00:00:00.000000&to=2026-01-01T00:00:00.000000",
			},
		}
		resp := h.HandleGetFromTo(nil, req)

		f.SetCellValue("Results", fmt.Sprintf("B%d", rowNum), resp[types.MongoDuration])
		f.SetCellValue("Results", fmt.Sprintf("C%d", rowNum), resp[types.BlockchainDuration])
		f.SetCellValue("Results", fmt.Sprintf("D%d", rowNum), resp[types.TotalDuration])
		f.SetCellValue("Results", fmt.Sprintf("E%d", rowNum), resp[types.Missed])
	}

	if err := f.SaveAs(fileName); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	if err := h.blockchain.StartMining(); err != nil {
		log.Println("Error starting mining:", err)
	}
	utils.ResumeSimulatorFile()

	return nil
}

func LogSendToFile(createdId uuid.UUID, mongoDuration, blockchainDuration, blockchainSendDuration time.Duration) {
	fileName := fmt.Sprintf("D:\\Studia_mgr\\Praca_magisterska\\Results\\send_results_%s.xlsx", configs.Envs.BlockchainMode.String())
	var f *excelize.File
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		f = excelize.NewFile()
		f.SetSheetName("Sheet1", "Results")
		f.SetCellValue("Results", "A1", "Timestamp")
		f.SetCellValue("Results", "B1", "Id")
		f.SetCellValue("Results", "C1", "Db duration")
		f.SetCellValue("Results", "D1", "Blockchain duration")
		f.SetCellValue("Results", "E1", "Blockchain send duration")
	} else {
		f, _ = excelize.OpenFile(fileName)
	}

	// Find next empty row
	rows, _ := f.GetRows("Results")
	rowNum := len(rows) + 1

	f.SetCellValue("Results", fmt.Sprintf("A%d", rowNum), time.Now().Format(types.TimeLayout))
	f.SetCellValue("Results", fmt.Sprintf("B%d", rowNum), createdId.String())
	f.SetCellValue("Results", fmt.Sprintf("C%d", rowNum), mongoDuration.String())
	f.SetCellValue("Results", fmt.Sprintf("D%d", rowNum), blockchainDuration.String())
	f.SetCellValue("Results", fmt.Sprintf("E%d", rowNum), blockchainSendDuration.String())

	f.SaveAs(fileName)
}
