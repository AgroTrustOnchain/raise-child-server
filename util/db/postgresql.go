package db

import (
	"os"

	"raise-child/constants/env"
)

type postgreSQL struct{}

func InitializePostgreSQL() ISQLServer {
	return &postgreSQL{}
}

// GetCnnStr implements ISQLServer.
func (p postgreSQL) GetCnnStr() string {
	return os.Getenv(env.POSTGRE_DB_CNNSTR)
}

// GetSQLServer implements ISQLServer.
func (p postgreSQL) GetSQLServer() string {
	return env.POSTGRE_SERVER
}
