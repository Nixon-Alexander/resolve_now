package repository

import (
	"database/sql"
	"time"

	"github.com/Nixon-Alexander/resolve_now.git/config"
	"github.com/Nixon-Alexander/resolve_now.git/model"
	"github.com/gin-gonic/gin"
)

type CompanyRepository struct {
	db     *config.MySQLDatabase
	logger *config.Log
}

func InitCompanyRepository(db *config.MySQLDatabase, logger *config.Log) *CompanyRepository {
	return &CompanyRepository{
		db:     db,
		logger: logger,
	}
}

func (cr *CompanyRepository) Add(c *gin.Context, addCompanyRequest *model.AddCompanyRequest) (sql.Result, error) {
	sql := "INSERT INTO companies (name, status, created_at, updated_at) VALUES (?, ?, ?, ?) "
	result, err := cr.db.Exec(
		c,
		sql,
		addCompanyRequest.Name,
		addCompanyRequest.Status,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		cr.logger.Error("Error when getting company", "error", err)
		return nil, err
	}

	return result, nil
}

func (cr *CompanyRepository) GetAll(c *gin.Context) (*sql.Rows, error) {
	sql := "SELECT * FROM companies"
	rows, err := cr.db.Query(c, sql)
	if err != nil {
		cr.logger.Error("Error when getting company", "error", err)
		return nil, err
	}

	return rows, nil
}

func (cr *CompanyRepository) GetById(c *gin.Context, id int) (*sql.Rows, error) {
	sql := "SELECT * FROM companies WHERE id = ?"
	rows, err := cr.db.Query(c, sql, id)
	if err != nil {
		cr.logger.Error("Error when getting company", "error", err)
		return nil, err
	}

	return rows, nil
}
