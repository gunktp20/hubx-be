package repository

import (
	"errors"
	"time"

	classRegistrationDto "github.com/gunktp20/digital-hubx-be/internal/modules/classRegistration/classRegistrationDto"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
	"gorm.io/gorm"
)

type (
	ClassRegistrationRepositoryService interface {
		CreateClassRegistration(createClassRegistrationReq *classRegistrationDto.CreateClassRegistrationReq, email string) (*classRegistrationDto.CreateClassRegistrationRes, error)
		HasUserRegistered(email string, classID string) (bool, error)
		GetUserRegistrations(email string, page int, limit int) (*[]classRegistrationDto.GetUserRegistrationsRes, int64, error)
		CountRegistrationWithClassSessionID(classSessionID string) (int, error)
		CancelClassRegistration(email, classID string) error
		CountUserCancelledRegistrationsByEmail(userEmail, classID string) (int, error)
		ResetCancelledQuota(resetCancelledQuotaReq *classRegistrationDto.ResetCancelledQuotaReq) error
		GetUserRegistrationsByClassSessionID(classSessionID string, page int, limit int) (*[]classRegistrationDto.GetUserRegistrationsRes, int64, error)
		HasUserRegisteredByClassSessionID(email string, classSessionID string) (bool, error)
	}
)

func (r *classRegistrationGormRepository) CreateClassRegistration(createClassRegistrationReq *classRegistrationDto.CreateClassRegistrationReq, email string) (*classRegistrationDto.CreateClassRegistrationRes, error) {

	classRegistration := models.UserClassRegistration{
		UserEmail:      email,
		ClassID:        createClassRegistrationReq.ClassID,
		ClassSessionID: createClassRegistrationReq.ClassSessionID,
	}

	if err := r.db.Create(&classRegistration).Error; err != nil {
		return &classRegistrationDto.CreateClassRegistrationRes{}, err
	}

	return &classRegistrationDto.CreateClassRegistrationRes{
		ID:              classRegistration.ID,
		Email:           classRegistration.UserEmail,
		ClassID:         classRegistration.ClassID,
		ClassSessionID:  classRegistration.ClassSessionID,
		UnattendedQuota: classRegistration.UnattendedQuota,
		IsBanned:        classRegistration.IsBanned,
		RegisteredAt:    classRegistration.RegisteredAt,
		RegStatus:       classRegistration.RegStatus,
		CreatedAt:       classRegistration.CreatedAt,
		UpdatedAt:       classRegistration.UpdatedAt,
	}, nil

}

// ? HasUserRegistered checks if a user has registered for a specific class session and is not cancelled
func (r *classRegistrationGormRepository) HasUserRegistered(email string, classID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserClassRegistration{}).
		Where("user_email = ? AND class_id = ? AND reg_status != ?", email, classID, models.Cancelled).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ? HasUserRegistered checks if a user has registered for a specific class session and is not cancelled
