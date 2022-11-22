package users

import (
	"github.com/gin-gonic/gin"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

func GET(c *gin.Context) {
	if !utils.AuthAdmin(c) {
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var users []schemas.User
	db.Model(&schemas.User{}).Find(&users)

	c.JSON(200, gin.H{
		"result": true,
		"users":  users,
	})
}
