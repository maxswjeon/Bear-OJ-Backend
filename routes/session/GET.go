package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GET(c *gin.Context) {
	session := sessions.Default(c)
	session.Save()

	c.JSON(200, gin.H{
		"result": true,
		"status": session.Get("user") != nil,
	});
}
