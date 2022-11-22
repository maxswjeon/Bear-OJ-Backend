package healthcheck

import (
	"github.com/gin-gonic/gin"
)

func GET(c *gin.Context) {
	c.JSON(200, gin.H{
		"result": true,
	})
}
