package internal

import (
	"net/http"

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

	var problems []schemas.Problem
	db.Model(&schemas.Problem{}).Find(&problems)

	c.JSON(http.StatusOK, gin.H{
		"result":   true,
		"problems": problems,
	})
}
