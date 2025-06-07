package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	blockchain "github.com/mbabinski218/BaaS-IoT-platform/blockchain"
	"github.com/mbabinski218/BaaS-IoT-platform/database"
	device "github.com/mbabinski218/BaaS-IoT-platform/services"
)

type APIServer struct {
	address    string
	blockchain *blockchain.Client
	database   *database.Client
}

func NewAPIServer(addr string, bc *blockchain.Client, db *database.Client) *APIServer {
	return &APIServer{
		address:    addr,
		blockchain: bc,
		database:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	deviceHandler := device.NewHandler(s.blockchain, s.database)
	deviceHandler.DataRoutes(subrouter)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.Println("Listening on", "http://"+s.address)

	return http.ListenAndServe(s.address, router)
}
