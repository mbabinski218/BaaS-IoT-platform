package device

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	blockchain "github.com/sikozonpc/ecom/blockchain"
	"github.com/sikozonpc/ecom/types"
	"github.com/sikozonpc/ecom/utils"
)

type Handler struct {
	client *blockchain.Client
}

func NewHandler(bcClient *blockchain.Client) *Handler {
	return &Handler{client: bcClient}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var device types.RegisterDevicePayload
	if err := utils.ParseJSON(r, &device); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(device); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	// DODAWANIE DO BLOCKCHAIN NOWEGO URZADZENIA!

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) SendData(w http.ResponseWriter, r *http.Request) {
	var data types.IotData
	if err := utils.ParseJSON(r, &data); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(data); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	h.client.Send(data)

	utils.WriteJSON(w, http.StatusOK, nil)
}
