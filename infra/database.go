package infra

import (
	"blockchain-scrap/entity"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&entity.CoinDetail{})
}

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
	db       *gorm.DB
	err      error
)

func GetDBConfig() gorm.Dialector {

	dbConfig := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		host,
		port,
		user,
		password,
		dbname,
	)

	return postgres.Open(dbConfig)
}

func GetDBInstance() *gorm.DB {
	return db
}

func init() {
	db, err = gorm.Open(GetDBConfig(), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Println("Connected to DB!")
}
