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

type SumbitData struct {
	Code      string `form:"code" json:"code" binding:"required"`
	Language  string `form:"language" json:"language" binding:"required"`
	ProblemID string `form:"problem_id" json:"problem_id" binding:"required"`
}

func POST(c *gin.Context) {
	session := sessions.Default(c)
	if !utils.AuthUser(c) {
		return
	}

	var data SumbitData
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  err.Error(),
		})
		return
	}

	if _, err := uuid.Parse(data.ProblemID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Invalid problem id",
		})
		return
	}

	student_number := session.Get("user").(string)

	db := c.MustGet("db").(*gorm.DB)

	var count int64

	// 문제 ID 확인
	db.Model(&schemas.Problem{}).Where("id = ?", data.ProblemID).Count(&count)
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": false,
			"error":  "Problem not found",
		})
		return
	}

	// 중복 제출 방지
	db.Model(&schemas.Submit{}).Where("student_number = ? AND problem_id = ? AND code = ?", student_number, data.ProblemID, data.Code).Count(&count)
	if count != 0 {
		c.JSON(http.StatusConflict, gin.H{
			"result": false,
			"error":  "Duplicate submission",
		})
		return
	}

	var user schemas.User
	db.Model(&schemas.User{}).Where("student_number = ?", student_number).First(user)

	var submit schemas.Submit
	submit.UserID = user.ID
	submit.ProblemID = uuid.MustParse(data.ProblemID)
	submit.Code = data.Code
	db.Model(&schemas.Submit{}).Create(submit)

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}
