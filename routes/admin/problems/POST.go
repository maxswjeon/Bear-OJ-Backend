package problems

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

type NewProblemData struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	ProblemID   string `json:"problem_id" binding:"required"`
}

func POST(c *gin.Context) {
	if !utils.AuthAdmin(c) {
		return
	}

	var data NewProblemData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Bad Request",
		})
		return
	}

	problemID, err := uuid.Parse(data.ProblemID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Bad Request",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var count int64
	db.Find(&schemas.InternalProblem{}, problemID).Count(&count)
	if count == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "No such problem found",
		})
		return
	}

	var problem schemas.Problem
	problem.Title = data.Title
	problem.Description = data.Description
	problem.InternalProblemID = problemID

	db.Model(&schemas.Problem{}).Create(&problem)

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}
