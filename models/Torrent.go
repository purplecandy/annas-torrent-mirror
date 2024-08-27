package models

import (
	"time"

	"gorm.io/gorm"
)

type Torrent struct {
	ID                uint `gorm:"primaryKey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	Url               string
	TopLevelGroupName string
	DisplayName       string
	AddedAt           time.Time
	IsMetadata        bool
	Btih              string `gorm:"uniqueIndex"`
	MagnetLink        string `gorm:"index"`
	TorrentSize       uint
	NumFiles          uint
	DataSize          uint
	CurrentlySeeding  bool
	Obsolete          bool
	Embargo           bool
	Seeders           uint
	Leechers          uint
	Completed         uint
	ScrapedAt         time.Time
	PartiallyBroken   bool
}

// {
// 	"url": "https://annas-archive.org/dyn/small_file/torrents/managed_by_aa/annas_archive_meta__aacid/annas_archive_meta__aacid__duxiu_records__20240130T000000Z--20240209T000000Z.jsonl.zst.torrent",
// 	"top_level_group_name": "managed_by_aa",
// 	"group_name": "duxiu",
// 	"display_name": "annas_archive_meta__aacid__duxiu_records__20240130T000000Z--20240209T000000Z.jsonl.zst.torrent",
// 	"added_to_torrents_list_at": "2024-02-21",
// 	"is_metadata": true,
// 	"btih": "47ac1ce69f07fb667fe364abf2961ffc90d81a64",
// 	"magnet_link": "magnet:?xt=urn:btih:47ac1ce69f07fb667fe364abf2961ffc90d81a64&dn=annas_archive_meta__aacid__duxiu_records__20240130T000000Z--20240209T000000Z.jsonl.zst.torrent&tr=udp://tracker.opentrackr.org:1337/announce",
// 	"torrent_size": 3488,
// 	"num_files": 1,
// 	"data_size": 10015727486,
// 	"aa_currently_seeding": true,
// 	"obsolete": true,
// 	"embargo": false,
// 	"seeders": 8,
// 	"leechers": 4,
// 	"completed": 37,
// 	"stats_scraped_at": "2024-08-27T12:01:15",
// 	"partially_broken": false,
// 	"random": "jFBSNach3ckorLVYozttcX"
// }
