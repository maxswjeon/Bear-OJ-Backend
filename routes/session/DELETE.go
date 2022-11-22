package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/maxswjeon/contest-backend/schemas"
	"gorm.io/gorm"
)

func DELETE(c *gin.Context) {
	session := sessions.Default(c)

	db := c.MustGet("db").(*gorm.DB)
	db.Model(&schemas.User{}).
	   Where("student_number = ?", session.Get("user")).
		 Update("screen_size", nil).
		 Update("focus_alert", false).
		 Update("screen_size_alert", false)

	session.Clear()
	session.Save()

	c.JSON(200, gin.H{
		"result": true,
	});
}
