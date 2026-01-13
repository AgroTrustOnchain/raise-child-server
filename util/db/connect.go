package db

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"

	"raise-child/constants/noti"
)

var _cnn *sql.DB

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

	_cnn = cnn

	return _cnn, nil
}
