package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	blockchain "github.com/mbabinski218/BaaS-IoT-platform/blockchain"
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

	mongoDuration, err := h.database.Add(payload.DataId, payload.Data, payload.DeviceId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add data to database: %v", err))
		return
	}

	blockchainDuration, blockchainSendDuration, blockchainMinedDuration, err := h.blockchain.Send(payload.DataId, hash, payload.DeviceId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add data hash to blockchain: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

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

	doc, mongoDuration, err := h.database.Get(uuid)
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

	success, blockchainDuration, err := h.blockchain.VerifyHash(uuid, hash)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("blockchain error: %v", err))
		return
	}
	if !success {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("blockchain - hash not found or invalid for dataId: %s", dataId))
		return
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
	fromTimestamp, err := time.Parse(time.RFC3339, from)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid 'from' timestamp format: %v", err))
		return
	}
	toTimestamp, err := time.Parse(time.RFC3339, to)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid 'to' timestamp format: %v", err))
		return
	}

	docs, mongoDuration, err := h.database.GetFromTo(fromTimestamp, toTimestamp)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get data from database: %v", err))
		return
	}

	if len(docs) == 0 {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("no data found between %s and %s", from, to))
		return
	}

	success, blockchainDuration, err := h.blockchain.VerifyHashes(docs)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("blockchain error: %v", err))
		return
	}
	if !success {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("blockchain - failed to verify hashes for documents in the specified range"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, docs)

	duration := time.Since(start)
	fmt.Println("-------- Data retrieved successfully --------")
	fmt.Println("MongoDB duration:", mongoDuration)
	fmt.Println("Blockchain duration:", blockchainDuration)
	fmt.Println("Total duration:", duration)
}
