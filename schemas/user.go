package schemas

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID              uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	StudentNumber   string    `json:"student_number" gorm:"uniqueIndex;not null"`
	Password        string    `json:"password" gorm:"not null"`
	Name            string    `json:"name" gorm:"not null"`
	ScreenSize      string    `json:"screen_size"`
	LastScreenSize  string    `json:"last_screen_size"`
	FocusAlert      bool      `json:"alert_focus" gorm:"not null"`
	ScreenSizeAlert bool      `json:"alert_screen_size" gorm:"not null"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return nil
}
