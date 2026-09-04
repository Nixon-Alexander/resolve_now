package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Nixon-Alexander/resolve_now.git/config"
	"github.com/Nixon-Alexander/resolve_now.git/helper"
	"github.com/Nixon-Alexander/resolve_now.git/model"
	"github.com/gin-gonic/gin"
)

type DocumentRepository struct {
	db     *config.MySQLDatabase
	logger *config.Log
	qdrant *config.QDrant
}

func InitDocumentRepository(
	db *config.MySQLDatabase,
	log *config.Log,
	qdrant *config.QDrant,
) *DocumentRepository {
	return &DocumentRepository{
		db:     db,
		logger: log,
		qdrant: qdrant,
	}
}

func (dr *DocumentRepository) rollbackAndDeleteFiles(
	tx *sql.Tx,
	uploadedFiles []model.UploadedFile,
) {
	dr.db.Rollback(tx)

	for _, file := range uploadedFiles {
		err := helper.DeleteFile(file.FilePath)

		if err != nil {
			dr.logger.Error(
				"ERROR DELETE UPLOADED FILE",
				"filepath", file.FilePath,
				"error", err,
			)
		}
	}
}

// Add Files
func (dr *DocumentRepository) Add(
	c *gin.Context,
	req *model.AddDocumentsRequest,
) (*sql.Result, error) {
	tx, err := dr.db.BeginTransaction(c)

	if err != nil {
		dr.logger.Error("ERROR BEGIN TRANSACTION", "error", err)
		return nil, err
	}

	// Ambil file
	form, err := c.MultipartForm()

	if err != nil {
		dr.db.Rollback(tx)
		dr.logger.Error("ERROR GETTING MULTIPART FORM", "error", err)
		return nil, err
	}

	files := form.File["files"]

	if len(files) == 0 {
		dr.db.Rollback(tx)

		return nil, errors.New("file is required")
	}

	// Upload file
	uploadedFiles, err := helper.UploadFile(files)

	if err != nil {
		dr.db.Rollback(tx)
		dr.logger.Error("ERROR UPLOADING FILE", "error", err)
		return nil, err
	}

	// Query untuk insert document ke db
	documentQuery := `
		INSERT INTO documents (
			company_id,
			name,
			document_type,
			file_path,
			mime_type,
			created_at
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	var lastResult sql.Result

	// Insert document
	for _, file := range uploadedFiles {

		dr.logger.Info("FILE UPLOAD", "filepath", file.FilePath, "storage_path", file.FilePath, "filename", file.FileName, "mimetype", file.MimeType)

		result, err := dr.db.ExecTx(
			c,
			tx,
			documentQuery,
			req.CompanyId,
			file.FileName,
			req.DocumentType,
			file.FilePath,
			file.MimeType,
			time.Now(),
		)

		if err != nil {
			dr.rollbackAndDeleteFiles(tx, uploadedFiles)
			dr.logger.Error("ERROR INSERT DOCUMENT", "error", err)
			return nil, err
		}

		documentID, err := result.LastInsertId()

		if err != nil {
			dr.rollbackAndDeleteFiles(tx, uploadedFiles)
			dr.logger.Error("ERROR GETTING DOCUMENT ID", "error", err)
			return nil, err
		}

		// Query untuk insert document version
		versionQuery := `
			INSERT INTO document_versions (
				document_id,
				version,
				file_path,
				status,
				effective_from,
				effective_until
			)
			VALUES (?, ?, ?, ?, ?, ?)
		`

		versionResult, err := dr.db.ExecTx(
			c,
			tx,
			versionQuery,
			documentID,
			req.Version,
			file.FilePath,
			model.DOC_ACTIVE,
			req.EffectiveFrom,
			req.EffectiveUntil,
		)

		if err != nil {
			dr.rollbackAndDeleteFiles(tx, uploadedFiles)
			dr.logger.Error("ERROR INSERT DOCUMENT VERSION", "error", err)
			return nil, err
		}

		versionID, err := versionResult.LastInsertId()

		if err != nil {
			dr.rollbackAndDeleteFiles(tx, uploadedFiles)
			dr.logger.Error("ERROR GETTING DOCUMENT VERSION ID", "error", err)
			return nil, err
		}

		// Read file
		content, err := helper.ReadFile(file.FilePath)

		if err != nil {
			dr.rollbackAndDeleteFiles(tx, uploadedFiles)
			dr.logger.Error("ERROR READING DOCUMENT", "error", err)
			return nil, err
		}

		// Convert file ke text dan di parse
		text := string(content)
		sections, err := helper.ParseSections(text)
		if err != nil {
			dr.rollbackAndDeleteFiles(tx, uploadedFiles)
			dr.logger.Error("ERROR PARSING DOCUMENT", "error", err)
			return nil, err
		}

		maxChars := 1000

		// Potong document menjadi chunk
		chunks := helper.ChunkDocument(sections, maxChars)

		dr.logger.Info("DOCUMENT CHUNKED", "document_id", documentID, "version_id", versionID, "chunks", len(chunks))

		// Insert query untuk chunk
		chunkQuery := `
			INSERT INTO document_chunks (
				document_version_id,
				chunk_index,
				content,
				token_count,
				created_at
			)
			VALUES (?, ?, ?, ?, ?)
		`

		for _, chunk := range chunks {
			tokenCount := 0

			// Save chunk ke DB
			chunkResult, err := dr.db.ExecTx(
				c,
				tx,
				chunkQuery,
				versionID,
				chunk.Index,
				chunk.Content,
				tokenCount,
				time.Now(),
			)

			if err != nil {
				dr.rollbackAndDeleteFiles(tx, uploadedFiles)
				dr.logger.Error("ERROR INSERT DOCUMENT CHUNK", "error", err, "chunk_index", chunk.Index)
				return nil, err
			}

			chunkID, err := chunkResult.LastInsertId()

			if err != nil {
				dr.rollbackAndDeleteFiles(tx, uploadedFiles)
				dr.logger.Error("ERROR GETTING CHUNK ID", "error", err, "chunk_index", chunk.Index)
				return nil, err
			}

			// Embed chunk
			embedding, err := helper.EmbedText(c, chunk.Content)

			if err != nil {
				dr.rollbackAndDeleteFiles(tx, uploadedFiles)
				dr.logger.Error("ERROR EMBEDDING CHUNK", "error", err, "chunk_index", chunk.Index)
				return nil, err
			}

			// Upsert ke Qdrant
			err = dr.qdrant.UpsertChunk(
				c,
				uint64(chunkID),
				embedding,
				map[string]any{
					"company_id":          req.CompanyId,
					"document_id":         documentID,
					"document_version_id": versionID,
					"chunk_id":            chunkID,
					"chunk_index":         chunk.Index,
					"content":             chunk.Content,
				},
			)

			if err != nil {
				dr.rollbackAndDeleteFiles(tx, uploadedFiles)
				dr.logger.Error("ERROR UPSERT CHUNK TO QDRANT", "error", err, "chunk_index", chunk.Index)
				return nil, err
			}
		}

		lastResult = versionResult
	}

	if err := dr.db.Commit(tx); err != nil {
		for _, file := range uploadedFiles {
			deleteErr := helper.DeleteFile(file.FilePath)
			if deleteErr != nil {
				dr.logger.Error("ERROR DELETE FILE AFTER COMMIT ERROR", "filepath", file.FilePath, "error", deleteErr)
			}
		}

		dr.logger.Error("ERROR COMMIT TRANSACTION", "error", err)
		return nil, err
	}

	return &lastResult, nil
}
