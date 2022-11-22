package schemas

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Problem struct {
	gorm.Model
	ID                uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Title             string    `json:"title" gorm:"not null"`
	Description       string    `json:"description" gorm:"not null"`
	InternalProblemID uuid.UUID `json:"problem_id" gorm:"type:uuid;not null"`
	InternalProblem   InternalProblem
	Contests          []*Contest `json:"contests" gorm:"many2many:contest_problems;"`
	Valid             bool       `json:"valid"`
}

func (problem *Problem) BeforeCreate(tx *gorm.DB) (err error) {
	problem.ID = uuid.New()
	return nil
}

func (problem *Problem) Validate(db *gorm.DB) {
	var internalProblem InternalProblem

	db.Model(InternalProblem{}).Where("id = ?", problem.InternalProblemID).Find(&internalProblem)

	problem.Valid = internalProblem.Valid
}
