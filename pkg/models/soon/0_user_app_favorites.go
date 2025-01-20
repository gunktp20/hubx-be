package models

import (
	"time"
)

type UserAppFavorite struct {
	ID        string `gorm:"type:uuid;primaryKey;not null" json:"id"`
	UserEmail string `gorm:"type:varchar(255);index;not null;" json:"user_email"`
	AppID     string `gorm:"type:uuid;not null" json:"app_id"`
	App       App    `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"app"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
