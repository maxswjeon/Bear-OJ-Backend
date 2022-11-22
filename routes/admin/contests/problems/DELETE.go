package problems

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

func DELETE(c *gin.Context) {
	if !utils.AuthAdmin(c) {
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

	pid_raw := c.Param("pid")
	pid, err := uuid.Parse(pid_raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid Problem ID",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var contest schemas.Contest
	db.Model(&schemas.Contest{}).Preload("Problems").Where("id = ?", id).First(&contest)

	problemFilter := utils.Filter(contest.Problems, func(p *schemas.Problem) bool {
		return p.ID != pid
	})

	contest.Problems = problemFilter

	db.Save(&contest)

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}
