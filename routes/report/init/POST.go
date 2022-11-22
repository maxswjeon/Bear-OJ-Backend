package report

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

type ReportData struct {
	ScreenSize string `form:"screen_size" json:"screen_size" binding:"required"`
}

func POST(c *gin.Context) {
	session := sessions.Default(c)
	if !utils.AuthUser(c) {
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
	db.Model(&schemas.User{}).Where("student_number = ?", session.Get("user")).Update("screen_size", data.ScreenSize)

	c.JSON(200, gin.H{
		"result": true,
	})
}
