package database

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

var (
	ErrDBConnectionCouldNotBeEstablished = errors.New("database connection could not be established please check your env variables")
	ErrDBCouldNotBeInitialized           = errors.New("database could not be initialized")
)

func EstablishDBConnection() (*gorm.DB, error) {
	dbUrl := os.Getenv("DB_URL")

	var dsn string
	if dbUrl != "" {
		dsn = dbUrl
	} else {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASS")
		dbName := os.Getenv("DB_NAME")
		dbSSLMode := os.Getenv("DB_SSLMODE")

		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPass, dbName, dbSSLMode,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Join(ErrDBConnectionCouldNotBeEstablished, err)
	}

	maxIdle := os.Getenv("DB_MAXIDLE")
	maxOpenConn := os.Getenv("DB_MAXOPENCONN")
	maxLifeTime := os.Getenv("DB_MAXLIFETIME")

	if maxIdle != "" || maxOpenConn != "" || maxLifeTime != "" {
		database, err := db.DB()
		if err != nil {
			return nil, err
		}

		maxIdle, _ := strconv.Atoi(maxIdle)
		maxOpenConn, _ := strconv.Atoi(maxOpenConn)
		maxLifeTime, _ := strconv.Atoi(maxLifeTime)

		database.SetMaxIdleConns(maxIdle)
		database.SetMaxOpenConns(maxOpenConn)
		database.SetConnMaxLifetime(time.Duration(maxLifeTime) * time.Hour)
	}

	return db, nil
}
