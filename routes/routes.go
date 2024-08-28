package routes

import (
	"annas-mirror/database"
	"annas-mirror/models"

	"github.com/gofiber/fiber/v2"
)

type AllTorrentsGroupResult struct {
	GroupName string
	Count     uint
	Bytes     string
	Files     string
}

func AllTorrentGroups(c *fiber.Ctx) error {

	var results []AllTorrentsGroupResult

	err := database.DB.Model(&models.Torrent{}).Select("group_name, COUNT(*) as count, sum(data_size)::TEXT AS bytes, sum(num_files)::TEXT AS files").Group("group_name").Order("count asc").Scan(&results).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err})
	}

	var healthySeedsCount int64
	var poorSeedsCount int64
	var riskySeedsCount int64

	database.DB.Model(&models.Torrent{}).Where("seeders > ?", 10).Where("embargo = ?", false).Count(&healthySeedsCount)
	database.DB.Model(&models.Torrent{}).Where("seeders > ?", 4).Where("seeders < ?", 10).Where("embargo = ?", false).Count(&poorSeedsCount)
	database.DB.Model(&models.Torrent{}).Where("seeders < ?", 4).Where("embargo = ?", false).Count(&riskySeedsCount)

	return c.Render("index", fiber.Map{"Result": results, "healthySeedsCount": healthySeedsCount, "poorSeedsCount": poorSeedsCount, "riskySeedsCount": riskySeedsCount})

}

func TorrentGroupData(c *fiber.Ctx) error {

	group := c.Params("group")

	var results []map[string]interface{}

	err := database.DB.Model(&models.Torrent{}).Where("group_name = ?", group).Find(&results).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err})
	}

	var healthySeedsCount string
	var poorSeedsCount string
	var riskySeedsCount string

	database.DB.Model(&models.Torrent{}).Where("seeders > ?", 10).Where("embargo = ?", false).Select("sum(data_size)::TEXT AS bytes").Find(&healthySeedsCount)
	database.DB.Model(&models.Torrent{}).Where("seeders > ?", 4).Where("seeders < ?", 10).Where("embargo = ?", false).Select("sum(data_size)::TEXT AS bytes").Find(&riskySeedsCount)
	database.DB.Model(&models.Torrent{}).Where("seeders < ?", 4).Where("embargo = ?", false).Select("sum(data_size)::TEXT AS bytes").Find(&poorSeedsCount)

	return c.Render("torrent_group", fiber.Map{"Result": results, "Group": group, "healthySeedsCount": healthySeedsCount, "poorSeedsCount": poorSeedsCount, "riskySeedsCount": riskySeedsCount})

}
