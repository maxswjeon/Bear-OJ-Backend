package problems

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

func PUT_(c *gin.Context) {
	if !utils.AuthAdmin(c) {
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	id_raw := c.Param("id")
	id, err := uuid.Parse(id_raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid Contest ID",
		})
		return
	}

	var contest schemas.Contest
	err = db.Model(&schemas.Contest{}).Preload("Problems").Where("id = ?", id).First(&contest).Error
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

	var problem schemas.Problem
	err = db.Model(&schemas.Problem{}).Where("id = ?", pid).First(&problem).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid Problem ID",
		})
		return
	}

	db.Model(&contest).Omit("Problems.*").Association("Problems").Append(&problem)

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}
