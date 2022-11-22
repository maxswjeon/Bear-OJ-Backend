package scoreboard

import (
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/maxswjeon/contest-backend/schemas"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

const (
	NotSumbitted = iota
	Accepted
	WrongAnswer
)

// {
//   "problems": [
//	   "A",
//	   "B",
//	   "C",
//	   "D",
//   ],
//   "scoreboard": [
//	    {
//        "user": "U1",
//        "score": 3,
//        "problems": [
//	        {
//            "result": 1,
//            "count": 4
//          }
//        ]
//      }
//   ]
// }

type ScoreItem struct {
	Result int   `json:"result"`
	Count  int64 `json:"count"`
}

type TeamItem struct {
	User     string      `json:"user"`
	Score    int         `json:"score"`
	Problems []ScoreItem `json:"problems"`
}

type ScoreboardData struct {
	Problems   []string   `json:"problems"`
	Scoreboard []TeamItem `json:"scoreboard"`
}

func max(a float64, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func calculate_score(submits []schemas.Submit) float64 {
	// max(문제 배점 x 0.3, 문제 배점 - 시간 페널티 - 제출 페널티)
	// 시간 페널티 = 대회 시작 후 해당 문제에서 AC까지 걸린 시간 (분) * 문제 배점 / 대회 전체 시간(시간) / 125
	// 제출 페널티 = max(0, 결과가 AC인 제출 횟수 - 1) * 50
	sort.Slice(submits, func(i, j int) bool {
		return submits[i].CreatedAt.After(submits[j].CreatedAt)
	})

	submits = utils.Filter(submits, func(submit schemas.Submit) bool {
		return submit.Status == "AC"
	})

	// var penalty_time float64
	// var penalty_tries float64
	// if len(submits) == 0 {
	// 	penalty_time = 0
	// 	penalty_tries = 0
	// }

	return 0
}

func GET(c *gin.Context) {
	if !utils.AuthAll(c) {
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var problems []schemas.Problem
	db.Model(&schemas.Problem{}).Find(&problems)

	var data ScoreboardData

	for _, problem := range problems {
		data.Problems = append(data.Problems, problem.Title)
	}
	sort.Slice(data.Problems, func(i, j int) bool {
		return data.Problems[i] < data.Problems[j]
	})

	var users []schemas.User
	db.Model(&schemas.User{}).Find(&users)

	for _, user := range users {
		for _, problem := range problems {
			var submits []schemas.Submit
			db.Model(&schemas.Submit{}).Where("user_id = ? AND problem_id = ?", user.ID, problem.ID).Find(&submits)
		}
	}

}
