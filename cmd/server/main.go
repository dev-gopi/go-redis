package main

import (
	"log"

	"github.com/dev-gopi/go-redis/internal/network"
)

func main() {
	server := network.NewServer(":6379")

	log.Println("Redis clone running on :6379")

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
