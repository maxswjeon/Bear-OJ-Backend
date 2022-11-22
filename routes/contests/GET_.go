package contests

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

func GET_(c *gin.Context) {
	if !utils.AuthAll(c) {
		return
	}

	id_raw := c.Param("id")
	id, err := uuid.Parse(id_raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid Contest ID",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var contest schemas.Contest
	db.Model(&schemas.Contest{}).Preload("Problems").Where("id = ?", id).First(&contest)

	c.JSON(http.StatusOK, gin.H{
		"result":  true,
		"contest": contest,
	})
}
