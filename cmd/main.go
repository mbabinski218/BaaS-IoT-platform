package main

import (
	"fmt"
	"log"

	client "github.com/sikozonpc/ecom/blockchain"
	"github.com/sikozonpc/ecom/cmd/api"
	"github.com/sikozonpc/ecom/configs"
)

func main() {
	client, err := client.NewEthClient("testnet")
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), client)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
