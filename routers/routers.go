package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func helloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello gin",
	})
}

func demoHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello gin",
		"code":    "200",
		"success": true,
	})
}

// SetupRouter 配置路由信息
func SetupRouter(e *gin.Engine) {
	e.GET("/hello", helloHandler)
	e.GET("/demo", demoHandler)
}
