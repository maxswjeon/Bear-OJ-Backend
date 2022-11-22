package problems

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

func DELETE_(c *gin.Context) {
	if !utils.AuthAdmin(c) {
		return
	}

	id_raw := c.Param("id")
	id, err := uuid.Parse(id_raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid ID",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var problem schemas.Problem
	db.Model(&schemas.Problem{}).Where(id).First(&problem)
	db.Delete(&problem)

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}
