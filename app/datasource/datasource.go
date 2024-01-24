package datasource

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/okira-e/go-as-your-backend/app/tables"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect connects to the database and returns the gorm database object
func Connect() (*gorm.DB, error) {
	host, user, pass, database, port, err := getDBInfo()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, pass, database, port)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

// DisconnectOrPanic disconnects from the database or panics if an error occurs
func DisconnectOrPanic(gormDB *gorm.DB) {
	db, err := gormDB.DB()
	if err != nil {
		log.Panic(err)
	}

	err = db.Close()
	if err != nil {
		log.Panic(err)
	}
}

// Migrate runs the database migrations
func Migrate(gormDB *gorm.DB) error {
	var err error

	err = gormDB.AutoMigrate(tables.Organizations{}, tables.Products{}, tables.Projects{}, tables.Users{}, tables.Roles{})
	if err != nil {
		return err
	}

	return nil
}

// getDBInfo returns the database connection information from environment variables
func getDBInfo() (string, string, string, string, string, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return "", "", "", "", "", errors.New("environment variable DB_HOST is not set")
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		return "", "", "", "", "", errors.New("environment variable DB_USER is not set")
	}

	pass := os.Getenv("DB_PASS")
	if pass == "" {
		return "", "", "", "", "", errors.New("environment variable DB_PASS is not set")
	}

	database := os.Getenv("DB_NAME")
	if database == "" {
		return "", "", "", "", "", errors.New("environment variable DB_NAME is not set")
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		return "", "", "", "", "", errors.New("environment variable DB_PORT is not set")
	}

	return host, user, pass, database, port, nil
}
