package report

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/maxswjeon/contest-backend/schemas"
	"gorm.io/gorm"
)

type ReportData struct {
	Focus      *bool  `form:"focus" json:"focus" binding:"required"`
	ScreenSize string `form:"screen_size" json:"screen_size" binding:"required"`
}

func POST(c *gin.Context) {
	session := sessions.Default(c)

	if session.Get("user") == nil {
		c.JSON(http.StatusOK, gin.H{
			"result": true,
		})
		return
	}

	var data ReportData
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  err.Error(),
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	db.Model(&schemas.User{}).Where("student_number = ?", session.Get("user")).Update("last_screen_size", data.ScreenSize)

	if !(*data.Focus) {
		db.Model(&schemas.User{}).Where("student_number = ?", session.Get("user")).Update("focus_alert", true)
	}

	var user schemas.User
	db.Model(&schemas.User{}).Where("student_number = ?", session.Get("user")).First(&user)

	if user.ScreenSize != data.ScreenSize {
		db.Model(&schemas.User{}).Where("student_number = ?", session.Get("user")).Update("screen_size_alert", true)
	}

	c.JSON(200, gin.H{
		"result": true,
	})
}
