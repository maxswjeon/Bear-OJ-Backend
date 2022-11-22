package contests

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

type ContestData struct {
	Title      string `json:"title" binding:"required"`
	StartTime  int64  `json:"time_start" binding:"required"`
	FreezeTime int64  `json:"time_freeze" binding:"required"`
	EndTime    int64  `json:"time_end" binding:"required"`
}

func PATCH(c *gin.Context) {
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

	var data ContestData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Bad Request",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var contest schemas.Contest
	db.Model(&schemas.Contest{}).Where(id).First(&contest)

	contest.Title = data.Title
	contest.StartTime = time.Unix(data.StartTime/1000, 0)
	contest.FreezeTime = time.Unix(data.FreezeTime/1000, 0)
	contest.EndTime = time.Unix(data.EndTime/1000, 0)

	db.Save(&contest)

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})

}
