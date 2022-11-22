package problems

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

func GET(c *gin.Context) {
	session := sessions.Default(c)
	if !utils.AuthUser(c) {
		return
	}

	id, err := uuid.Parse(session.Get("user").(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid ID",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var problems []schemas.Problem
	db.Model(&schemas.Problem{}).Where("start_time < ? AND end_time > ?", time.Now(), time.Now()).Select("id", "title").Find(&problems)

	var submits []schemas.Submit
	for _, problem := range problems {
		var submit schemas.Submit
		db.Model(&schemas.Submit{}).
			Where("problem_id = ? AND user_id = ?", problem.ID, id).
			Order("create_time desc").
			Limit(1).
			Select("problem_id", "status", "create_time").
			First(&submit)
		submits = append(submits, submit)
	}

	c.JSON(http.StatusOK, gin.H{
		"result":   true,
		"problems": problems,
		"submits":  submits,
	})
}
