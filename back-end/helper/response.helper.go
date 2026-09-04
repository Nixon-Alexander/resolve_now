package helper

import (
	"net/http"

	"github.com/Nixon-Alexander/resolve_now.git/model"
	"github.com/gin-gonic/gin"
)

func SendErrorJSON(c *gin.Context, code int, message string, err error) {
	response := model.ErrorResponse{
		Message: message + err.Error(),
		Status:  model.FAILED,
		Code:    code,
	}
	c.JSON(http.StatusBadRequest, response)
}

func SendSuccessJSON(c *gin.Context, code int, message string, data any) {
	response := model.SuccessResponse{
		Message: message,
		Status:  model.SUCCESS,
		Code:    code,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}
