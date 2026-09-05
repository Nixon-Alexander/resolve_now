package repository

import (
	"errors"
	"strings"

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

	embedding, err := helper.EmbedText(c, req.Message)

	if err != nil {
		cr.logger.Error("ERROR EMBEDDING QUERY", "error", err)
		return nil, err
	}

	limit := uint64(5)
	if req.Limit > 0 {
		limit = uint64(req.Limit)
	}

	results, err := cr.qdrant.SearchChunks(c, embedding, req.CompanyId, limit)

	if err != nil {
		cr.logger.Error("ERROR SEARCHING QDRANT", "error", err)
		return nil, err
	}

	cr.logger.Info("CHUNKS FOUND", "company_id", req.CompanyId, "count", len(results))

	if len(results) == 0 {
		return nil, errors.New("tidak ditemukan informasi yang relevan")
	}

	sources := make([]model.ChatSourceChunk, 0, len(results))
	var contextBuilder strings.Builder

	for _, r := range results {
		sources = append(sources, model.ChatSourceChunk{
			DocumentID:        r.DocumentID,
			DocumentVersionID: r.DocumentVersionID,
			ChunkIndex:        r.ChunkIndex,
			Content:           r.Content,
			Score:             r.Score,
		})

		contextBuilder.WriteString(r.Content)
		contextBuilder.WriteString("\n\n---\n\n")
	}

	// Generate jawaban natural pakai LLM, berdasarkan chunk yang ditemukan
	answer, err := helper.AskLLM(c, req.Message, contextBuilder.String())

	if err != nil {
		cr.logger.Error("ERROR ASKING LLM", "error", err)
		answer = results[0].Content
	}

	return &model.ChatResponse{
		Answer:  answer,
		Sources: sources,
	}, nil
}
