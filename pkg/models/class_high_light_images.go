package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClassHighLightImage struct {
	ID       string `gorm:"type:uuid;primaryKey;not null" json:"id"`
	ImageURL string `gorm:"not null" json:"image_url"`
	ClassID  string `gorm:"type:uuid;index;not null" json:"class_id"`
	Class    Class  `gorm:"foreignKey:ClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"class"`

	Order int `gorm:"not null" json:"order"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (highLightImage *ClassHighLightImage) BeforeCreate(tx *gorm.DB) (err error) {
	if highLightImage.ID == "" {
		highLightImage.ID = uuid.New().String()
	}

	var maxOrder int
	if err := tx.Model(&ClassHighLightImage{}).
		Where("class_id = ?", highLightImage.ClassID).
		Select("COALESCE(MAX(\"order\"), 0)").
		Scan(&maxOrder).Error; err != nil {
		return err
	}

	highLightImage.Order = maxOrder + 1

	return
}
