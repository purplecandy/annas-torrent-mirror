package database

import (
	"annas-mirror/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dbUri := os.Getenv("DATABASE_URI")

	db, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Unable to connect to the database \n", err)
	}

	log.Println("Connected to database!")
	db.Logger = logger.Default.LogMode(logger.Silent)

	log.Println("Performing migrations")
	db.AutoMigrate(&models.Torrent{})

	DB = db
}
