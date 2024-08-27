package main

import (
	"annas-mirror/actions"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Println("Warning: Unable to load .env file")
	}

	server := asynq.NewServer(asynq.RedisClientOpt{Addr: os.Getenv("REDIS_URI")}, asynq.Config{
		Concurrency: 2,
	})

	mux := asynq.NewServeMux()

	mux.HandleFunc(actions.TypeSyncTorrents, actions.SyncTorrents)

	if err := server.Run(mux); err != nil {
		log.Fatalf("Something went wrong with worker: %v", err)
	}

}
