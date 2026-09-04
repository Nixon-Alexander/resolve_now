package repository

import (
	"errors"

	"github.com/Nixon-Alexander/resolve_now.git/config"
	"github.com/Nixon-Alexander/resolve_now.git/helper"
	"github.com/Nixon-Alexander/resolve_now.git/model"
	"github.com/gin-gonic/gin"
)

type ChatRepository struct {
	db     *config.MySQLDatabase
	qdrant *config.QDrant
	logger *config.Log
}

func InitChatRepository(
	db *config.MySQLDatabase,
	qdrant *config.QDrant,
	log *config.Log,
) *ChatRepository {
	return &ChatRepository{
		db:     db,
		qdrant: qdrant,
		logger: log,
	}
}

// Search: embed pertanyaan user, cari chunk paling mirip di Qdrant, balikin langsung sebagai jawaban
func (cr *ChatRepository) Search(
	c *gin.Context,
	req *model.ChatRequest,
) (*model.ChatResponse, error) {

	// Embed pertanyaan
	embedding, err := helper.EmbedText(c, req.Message)

	if err != nil {
		cr.logger.Error("ERROR EMBEDDING QUERY", "error", err)
		return nil, err
	}

	limit := uint64(5)
	if req.Limit > 0 {
		limit = uint64(req.Limit)
	}

	// Search ke qdrant, difilter by company_id biar tidak bocor antar tenant
	results, err := cr.qdrant.SearchChunks(
		c,
		embedding,
		req.CompanyId,
		limit,
	)

	if err != nil {
		cr.logger.Error("ERROR SEARCHING QDRANT", "error", err)
		return nil, err
	}

	cr.logger.Info("CHUNKS FOUND", "company_id", req.CompanyId, "count", len(results))

	if len(results) == 0 {
		return nil, errors.New("tidak ditemukan informasi yang relevan")
	}

	// Susun sources
	sources := make([]model.ChatSourceChunk, 0, len(results))

	for _, r := range results {
		sources = append(sources, model.ChatSourceChunk{
			DocumentID:        r.DocumentID,
			DocumentVersionID: r.DocumentVersionID,
			ChunkIndex:        r.ChunkIndex,
			Content:           r.Content,
			Score:             r.Score,
		})
	}

	// results dari Qdrant sudah terurut dari score tertinggi -> ambil yang pertama sebagai "jawaban" utama
	response := &model.ChatResponse{
		Answer:  results[0].Content,
		Sources: sources,
	}

	return response, nil
}
