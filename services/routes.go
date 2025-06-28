package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	blockchain "github.com/mbabinski218/BaaS-IoT-platform/blockchain"
	"github.com/mbabinski218/BaaS-IoT-platform/configs"
	"github.com/mbabinski218/BaaS-IoT-platform/database"
	"github.com/mbabinski218/BaaS-IoT-platform/types"
	"github.com/mbabinski218/BaaS-IoT-platform/utils"
)

type Handler struct {
	blockchain *blockchain.Client
	database   *database.Client
}

func NewHandler(bc *blockchain.Client, db *database.Client) *Handler {
	return &Handler{
		blockchain: bc,
		database:   db,
	}
}

func (h *Handler) DataRoutes(router *mux.Router) {
	router.HandleFunc("/send", h.handleSend).Methods("POST")
	router.HandleFunc("/get/{dataId}", h.handleGet).Methods("GET")
	router.HandleFunc("/get", h.handleGetFromTo).Methods("GET")
	router.HandleFunc("/blocknumber", h.handleGetBlockNumber).Methods("GET")
}

func (h *Handler) handleSend(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var payload types.NewDataPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	hash, err := utils.StringToBytes32(payload.Hash)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	createdId, mongoDuration, err := h.database.Add(payload.DataId, payload.Data, payload.DeviceId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add data with id: %s to database: %w", payload.DataId, err))
		return
	}

	blockchainDuration, blockchainSendDuration, blockchainMinedDuration, err := h.blockchain.Send(payload.DataId, hash, payload.DeviceId)
	if err != nil {
		h.database.Delete(createdId)

		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add data with id: %s hash to blockchain: %v", payload.DataId, err))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, createdId)

	duration := time.Since(start)

	fmt.Println("-------- Data sent successfully --------")
	fmt.Println("MongoDB duration:", mongoDuration)
	fmt.Println("Blockchain duration:", blockchainDuration)
	fmt.Println("Blockchain send duration:", blockchainSendDuration)
	fmt.Println("Blockchain mined duration:", blockchainMinedDuration)
	fmt.Println("Total duration:", duration)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	vars := mux.Vars(r)
	dataId := vars["dataId"]
	if dataId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("dataId is required"))
		return
	}

	uuid, err := uuid.Parse(dataId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid dataId format: %v", err))
		return
	}

	doc, proof, mongoDuration, err := h.database.Get(uuid)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get data from database: %v", err))
		return
	}

	if doc == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("data not found for id: %s", dataId))
		return
	}

	hash, err := utils.CalculateHash(doc)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to calculate hash: %v", err))
		return
	}

	var timestamp time.Time
	if proof != nil {
		timestamp, _ = time.Parse(types.TimeLayout, doc["timestamp"].(string))
	}

	success, blockchainDuration, err := h.blockchain.VerifyHash(uuid, hash, timestamp, proof)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("blockchain error: %v", err))
		return
	}
	if !success {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("blockchain - hash not found or invalid for dataId: %s", dataId))
		return
	}

	if configs.Envs.BlockchainMode == types.BCBatchCheck && proof == nil {
		doc["verified"] = false
	}

	utils.WriteJSON(w, http.StatusOK, doc)

	duration := time.Since(start)
	fmt.Println("-------- Data retrieved successfully --------")
	fmt.Println("MongoDB duration:", mongoDuration)
	fmt.Println("Blockchain duration:", blockchainDuration)
	fmt.Println("Total duration:", duration)
}

func (h *Handler) handleGetFromTo(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("both 'from' and 'to' query parameters are required"))
		return
	}
	fromTimestamp, err := time.Parse(types.TimeLayout, from)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid 'from' timestamp format: %v", err))
		return
	}
	toTimestamp, err := time.Parse(types.TimeLayout, to)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid 'to' timestamp format: %v", err))
		return
	}

	if fromTimestamp.IsZero() || toTimestamp.IsZero() {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("timestamps cannot be zero"))
		return
	}
	if fromTimestamp.After(toTimestamp) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("'from' timestamp cannot be after 'to' timestamp"))
		return
	}
	if toTimestamp.Before(fromTimestamp) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("'to' timestamp cannot be before 'from' timestamp"))
		return
	}

	fixedFromTimestamp, fixedToTimestamp, err := utils.FixTimestamps(fromTimestamp, toTimestamp, time.Duration(configs.Envs.BlockchainBatchInterval))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to fix timestamps: %v", err))
		return
	}

	docs, mongoDuration, err := h.database.GetFromTo(fixedFromTimestamp, fixedToTimestamp)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get data from database: %v", err))
		return
	}

	if len(docs) == 0 {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("no data found between %s and %s", from, to))
		return
	}

	success, blockchainDuration, err := h.blockchain.VerifyHashes(docs, fixedFromTimestamp, fixedToTimestamp)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("blockchain error: %v", err))
		return
	}
	if !success {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("blockchain - failed to verify hashes for documents in the specified range"))
		return
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

	utils.WriteJSON(w, http.StatusOK, fixedDocs)

	duration := time.Since(start)
	fmt.Println("-------- Data retrieved successfully --------")
	fmt.Println("MongoDB duration:", mongoDuration)
	fmt.Println("Blockchain duration:", blockchainDuration)
	fmt.Println("Total duration:", duration)
}

func (h *Handler) handleGetBlockNumber(w http.ResponseWriter, r *http.Request) {
	blockNumber, err := h.blockchain.GetBlockNumber()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("blockchain error: %v", err))
	} else {
		utils.WriteJSON(w, http.StatusOK, blockNumber)
	}
}
