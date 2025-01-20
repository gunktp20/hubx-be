package repository

import (
	"github.com/gunktp20/digital-hubx-be/pkg/models"
	"gorm.io/gorm"
)

func (r *classGormRepository) GetAllClasses(class_tier, keyword string, class_level *int, class_category string, page int, limit int) (*[]models.Class, int64, error) {
	var classes []models.Class
	var total int64

	query := r.db.Model(&models.Class{})

	// Filter by keyword
	if keyword != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// Filter by class_tier
	if class_tier != "" {
		query = query.Where("class_tier = ?", class_tier)
	}

	// Filter by class_level
	if class_level != nil {
		query = query.Where("class_level = ?", *class_level)
	}

	// Filter by class_category
	if class_category != "" {
		query = query.Joins("JOIN class_categories ON class_categories.id = classes.class_category_id").
			Where("class_categories.category_name = ?", class_category)
	}

	// Count total results
	query.Count(&total)

	// Calculate OFFSET and apply Limit and Offset
	offset := (page - 1) * limit
	result := query.
		Preload("ClassCategory").
		Preload("ClassSessions", func(db *gorm.DB) *gorm.DB {
			// Sort ClassSessions by date (ascending)
			return db.Order("date ASC")
		}).
		Limit(limit).   // Apply Limit
		Offset(offset). // Apply Offset
		Find(&classes)

	if result.Error != nil {
		return &[]models.Class{}, 0, result.Error
	}

	return &classes, total, nil
}
