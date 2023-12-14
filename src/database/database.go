package database

import (
	"fmt"
	"jonathantyar/tokopedia-crawler/src/constant"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func InitDB() (*gorm.DB, error) {
	config := Config{
		Host:     viper.GetString(constant.DB_HOST),
		Port:     viper.GetInt(constant.DB_PORT),
		User:     viper.GetString(constant.DB_USER),
		Password: viper.GetString(constant.DB_PASS),
		Database: viper.GetString(constant.DB_NAME),
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
