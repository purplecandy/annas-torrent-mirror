package cache

import (
	"os"

	"github.com/hibiken/asynq"
)

var Dispatcher *asynq.Client

func ConnectCache() {
	Dispatcher = asynq.NewClient(asynq.RedisClientOpt{
		Addr: os.Getenv("REDIS_URI"),
	})
}
