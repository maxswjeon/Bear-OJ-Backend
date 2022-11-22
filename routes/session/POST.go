package session

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/maxswjeon/contest-backend/schemas"
	"gorm.io/gorm"
)

type LoginData struct {
	StudentNumber string `form:"student_nubmer" json:"student_number" binding:"required"`
	Password      string `form:"password" json:"password" binding:"required"`
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
	if err := db.Where("student_number = ? AND password = ?", data.StudentNumber, data.Password).First(&schemas.User{}).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": false,
			"error":  "Authentication failed",
		})
		return
	}

	session.Set("user", data.StudentNumber)
	session.Save()

	c.JSON(200, gin.H{
		"result": true,
		"status": true,
	})
}
