package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	blockchain "github.com/sikozonpc/ecom/blockchain"
	device "github.com/sikozonpc/ecom/services/device"
)

type APIServer struct {
	addr   string
	client *blockchain.Client
}

func NewAPIServer(addr string, client *blockchain.Client) *APIServer {
	return &APIServer{
		addr:   addr,
		client: client,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	deviceHandler := device.NewHandler(s.client)
	deviceHandler.RegisterRoutes(subrouter)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
