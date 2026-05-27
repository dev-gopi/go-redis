package main

import (
	"log"
	"os"
	"time"

	"github.com/dev-gopi/go-redis/internal/auth"
	"github.com/dev-gopi/go-redis/internal/logger"
	"github.com/dev-gopi/go-redis/internal/network"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/snapshot"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func main() {

	logger.Init()
	auth.LoadFromEnv()

	logger.InfoLogger.Println(
		"Starting Redis Clone Server...",
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("REDIS_PORT")
	}
	if port == "" {
		port = "6379"
	}

	server := network.NewServer(":" + port)

	storage.StartTTLWorker()

	err := snapshot.Load("data/dump.rdb")
	if err != nil {
		panic(err)
	}

	err = wal.Init("data/wal.log")
	if err != nil {
		panic(err)
	}

	err = aof.Init("data/appendonly.aof")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(time.Minute)
			if err := snapshot.Save("data/dump.rdb"); err != nil {
				logger.ErrorLogger.Printf("Snapshot save failed: %v", err)
				continue
			}

			if err := aof.Reset("data/appendonly.aof"); err != nil {
				logger.ErrorLogger.Printf("AOF reset failed after snapshot: %v", err)
			}
		}
	}()

	err = aof.Replay("data/appendonly.aof")
	if err != nil {
		panic(err)
	}

	aof.StartAutoRotate(
		"data/appendonly.aof",
		10*1024*1024,
	)

	logger.InfoLogger.Printf("Redis clone running on :%s", port)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
