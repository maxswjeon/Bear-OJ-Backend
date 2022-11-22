package users

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

	db := c.MustGet("db").(*gorm.DB)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid ID",
		})
		return
	}

	db.Model(&schemas.User{}).Delete(id)

	c.JSON(200, gin.H{
		"result": true,
	})
}
