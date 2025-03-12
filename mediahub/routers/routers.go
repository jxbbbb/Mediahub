package routers

import (
	"github.com/gin-gonic/gin"
	"mediahub/controller"
)

func InitRouters(api *gin.RouterGroup, c *controller.Controller) {
	v1 := api.Group("/v1")
	fileGroup := v1.Group("/file")
	fileGroup.POST("/upload", c.Upload)
}
