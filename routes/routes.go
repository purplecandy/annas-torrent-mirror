package routes

import (
	"annas-mirror/database"
	"annas-mirror/models"

	"github.com/gofiber/fiber/v2"
)

type AllTorrentsGroupResult struct {
	GroupName string
	Count     uint
}

func AllTorrentGroups(c *fiber.Ctx) error {

	var results []AllTorrentsGroupResult

	err := database.DB.Model(&models.Torrent{}).Select("group_name, COUNT(*) as count").Group("group_name").Scan(&results).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err})
	}

	var healthySeedsCount int64
	var poorSeedsCount int64
	var riskySeedsCount int64

	database.DB.Model(&models.Torrent{}).Where("seeders > ?", 10).Count(&healthySeedsCount)
	database.DB.Model(&models.Torrent{}).Where("seeders > ?", 4).Where("seeders < ?", 10).Count(&poorSeedsCount)
	database.DB.Model(&models.Torrent{}).Where("seeders < ?", 4).Count(&riskySeedsCount)

	return c.Render("index", fiber.Map{"Result": results, "healthySeedsCount": healthySeedsCount, "poorSeedsCount": poorSeedsCount, "riskySeedsCount": riskySeedsCount})

}
