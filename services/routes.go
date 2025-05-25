package services

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	blockchain "github.com/sikozonpc/ecom/blockchain"
	"github.com/sikozonpc/ecom/database"
	"github.com/sikozonpc/ecom/types"
	"github.com/sikozonpc/ecom/utils"
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

func (h *Handler) TestRoutes(router *mux.Router) {
	router.HandleFunc("/test", h.handleTest).Methods("GET")
}

func (h *Handler) DataRoutes(router *mux.Router) {
	router.HandleFunc("/send", h.handleSend).Methods("POST")
	router.HandleFunc("/get/{dataId}", h.handleGet).Methods("GET")
}

func (h *Handler) handleSend(w http.ResponseWriter, r *http.Request) {
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

	if err := h.database.Add(payload.DataId, payload.Data, payload.DeviceId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add data to database: %v", err))
		return
	}

	if err := h.blockchain.Send(payload.DataId, hash, payload.DeviceId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add data hash to blockchain: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleTest(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
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

	doc, err := h.database.Get(uuid)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get data from database: %v", err))
		return
	}

	if doc == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("data not found for id: %s", dataId))
		return
	}

	// data, ok := doc["data"]
	// if !ok {
	// 	utils.WriteError(w, http.StatusNotFound, fmt.Errorf("data field not found for id: %s", dataId))
	// 	return
	// }

	hash, err := utils.CalculateHash(doc)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to calculate hash: %v", err))
		return
	}

	success, err := h.blockchain.VerifyHash(uuid, hash)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to verify hash on blockchain: %v", err))
		return
	}
	if !success {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("hash not found on blockchain for dataId: %s", dataId))
		return
	}

	utils.WriteJSON(w, http.StatusOK, doc)
}
