package repository

import (
	"database/sql"

	"github.com/Nixon-Alexander/resolve_now.git/config"
	"github.com/Nixon-Alexander/resolve_now.git/model"
	"github.com/gin-gonic/gin"
)

type DocVersionRepository struct {
	db  *config.MySQLDatabase
	log *config.Log
}

func InitDocVersionRepository(
	db *config.MySQLDatabase,
	log *config.Log,
) *DocVersionRepository {
	return &DocVersionRepository{
		db:  db,
		log: log,
	}
}

func (dvr *DocVersionRepository) Add(c *gin.Context, req *model.AddDocumentVersionRequest) (sql.Result, error) {
	sql := "INSERT INTO document_versions (document_id, version, file_path, status, effective_from, effective_until) VALUE(?, ?, ?, ?, ?, ?)"
	result, err := dvr.db.Exec(c, sql, req.DocumentID, req.Version, req.FilePath, req.EffectiveFrom, req.EffeciveUntil)
	if err != nil {
		return nil, err
	}

	return result, nil
}
