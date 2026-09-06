package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLDatabase struct {
	db     *sql.DB
	logger *Log
}

func InitMySqlDatabase(log *Log) *MySQLDatabase {
	return &MySQLDatabase{
		logger: log,
	}
}

func (mysqlDB *MySQLDatabase) Connect() (*MySQLDatabase, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	connection := os.Getenv("DB_CONNECTION")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=true", user, password, connection, host, port, name)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		mysqlDB.logger.Error("Database cannot connect", "error", err)
		return nil, err
	}

	mysqlDB.db = db
	mysqlDB.logger.Info("Database connect successfully", "conn string", connString)
	return mysqlDB, nil
}

func (mysqlDB *MySQLDatabase) BeginTransaction(ctx *gin.Context) (*sql.Tx, error) {
	return mysqlDB.db.BeginTx(ctx, nil)
}

func (mysqlDB *MySQLDatabase) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (mysqlDB *MySQLDatabase) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func (mysqlDB *MySQLDatabase) Query(
	ctx context.Context,
	query string,
	args ...any,
) (*sql.Rows, error) {
	return mysqlDB.db.QueryContext(ctx, query, args...)
}

func (mysqlDB *MySQLDatabase) Exec(
	ctx context.Context,
	query string,
	args ...any,
) (sql.Result, error) {
	return mysqlDB.db.ExecContext(ctx, query, args...)
}

func (mysqlDB *MySQLDatabase) ExecTx(
	ctx context.Context,
	tx *sql.Tx,
	query string,
	args ...any,
) (sql.Result, error) {
	return tx.ExecContext(ctx, query, args...)
}

func (mysqlDB *MySQLDatabase) Close() error {
	return mysqlDB.db.Close()
}
