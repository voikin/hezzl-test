package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/voikin/hezzl-test/internal/controller/good"
	"github.com/voikin/hezzl-test/internal/controller/project"
	"github.com/voikin/hezzl-test/internal/service"
)

func RegisterRoutes(route *gin.Engine, service *service.Service) {
	baseRoute := route.Group("/")

	projectHandlers := project.NewProjectController(service.ProjectService)
	projectRoute := baseRoute.Group("/project")
	{
		projectRoute.POST("/create", projectHandlers.Create)
		projectRoute.PATCH("/update", projectHandlers.Update)
		projectRoute.DELETE("/remove", projectHandlers.Delete)
		projectRoute.GET("/", projectHandlers.GetProject)
	}
	baseRoute.GET("/projects/list", projectHandlers.GetProjects)

	goodHandlers := good.NewGoodController(service.GoodService)
	goodRoute := baseRoute.Group("/good")
	{
		goodRoute.POST("/create", goodHandlers.Create)
		goodRoute.PATCH("/update", goodHandlers.Update)
		goodRoute.PATCH("/reprioritize", goodHandlers.UpdateGoodPriority)
		goodRoute.DELETE("/remove", goodHandlers.Delete)
		goodRoute.GET("/", goodHandlers.GetGood)
	}
	baseRoute.GET("/goods/list", goodHandlers.GetGoods)
}
