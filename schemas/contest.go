package schemas

import (
	"time"

	"github.com/google/uuid"
	"github.com/maxswjeon/contest-backend/utils"
	"gorm.io/gorm"
)

type Contest struct {
	gorm.Model
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Title      string     `json:"title" gorm:"not null"`
	StartTime  time.Time  `json:"time_start" gorm:"not null"`
	EndTime    time.Time  `json:"time_end" gorm:"not null"`
	FreezeTime time.Time  `json:"time_freeze" gorm:"not null"`
	Problems   []*Problem `json:"problems" gorm:"many2many:contest_problems;"`
	Valid      bool       `json:"valid"`
}

func (contest *Contest) BeforeCreate(tx *gorm.DB) (err error) {
	contest.ID = uuid.New()
	return nil
}

func (contest *Contest) Validate(db *gorm.DB) {
	contest.Valid = utils.All(utils.Map(contest.Problems, func(problem *Problem) bool {
		problem.Validate(db)
		return problem.Valid
	}))
}
