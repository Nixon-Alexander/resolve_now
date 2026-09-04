package controller

import (
	"github.com/Nixon-Alexander/resolve_now.git/config"
	"github.com/Nixon-Alexander/resolve_now.git/helper"
	"github.com/Nixon-Alexander/resolve_now.git/model"
	"github.com/Nixon-Alexander/resolve_now.git/repository"
	"github.com/gin-gonic/gin"
)

type DocumentController struct {
	dr     *repository.DocumentRepository
	logger *config.Log
}

func InitDocumentController(
	dr *repository.DocumentRepository,
	log *config.Log,
) *DocumentController {
	return &DocumentController{
		dr:     dr,
		logger: log,
	}
}

func (dc *DocumentController) AddFiles(c *gin.Context) {
	var req model.AddDocumentsRequest
	dc.logger.Info("INSERT VALUE ADD DOCUMENT", "value", req)
	if err := c.ShouldBind(&req); err != nil {
		dc.logger.Error("ERROR BINDING ADD DOC REQ", "error", err)
		helper.SendErrorJSON(c, 400, "Something wrong happened", err)
		return
	}

	result, err := dc.dr.Add(c, &req)
	dc.logger.Info("RESPONSE ADD DOC", "result", result, "error", err)
	if err != nil {
		dc.logger.Error("ERROR ADD DOC REQ", "error", err)
		helper.SendErrorJSON(c, 400, "Something wrong happend: ", err)
		return
	}

	helper.SendSuccessJSON(c, 202, "File Added", result)
}
