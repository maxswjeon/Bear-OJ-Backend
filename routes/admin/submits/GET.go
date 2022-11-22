package submits

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
	var submits []schemas.Submit
	db.Model(&schemas.Submit{}).Find(&submits)

	c.JSON(200, gin.H{
		"result":  true,
		"submits": submits,
	})
}
