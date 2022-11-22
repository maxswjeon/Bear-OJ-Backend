package schemas

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InternalProblem struct {
	gorm.Model
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Name       string    `json:"name" gorm:"not null"`
	LastUpdate float64   `json:"last_update" gorm:"not null"`
	Valid      bool      `json:"valid" gorm:"not null"`
}

func (internalProblem *InternalProblem) BeforeCreate(tx *gorm.DB) (err error) {
	internalProblem.ID = uuid.New()
	internalProblem.Valid = true
	return nil
}
