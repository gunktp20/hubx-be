package repository

import (
	classDto "github.com/gunktp20/digital-hubx-be/internal/modules/class/classDto"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
	"gorm.io/gorm"
)

func (r *classGormRepository) GetAllClasses(class_tier, keyword string, class_level *int, class_category, userEmail string, page, limit int) (*[]classDto.ClassRes, int64, error) {
	var classes []models.Class
	var total int64

	query := r.db.Model(&models.Class{})

	// Filter by is_remove (Exclude soft-deleted records)
	query = query.Where("is_remove = ?", false)

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
			// Filter sessions that are not in the past
			return db.Where("date >= CURRENT_DATE").Order("date ASC")
		}).
		Limit(limit).
		Offset(offset).
		Find(&classes)

	if result.Error != nil {
		return &[]classDto.ClassRes{}, 0, result.Error
	}

	// Map Classes to DTOs
	var classDTOs []classDto.ClassRes
	for _, class := range classes {
		isRegistered := false

		// Check registration if userEmail is provided
		if userEmail != "" {
			var count int64
			r.db.Model(&models.UserClassRegistration{}).
				Where("user_email = ? AND class_id = ?", userEmail, class.ID).
				Count(&count)

			isRegistered = count > 0
		}

		// Map class to DTO
		classDTOs = append(classDTOs, classDto.ClassRes{
			ID:             class.ID,
			Title:          class.Title,
			Description:    class.Description,
			CoverImage:     class.CoverImage,
			ClassTier:      class.ClassTier,
			ClassLevel:     class.ClassLevel,
			IsActive:       class.IsActive,
			IsRemove:       class.IsRemove,
			EnableQuestion: class.EnableQuestion,
			Order:          class.Order,
			CreatedAt:      class.CreatedAt,
			UpdatedAt:      class.UpdatedAt,
			ClassCategory:  class.ClassCategory,
			ClassSessions:  class.ClassSessions,
			IsRegistered:   isRegistered, // Include registration status
		})
	}

	return &classDTOs, total, nil
}
