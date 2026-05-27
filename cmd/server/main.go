package main

import (
	"log"
	"time"

	"github.com/dev-gopi/go-redis/internal/logger"
	"github.com/dev-gopi/go-redis/internal/network"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/snapshot"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func main() {

	logger.Init()

	logger.InfoLogger.Println(
		"Starting Redis Clone Server...",
	)

	server := network.NewServer(":6379")

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
			_ = snapshot.Save("data/dump.rdb")
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

	logger.InfoLogger.Println("Redis clone running on :6379")

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
