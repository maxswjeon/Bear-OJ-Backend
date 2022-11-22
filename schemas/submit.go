package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Submit struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	ProblemID uuid.UUID `gorm:"type:uuid;not null"`
	Problem   Problem
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	User      User
	Language  string    `gorm:"not null"`
	Code      string    `gorm:"not null"`
	Status    string    `gorm:"not null;default:W"`
	Time      time.Time `gorm:"not null"`
}

func (submit *Submit) BeforeCreate(tx *gorm.DB) (err error) {
	submit.ID = uuid.New()
	return nil
}
