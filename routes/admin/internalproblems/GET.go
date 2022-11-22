package internalproblems

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

	var internalproblems []schemas.InternalProblem
	db.Model(&schemas.InternalProblem{}).Find(&internalproblems)

	c.JSON(http.StatusOK, gin.H{
		"result":           true,
		"internalproblems": internalproblems,
	})
}
