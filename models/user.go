package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          string     `gorm:"primaryKey;type:uuid;index"`
	Email       string     `gorm:"column:email;type:varchar;not null;index"`
	Username    string     `gorm:"column:username;type:varchar;not null;uniqueIndex;index"`
	Password    string     `gorm:"column:password;type:varchar;not null;"`
	IsActive    bool       `gorm:"column:is_active;default:true"`
	IsSuperuser bool       `gorm:"column:is_superuser;default:false"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp with time zone;"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;type:timestamp with time zone;default null"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;type:timestamp with time zone;default null"`
}

func (User) TableName() string {
	return "user"
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.ID = uuid.NewV4().String()
	return nil
}