func (r *classRegistrationGormRepository) HasUserRegisteredByClassSessionID(email string, classSessionID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserClassRegistration{}).
		Where("user_email = ? AND class_session_id = ? AND reg_status != ?", email, classSessionID, models.Cancelled).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *classRegistrationGormRepository) GetUserRegistrations(email string, page int, limit int) (*[]classRegistrationDto.GetUserRegistrationsRes, int64, error) {
	var userClassRegistrations []models.UserClassRegistration
	var total int64

	query := r.db.Model(&models.UserClassRegistration{})

	// Filter by class_level
	if email != "" {
		query = query.Where("user_email = ? AND reg_status != ?", email, models.Cancelled)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	result := query.
		// Preload("ClassSessions", func(db *gorm.DB) *gorm.DB {
		// 	return db.Select("id,class_id,max_capacity,date,start_time,end_time")
		// }).
		Limit(limit).
		Offset(offset).
		Find(&userClassRegistrations)

	var userClassRegistrationsRes []classRegistrationDto.GetUserRegistrationsRes

	for i := range userClassRegistrations {

		userClassRegistrationsRes = append(userClassRegistrationsRes, classRegistrationDto.GetUserRegistrationsRes{
			ID:              userClassRegistrations[i].ID,
			Email:           userClassRegistrations[i].UserEmail,
			ClassID:         userClassRegistrations[i].ClassID,
			ClassSessionID:  userClassRegistrations[i].ClassSessionID,
			UnattendedQuota: userClassRegistrations[i].UnattendedQuota,
			IsBanned:        userClassRegistrations[i].IsBanned,
			RegisteredAt:    userClassRegistrations[i].RegisteredAt,
			RegStatus:       userClassRegistrations[i].RegStatus,
			CreatedAt:       userClassRegistrations[i].CreatedAt,
			UpdatedAt:       userClassRegistrations[i].UpdatedAt,
		})
	}

	if result.Error != nil {
		return &[]classRegistrationDto.GetUserRegistrationsRes{}, 0, result.Error
	}

	return &userClassRegistrationsRes, total, nil
}

func (r *classRegistrationGormRepository) GetUserRegistrationsByClassSessionID(classSessionID string, page int, limit int) (*[]classRegistrationDto.GetUserRegistrationsRes, int64, error) {
	var userClassRegistrations []models.UserClassRegistration
	var total int64

	query := r.db.Model(&models.UserClassRegistration{})

	// Filter by class_level
	if classSessionID != "" {
		query = query.Where("class_session_id = ? AND reg_status != ?", classSessionID, models.Cancelled)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	result := query.
		// Preload("ClassSessions", func(db *gorm.DB) *gorm.DB {
		// 	return db.Select("id,class_id,max_capacity,date,start_time,end_time")
		// }).
		Limit(limit).
		Offset(offset).
		Find(&userClassRegistrations)

	var userClassRegistrationsRes []classRegistrationDto.GetUserRegistrationsRes

	for i := range userClassRegistrations {

		userClassRegistrationsRes = append(userClassRegistrationsRes, classRegistrationDto.GetUserRegistrationsRes{
			ID:              userClassRegistrations[i].ID,
			Email:           userClassRegistrations[i].UserEmail,
			ClassID:         userClassRegistrations[i].ClassID,
			ClassSessionID:  userClassRegistrations[i].ClassSessionID,
			UnattendedQuota: userClassRegistrations[i].UnattendedQuota,
			IsBanned:        userClassRegistrations[i].IsBanned,
			RegisteredAt:    userClassRegistrations[i].RegisteredAt,
			RegStatus:       userClassRegistrations[i].RegStatus,
			CreatedAt:       userClassRegistrations[i].CreatedAt,
			UpdatedAt:       userClassRegistrations[i].UpdatedAt,
		})
	}

	if result.Error != nil {
		return &[]classRegistrationDto.GetUserRegistrationsRes{}, 0, result.Error
	}

	return &userClassRegistrationsRes, total, nil
}

func (r *classRegistrationGormRepository) CountRegistrationWithClassSessionID(classSessionID string) (int, error) {

	var totalRegistrations int64
	result := r.db.Model(&models.UserClassRegistration{}).
		Where("class_session_id = ? AND reg_status != ?", classSessionID, models.Cancelled).
		Count(&totalRegistrations)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(totalRegistrations), nil
}

// ? CancelClassRegistration cancels a class registration by updating its status to "cancelled" if allowed
func (r *classRegistrationGormRepository) CancelClassRegistration(userEmail, classID string) error {
	// Find the UserClassRegistration record
	var registration models.UserClassRegistration
	result := r.db.Where("user_email = ? AND class_id = ? AND reg_status = ?", userEmail, classID, models.Registered).First(&registration)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("registration not found")
		}
		return result.Error
	}

	// Update the status to "cancelled" and update timestamps
	update := r.db.Model(&registration).Updates(map[string]interface{}{
		"reg_status": models.Cancelled,
		"updated_at": time.Now(),
	})

	if update.Error != nil {
		return update.Error
	}

	return nil
}

func (r *classRegistrationGormRepository) CountUserCancelledRegistrationsByEmail(userEmail, classID string) (int, error) {
	var totalCancelled int64
	result := r.db.Model(&models.UserClassRegistration{}).
		Where("user_email = ? AND class_id = ? AND reg_status = ?", userEmail, classID, models.Cancelled).
		Count(&totalCancelled)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(totalCancelled), nil
}

// ResetCancelledQuota removes all cancelled registrations for a user in a specific class
func (r *classRegistrationGormRepository) ResetCancelledQuota(resetCancelledQuotaReq *classRegistrationDto.ResetCancelledQuotaReq) error {
	// Delete all registrations with status "cancelled" for the user in the specified class
	result := r.db.Where("user_email = ? AND class_id = ? AND reg_status = ?", resetCancelledQuotaReq.UserEmail, resetCancelledQuotaReq.ClassID, models.Cancelled).
		Delete(&models.UserClassRegistration{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *classRegistrationGormRepository) CountRegistrationsByClassSessionID(classSessionID string) (int, error) {
	var count int64

	err := r.db.Model(&models.UserClassRegistration{}).
		Where("class_session_id = ?", classSessionID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
