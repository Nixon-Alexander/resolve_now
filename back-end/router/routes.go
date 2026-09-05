package router

import (
	"github.com/Nixon-Alexander/resolve_now.git/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRoutes(
	companyController *controller.CompanyController,
	documentControler *controller.DocumentController,
	chatController *controller.ChatController,
) *gin.Engine {
	var router *gin.Engine = gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false,
	}))

	{
		v1 := router.Group("/api/v1")
		// Companies
		v1.POST("/add-company", companyController.AddCompany)
		v1.GET("/get-companies", companyController.GetAllCompany)
		v1.GET("/get-company/:id", companyController.GetCompany)

		// Documents
		v1.POST("/doc/add", documentControler.AddFiles)

		// User
		v1.POST("/add-user")

		// Chat
		v1.POST("/chat/search", chatController.Search)
	}

	return router
}
