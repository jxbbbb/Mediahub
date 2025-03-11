package routers

import (
	"github.com/gin-gonic/gin"
	"mediahub/Controller"
)

func InitRouters(api *gin.RouterGroup, c *Controller.Controller) {
	v1 := api.Group("/v1")
	fileGroup := v1.Group("/file")
	fileGroup.POST("/upload", c.Upload)
}
