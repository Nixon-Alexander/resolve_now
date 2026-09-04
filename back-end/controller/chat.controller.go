package controller

import (
	"net/http"

	"github.com/Nixon-Alexander/resolve_now.git/config"
	"github.com/Nixon-Alexander/resolve_now.git/model"
	"github.com/Nixon-Alexander/resolve_now.git/repository"
	"github.com/gin-gonic/gin"
)

type ChatController struct {
	chatRepo *repository.ChatRepository
	logger   *config.Log
}

func InitChatController(
	chatRepo *repository.ChatRepository,
	logger *config.Log,
) *ChatController {
	return &ChatController{
		chatRepo: chatRepo,
		logger:   logger,
	}
}

func (cc *ChatController) Search(c *gin.Context) {
	var req model.ChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := cc.chatRepo.Search(c, &req)

	if err != nil {
		cc.logger.Error("ERROR CHAT SEARCH", "error", err)
		c.JSON(http.StatusOK, gin.H{
			"answer":  "Maaf, saya tidak menemukan informasi yang relevan di dokumen.",
			"sources": []model.ChatSourceChunk{},
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
