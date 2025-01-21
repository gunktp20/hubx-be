package repository

import (
	"errors"
	"time"

	classSessionDto "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionDto"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
	"gorm.io/gorm"
)

type (
	ClassSessionRepositoryService interface {
		CreateClassSession(createClassSessionReq *classSessionDto.CreateClassSessionReq) (*classSessionDto.CreateClassSessionRes, error)
		GetAllClassSessions(class_id, class_tier string, page int, limit int) (*[]classSessionDto.ClassSessionsRes, int64, error)
		CheckSessionDateConflict(classID, classTier string, date time.Time) (bool, error)
		GetClassSessionById(classSessionID string) (*models.ClassSession, error)
		GetMaxCapacityOfClassSessionById(classSessionID string) (int, error)
		SetMaxCapacity(classSessionID string, newCapacity int) error
		CheckDateConflictForMultipleClassTiers(date time.Time, tiers []models.ClassTier) (int64, error)
		CountSessionsByDate(date time.Time) (int64, error)
		CheckClassTierDateConflict(classTier string, date time.Time) (bool, error)
		UpdateLocation(classSessionID, newLocation string) error
		DeleteClassSessionByID(classSessionID string) error
	}
)

func (r *classSessionGormRepository) CheckSessionDateConflict(classID, classTier string, date time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&models.ClassSession{}).
		Joins("JOIN classes ON classes.id = class_sessions.class_id").
		Where("(class_sessions.class_id = ? OR classes.class_tier = ?) AND DATE(class_sessions.date) = DATE(?)", classID, classTier, date).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *classSessionGormRepository) CheckClassTierDateConflict(classTier string, date time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&models.ClassSession{}).
		Joins("JOIN classes ON classes.id = class_sessions.class_id").
		Where("DATE(class_sessions.date) = DATE(?) AND classes.class_tier = ? AND classes.is_remove = false", date, classTier).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil // ? If count > 0, there are repeating days.
}

func (r *classSessionGormRepository) CreateClassSession(createClassSessionReq *classSessionDto.CreateClassSessionReq) (*classSessionDto.CreateClassSessionRes, error) {

	classSession := models.ClassSession{
		ClassID:     createClassSessionReq.ClassID,
		Date:        createClassSessionReq.Date,
		MaxCapacity: createClassSessionReq.MaxCapacity,
		StartTime:   createClassSessionReq.StartTime,
		EndTime:     createClassSessionReq.EndTime,
		Location:    createClassSessionReq.Location,
	}

	if err := r.db.Create(&classSession).Error; err != nil {
		return &classSessionDto.CreateClassSessionRes{}, err
	}

	return &classSessionDto.CreateClassSessionRes{
		ID:                 classSession.ID,
		ClassID:            classSession.ClassID,
		Date:               classSession.Date,
		MaxCapacity:        classSession.MaxCapacity,
		ClassSessionStatus: classSession.ClassSessionStatus,
		StartTime:          classSession.StartTime,
		EndTime:            classSession.EndTime,
		Location:           classSession.Location,
		CreatedAt:          classSession.CreatedAt,
		UpdatedAt:          classSession.UpdatedAt,
	}, nil
}

