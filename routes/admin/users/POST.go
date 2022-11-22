package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

type UserCreationData struct {
	StudentNumber string `form:"studentNumber" json:"studentNumber" binding:"required"`
	Password      string `form:"password" json:"password" binding:"required"`
	Name          string `form:"name" json:"name" binding:"required"`
}

func POST(c *gin.Context) {
	if !utils.AuthAdmin(c) {
		return
	}

	var data UserCreationData
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  err.Error(),
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&schemas.User{
		StudentNumber:   data.StudentNumber,
		Password:        data.Password,
		Name:            data.Name,
		FocusAlert:      false,
		ScreenSizeAlert: false,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": false,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"result": true,
	})
}
