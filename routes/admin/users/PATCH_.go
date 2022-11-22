package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

type UserData struct {
	StudentNumber string `form:"studentNumber" json:"studentNumber" binding:"required"`
	Password      string `form:"password" json:"password" binding:"required"`
	Name          string `form:"name" json:"name" binding:"required"`
}

func PATCH_(c *gin.Context) {
	if !utils.AuthAdmin(c) {
		return
	}

	var data UserData
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  err.Error(),
		})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid ID",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user schemas.User
	db.Model(&schemas.User{}).Where(id).First(&user)

	user.StudentNumber = data.StudentNumber
	user.Password = data.Password
	user.Name = data.Name

	db.Save(user)

	c.JSON(200, gin.H{
		"result": true,
	})
}
