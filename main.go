package main

import (
	"log"

	"github.com/raphaelbertoldo/scraping-rentals-br/api"
)

func main() {
	server := api.NewServer()
	if err := server.Router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
