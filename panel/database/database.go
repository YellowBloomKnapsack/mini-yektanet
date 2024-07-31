package database

import (
	"YellowBloomKnapsack/mini-yektanet/common/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initPublishers() error {
	var count int64
	DB.Model(&models.Publisher{}).Count(&count)

	if count != 0 {
		return nil
	}

	publishers := []models.Publisher{
		{Username: "varzesh3"},
		{Username: "digikala"},
		{Username: "zoomit"},
		{Username: "sheypoor"},
		{Username: "filimo"},
	}

	if err := DB.Create(&publishers).Error; err != nil {
		return err
	}

	return nil
}

func InitDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB.AutoMigrate(&models.Publisher{})
	DB.AutoMigrate(&models.Ad{})
	DB.AutoMigrate(&models.AdsInteraction{})
	DB.AutoMigrate(&models.Advertiser{})
	DB.AutoMigrate(&models.Transaction{})
	DB.AutoMigrate(&models.Keyword{})
	err = initPublishers()
	if err != nil {
		log.Fatal(err)
	}
}

func InitTestDB() {
	dbPath := os.Getenv("TEST_DB_PATH")
	if dbPath == "" {
		dbPath = "test.db"
	}

	if _, err := os.Stat(dbPath); err == nil {
		os.Remove(dbPath)
	}
	// Create the SQLite database file if it doesn't exist
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		_, err = os.Create(dbPath)
		if err != nil {
			log.Fatal("Failed to create test database file:", err)
		}
	}

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	DB.AutoMigrate(&models.Publisher{})
	DB.AutoMigrate(&models.Ad{})
	DB.AutoMigrate(&models.AdsInteraction{})
	DB.AutoMigrate(&models.Advertiser{})
	DB.AutoMigrate(&models.Transaction{})
	DB.AutoMigrate(&models.Keyword{})
	if err != nil {
		log.Fatal(err)
	}
}
