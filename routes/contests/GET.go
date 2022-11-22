package contests

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

func GET(c *gin.Context) {
	if !utils.AuthAll(c) {
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var contests []schemas.Contest
	db.Model(&schemas.Contest{}).Preload("Problems").Find(&contests)

	contests = utils.Filter(contests, func(contest schemas.Contest) bool {
		contest.Validate(db)
		return contest.Valid
	})

	c.JSON(http.StatusOK, gin.H{
		"result":   true,
		"contests": contests,
	})
}
