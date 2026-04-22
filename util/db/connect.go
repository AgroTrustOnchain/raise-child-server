package db

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"

	"raise-child/constants/noti"
)

var (
	_cnn *sql.DB
)

// Database connection
func ConnectDB(logger *log.Logger, server ISQLServer) (*sql.DB, error) {
	if _cnn != nil {
		return _cnn, nil
	}

	cnn, err := sql.Open(server.GetSQLServer(), server.GetCnnStr())

	if err != nil {
		logger.Println(noti.DB_CONNECTION_ERR_MSG + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	cnn.SetMaxOpenConns(10)
	cnn.SetMaxIdleConns(5)
	cnn.SetConnMaxLifetime(2 * time.Minute)
	cnn.SetConnMaxIdleTime(30 * time.Second)

	_cnn = cnn

	return _cnn, nil
}