func (r *classSessionGormRepository) GetAllClassSessions(class_id, class_tier string, page int, limit int) (*[]classSessionDto.ClassSessionsRes, int64, error) {
	var classSessions []models.ClassSession
	var total int64

	query := r.db.Model(&models.ClassSession{}).
		Joins("JOIN classes ON classes.id = class_sessions.class_id"). // Join with Class table
		Where("classes.is_remove = false")                             // Ensure Class is not removed

	// Filter by class_id
	if class_id != "" {
		query = query.Where("class_id = ?", class_id)
	}

	// Filter by class_tier
	if class_tier != "" {
		query = query.Where("classes.class_tier = ?", class_tier)
	}

	// Filter sessions that are not in the past
	query = query.Where("date >= CURRENT_DATE")

	// Count total results
	query.Count(&total)

	// Sort by date (ascending order)
	query = query.Order("date ASC")

	offset := (page - 1) * limit
	result := query.
		Preload("Class"). // Load related Class data
		Limit(limit).
		Offset(offset).
		Find(&classSessions)

	if result.Error != nil {
		return &[]classSessionDto.ClassSessionsRes{}, 0, result.Error
	}

	// Add remaining slots field
	var classSessionsRes []classSessionDto.ClassSessionsRes

	for i := range classSessions {
		var totalRegistrations int64

		r.db.Model(&models.UserClassRegistration{}).
			Where("class_session_id = ? AND reg_status != ?", classSessions[i].ID, models.Cancelled).
			Count(&totalRegistrations)

		// Append session details
		classSessionsRes = append(classSessionsRes, classSessionDto.ClassSessionsRes{
			ID:                 classSessions[i].ID,
			ClassID:            classSessions[i].ClassID,
			Date:               classSessions[i].Date,
			MaxCapacity:        classSessions[i].MaxCapacity,
			ClassSessionStatus: classSessions[i].ClassSessionStatus,
			Class:              classSessions[i].Class,
			StartTime:          classSessions[i].StartTime,
			EndTime:            classSessions[i].EndTime,
			Location:           classSessions[i].Location,
			RemainingSeats:     int(classSessions[i].MaxCapacity) - int(totalRegistrations),
			CreatedAt:          classSessions[i].CreatedAt,
			UpdatedAt:          classSessions[i].UpdatedAt,
		})
	}

	return &classSessionsRes, total, nil
}

func (r *classSessionGormRepository) GetClassSessionById(classSessionID string) (*models.ClassSession, error) {
	var classSession = new(models.ClassSession)

	// Query with filtering to ensure the session is not past and the class is not removed
	result := r.db.
		Joins("JOIN classes ON classes.id = class_sessions.class_id").                                                        // Join with Class table
		Where("class_sessions.id = ? AND class_sessions.date >= CURRENT_DATE AND classes.is_remove = false", classSessionID). // Add is_remove check
		First(&classSession)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("class session not found, already past, or associated class has been removed")
		}
		return nil, result.Error
	}

	return classSession, nil
}

func (r *classSessionGormRepository) GetMaxCapacityOfClassSessionById(classSessionID string) (int, error) {
	var classSession = new(models.ClassSession)
	result := r.db.First(&classSession, "id = ?", classSessionID)

	if result.Error != nil {
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		return 0, errors.New("class session record not found")
	}

	return classSession.MaxCapacity, nil
}

func (r *classSessionGormRepository) SetMaxCapacity(classSessionID string, newCapacity int) error {
	var classSession = new(models.ClassSession)
	if newCapacity <= 0 {
		return errors.New("max capacity must be greater than zero")
	}
	result := r.db.First(&classSession, "id = ?", classSessionID)

	if result.Error != nil {
		return result.Error
	}

	err := r.db.Model(&models.ClassSession{}).
		Where("id = ?", classSessionID).
		Update("max_capacity", newCapacity).Error

	if err != nil {
		return err
	}

	classSession.MaxCapacity = newCapacity

	return nil
}

func (r *classSessionGormRepository) CheckDateConflictForMultipleClassTiers(date time.Time, tiers []models.ClassTier) (int64, error) {
	var count int64
	result := r.db.Raw(`
        SELECT COUNT(DISTINCT classes.class_tier) 
        FROM class_sessions 
        JOIN classes ON classes.id = class_sessions.class_id 
        WHERE class_sessions.date = ? 
        AND classes.class_tier IN ? 
        AND classes.is_remove = false
    `, date, tiers).Scan(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (r *classSessionGormRepository) CountSessionsByDate(date time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.ClassSession{}).
		Where("DATE(class_sessions.date) = DATE(?)", date).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *classSessionGormRepository) UpdateLocation(classSessionID, newLocation string) error {
	result := r.db.Model(&models.ClassSession{}).
		Where("id = ?", classSessionID).
		Update("location", newLocation)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no class session found to update location")
	}

	return nil
}

func (r *classSessionGormRepository) DeleteClassSessionByID(classSessionID string) error {
	result := r.db.Where("id = ?", classSessionID).Delete(&models.ClassSession{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("class session not found or already deleted")
	}
	return nil
}
