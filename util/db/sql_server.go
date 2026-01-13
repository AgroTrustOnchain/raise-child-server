package db

type ISQLServer interface {
	GetSQLServer() string
	GetCnnStr() string
}
