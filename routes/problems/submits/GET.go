package submits

import (
	"net/http"

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

	id_problem_raw := c.Param("id")
	id_problem, err := uuid.Parse(id_problem_raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid Problem ID",
		})
		return
	}

	id_user, err := uuid.Parse(session.Get("user").(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"result": false,
			"error":  "Invalid User ID",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var submit schemas.Submit
	db.Model(&schemas.Submit{}).
		Where("user_id = ?, problem_id = ?", id_user, id_problem).
		Order("created_at desc").
		First(&submit)

	c.JSON(http.StatusOK, gin.H{
		"result": true,
		"submit": submit,
	})
}
