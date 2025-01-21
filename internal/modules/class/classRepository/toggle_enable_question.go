package repository

import (
	"fmt"

	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

// ToggleEnableQuestion toggles the EnableQuestion field for a class
// Returns the updated state of EnableQuestion (true = enabled, false = disabled)
func (r *classGormRepository) ToggleClassEnableQuestion(classID string) (bool, error) {
	var class models.Class

	if err := r.db.First(&class, "id = ?", classID).Error; err != nil {
		return false, fmt.Errorf("class with ID %s not found: %w", classID, err)
	}

	class.EnableQuestion = !class.EnableQuestion

	if err := r.db.Save(&class).Error; err != nil {
		return false, fmt.Errorf("failed to toggle EnableQuestion for class with ID %s: %w", classID, err)
	}

	return class.EnableQuestion, nil
}
