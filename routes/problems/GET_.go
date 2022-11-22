package problems

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

func GET_(c *gin.Context) {
	if !utils.AuthUser(c) {
		return
	}

	id_raw := c.Param("id")
	id, err := uuid.Parse(id_raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid Problem ID",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var problem schemas.Problem
	db.Model(&schemas.Problem{}).Where("id = ?", id).First(&problem)

	c.JSON(http.StatusOK, gin.H{
		"result":  true,
		"problem": problem,
	})
}
