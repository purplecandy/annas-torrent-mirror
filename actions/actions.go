package actions

import (
	"annas-mirror/database"
	"annas-mirror/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	TypeSyncTorrents = "sync:torrents" // sync torrent from anna's archive
)

type TorrentData struct {
	URL                 string `json:"url"`
	TopLevelGroupName   string `json:"top_level_group_name"`
	GroupName           string `json:"group_name"`
	DisplayName         string `json:"display_name"`
	AddedToTorrentsList string `json:"added_to_torrents_list_at"`
	IsMetadata          bool   `json:"is_metadata"`
	BTIH                string `json:"btih"`
	MagnetLink          string `json:"magnet_link"`
	TorrentSize         uint   `json:"torrent_size"`
	NumFiles            uint   `json:"num_files"`
	DataSize            int64  `json:"data_size"`
	AACurrentlySeeding  bool   `json:"aa_currently_seeding"`
	Obsolete            bool   `json:"obsolete"`
	Embargo             bool   `json:"embargo"`
	Seeders             uint   `json:"seeders"`
	Leechers            uint   `json:"leechers"`
	Completed           uint   `json:"completed"`
	StatsScrapedAt      string `json:"stats_scraped_at"`
	PartiallyBroken     bool   `json:"partially_broken"`
}

func DispatchSyncTorrents() *asynq.Task {
	return asynq.NewTask(TypeSyncTorrents, nil, asynq.MaxRetry(1))
}

func SyncTorrents(ctx context.Context, t *asynq.Task) error {
	log.Printf("Received job - %v", t.ResultWriter().TaskID())

	resp, err := http.Get(os.Getenv("TORRENTS_URL"))

	if err != nil {
		return fmt.Errorf("error trying to fetch torrent from source:  %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("error trying to read response body: %v", err)
	}

	var jsonData []TorrentData

	if err := json.Unmarshal(body, &jsonData); err != nil {
		return fmt.Errorf("error parsing json payload: %v", err)
	}

	var payload []models.Torrent
	count := 0

	log.Printf("Total torrents: %v", len(jsonData))

	for i, data := range jsonData {
		addedAt, err := time.Parse("2006-01-02", data.AddedToTorrentsList)

		if err != nil {
			return fmt.Errorf("enable to parse Date for item: %v", err)
		}

		scrapedAt, err := time.Parse("2006-01-02T15:04:05", data.StatsScrapedAt)

		if err != nil {
			return fmt.Errorf("enable to parse Date for item: %v", err)
		}

		payload = append(payload, models.Torrent{
			Url:               data.URL,
			Btih:              data.BTIH,
			TopLevelGroupName: data.TopLevelGroupName,
			GroupName:         data.GroupName,
			DisplayName:       data.DisplayName,
			AddedAt:           addedAt,
			IsMetadata:        data.IsMetadata,
			MagnetLink:        data.MagnetLink,
			TorrentSize:       data.TorrentSize,
			NumFiles:          data.NumFiles,
			DataSize:          uint(data.DataSize),
			CurrentlySeeding:  data.AACurrentlySeeding,
			Obsolete:          data.Obsolete,
			Embargo:           data.Embargo,
			Seeders:           data.Seeders,
			Leechers:          data.Leechers,
			Completed:         data.Completed,
			PartiallyBroken:   data.PartiallyBroken,
			ScrapedAt:         scrapedAt,
		})

		// We can't save all the records at once, so saving them in chunk of 1000
		if i-count == 1000 || len(jsonData) == i {
			txError := database.DB.Transaction(func(tx *gorm.DB) error {
				err := tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "btih"}},
					UpdateAll: true,
				}).Create(payload).Error

				return err
			})

			if txError != nil {
				fmt.Printf("error - %v", txError.Error())
				return fmt.Errorf("error while creating records: %v", txError)
			}

			count = i
			payload = nil
			log.Printf("Saved chunk: %v", count/1000)
		}

	}

	log.Println("Successfully synced torrents from source!")

	return nil
}
