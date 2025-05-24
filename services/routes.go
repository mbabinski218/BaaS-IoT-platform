package services

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
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

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) SendRoutes(router *mux.Router) {
	router.HandleFunc("/send", h.handleSend).Methods("POST")
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

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// 	var device types.RegisterDevicePayload
	// 	if err := utils.ParseJSON(r, &device); err != nil {
	// 		utils.WriteError(w, http.StatusBadRequest, err)
	// 		return
	// 	}

	// 	if err := utils.Validate.Struct(device); err != nil {
	// 		errors := err.(validator.ValidationErrors)
	// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
	// 		return
	// 	}

	// 	// DODAWANIE DO BLOCKCHAIN NOWEGO URZADZENIA!

	// 	utils.WriteJSON(w, http.StatusCreated, nil)
	// }

	// func (h *Handler) SendData(w http.ResponseWriter, r *http.Request) {
	// 	var data types.IotData
	// 	if err := utils.ParseJSON(r, &data); err != nil {
	// 		utils.WriteError(w, http.StatusBadRequest, err)
	// 		return
	// 	}

	// 	if err := utils.Validate.Struct(data); err != nil {
	// 		errors := err.(validator.ValidationErrors)
	// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
	// 		return
	// 	}

	// 	// h.client.Send(data)

	// 	utils.WriteJSON(w, http.StatusOK, nil)
}
