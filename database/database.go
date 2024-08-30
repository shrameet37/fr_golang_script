package database

import (
	"context"
	"face_management/logger"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	DAO_COULD_NOT_FIND_RESOURCE = 0x00000005 //05 means error getting resource from table
)

var dbPool *pgxpool.Pool

func InitializeDatabasePool() error {

	var err error
	var databaseUrl string

	databaseUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require&sslrootcert=%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_SSL_CERT_PATH"))

	logger.Log.Debug(databaseUrl)

	dbconfig, _ := pgxpool.ParseConfig(databaseUrl)
	dbPool, err = pgxpool.NewWithConfig(context.Background(), dbconfig)

	if err != nil {

		errorMsg := fmt.Sprintf("Cannot connect to database %s.Error:%s!\n", databaseUrl, err.Error())
		logger.Log.Error(errorMsg)
		return err

	}

	return nil

}

func PingDatabasePool() error {

	return dbPool.Ping(context.Background())

}

func CloseDatabasePool() {

	dbPool.Close()

}
