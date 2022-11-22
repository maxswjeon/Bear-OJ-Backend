package problems

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

type ProblemData struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	ProblemID   string `json:"problem_id" binding:"required"`
}

func PATCH_(c *gin.Context) {
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

	var data ProblemData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Bad Request",
		})
		return
	}

	problemID, err := uuid.Parse(data.ProblemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid Problem ID",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var problem schemas.Problem
	db.Model(&schemas.Problem{}).Where(id).First(&problem)

	problem.Title = data.Title
	problem.Description = data.Description
	problem.InternalProblemID = problemID

	db.Save(problem)

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}
