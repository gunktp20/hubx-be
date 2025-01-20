package repository

import (
	"errors"

	classDto "github.com/gunktp20/digital-hubx-be/internal/modules/class/classDto"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
	"gorm.io/gorm"
)

type (
	ClassRepositoryService interface {
		CreateClass(createClassReq *classDto.CreateClassReq, ThumbnailUrl string) (*classDto.CreateClassRes, error)
		IsClassTitleExists(classTitle string) bool
		GetAllClasses(class_tier, keyword string, class_level *int, class_category string, page int, limit int) (*[]models.Class, int64, error)
		GetClassById(classId string) (*models.Class, error)
		DeleteClassById(classId string) (*models.Class, error)
		ToggleClassEnableQuestion(classID string) (bool, error)
		UpdateClassDetailsWithTransaction(tx *gorm.DB, classID string, updates map[string]interface{}) error
		UpdateClassCoverImage(classID string, coverImagePath string) error
		SoftDeleteClass(classID string) error
	}
)

func (r *classGormRepository) UpdateClassDetailsWithTransaction(tx *gorm.DB, classID string, updates map[string]interface{}) error {
	result := tx.Model(&models.Class{}).Where("id = ?", classID).Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no class found with the given ID")
	}

	return nil
}

func (r *classGormRepository) UpdateClassCoverImage(classID string, coverImagePath string) error {
	// อัปเดต cover_image ของคลาสที่ตรงกับ classID
	result := r.db.Model(&models.Class{}).
		Where("id = ?", classID).
		Update("cover_image", coverImagePath)

	if result.Error != nil {
		return result.Error
	}

	// ตรวจสอบว่ามีการอัปเดตหรือไม่
	if result.RowsAffected == 0 {
		return errors.New("no class found with the given ID")
	}

	return nil
}

func (r *classGormRepository) SoftDeleteClass(classID string) error {
	// อัปเดตค่า is_remove เป็น true
	result := r.db.Model(&models.Class{}).
		Where("id = ?", classID).
		Update("is_remove", true)

	if result.Error != nil {
		return result.Error
	}

	// ตรวจสอบว่ามีการอัปเดตสำเร็จหรือไม่
	if result.RowsAffected == 0 {
		return errors.New("no class found with the given ID")
	}

	return nil
}
