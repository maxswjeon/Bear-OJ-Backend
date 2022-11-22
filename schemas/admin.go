package schemas

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username string    `gorm:"uniqueIndex;not null"`
	Password string    `gorm:"not null"`
}

func (admin *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	admin.ID = uuid.New()
	return nil
}
