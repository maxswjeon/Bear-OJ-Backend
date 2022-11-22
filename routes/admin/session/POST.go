package session

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/maxswjeon/contest-backend/schemas"
	"gorm.io/gorm"
)

type LoginData struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func POST(c *gin.Context) {
	session := sessions.Default(c)

	var data LoginData
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  err.Error(),
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("username = ? AND password = ?", data.Username, data.Password).First(&schemas.Admin{}).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": false,
			"error":  "Authentication failed",
		})
		return
	}

	session.Set("admin", data.Username)
	session.Save()

	c.JSON(200, gin.H{
		"result": true,
		"status": true,
	})
}
