package reset

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

type ResetReportData struct {
	User string `json:"user"`
}

func POST(c *gin.Context) {
	if !utils.AuthAdmin(c) {
		return
	}

	var data ResetReportData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Bad Request",
		})
		return
	}

	userID, err := uuid.Parse(data.User)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid User ID",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user schemas.User
	db.Model(&schemas.User{}).Where(userID).First(&user)
	user.ScreenSize = user.LastScreenSize
	user.ScreenSizeAlert = false
	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}
