package actions

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
)

const (
	TypeSyncTorrents = "sync:torrents" // sync torrent from anna's archive
)

func DispatchSyncTorrents() *asynq.Task {
	return asynq.NewTask(TypeSyncTorrents, nil, asynq.MaxRetry(2))
}

func SyncTorrents(ctx context.Context, t *asynq.Task) error {

	log.Println("Syncing torrents from anaa")
	return nil
}
