package models

import (
	"time"

	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

type UserAppFavorite struct {
	ID        string      `gorm:"type:uuid;primaryKey;not null" json:"id"`
	UserID    string      `gorm:"type:uuid;not null" json:"user_id"`                                           // เพิ่ม
	AppID     string      `gorm:"type:uuid;not null" json:"app_id"`                                            // เพิ่ม
	App       App         `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"app"`   // เพิ่ม
	User      models.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"` // เพิ่ม
	CreatedAt time.Time
	UpdatedAt time.Time
}
