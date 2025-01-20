package usecase

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/gunktp20/digital-hubx-be/external/gcs"
	classDto "github.com/gunktp20/digital-hubx-be/internal/modules/class/classDto"
	classRepository "github.com/gunktp20/digital-hubx-be/internal/modules/class/classRepository"
	classCategoryRepository "github.com/gunktp20/digital-hubx-be/internal/modules/classCategory/classCategoryRepository"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
	"github.com/gunktp20/digital-hubx-be/pkg/utils"
	"gorm.io/gorm"
)

type (
	ClassUsecaseService interface {
		CreateClass(createClassReq *classDto.CreateClassReq, fileBytes []byte, fileHeader *multipart.FileHeader) (*classDto.CreateClassRes, error)
		// CreateClass(createClassReq *classDto.CreateClassReq, fileBytes []byte) (*classDto.CreateClassRes, error)
		GetAllClasses(class_tier, keyword string, class_level *int, class_category string, page int, limit int) (*[]models.Class, int64, error)
		GetClassById(classId string) (*models.Class, error)
		ToggleClassEnableQuestion(classID string) (bool, error)
		UpdateClassDetails(classID string, title *string, description *string, classCategoryName *string, classTier *string, classLevel *int) error
		UpdateClassCoverImage(classID string, fileBytes []byte, fileHeader *multipart.FileHeader) error
		SoftDeleteClass(classID string) error
	}

	classUsecase struct {
		classRepo         classRepository.ClassRepositoryService
		classCategoryRepo classCategoryRepository.ClassCategoryRepositoryService
		gcsClient         gcs.GcsClientService
		db                *gorm.DB
	}
)

func NewClassUsecase(classRepo classRepository.ClassRepositoryService, classCategoryRepo classCategoryRepository.ClassCategoryRepositoryService, gcsClient gcs.GcsClientService, db *gorm.DB) ClassUsecaseService {
	return &classUsecase{
		classRepo:         classRepo,
		classCategoryRepo: classCategoryRepo,
		gcsClient:         gcsClient,
		db:                db,
	}
}

func (u *classUsecase) CreateClass(createClassReq *classDto.CreateClassReq, fileBytes []byte, fileHeader *multipart.FileHeader) (*classDto.CreateClassRes, error) {

	// ? Check is new app group name is taken yet ?
	classTitleExists := u.classRepo.IsClassTitleExists(createClassReq.Title)
	if classTitleExists {
		return &classDto.CreateClassRes{}, errors.New("class title was taken")
	}

	// ? Is class category id that user provided is exists
	if createClassReq.ClassCategoryID != "" {
		classCategoryExists := u.classCategoryRepo.IsClassCategoryIdExists(createClassReq.ClassCategoryID)
		if !classCategoryExists {
			return &classDto.CreateClassRes{}, errors.New("class category that you provided doesn't exists")
		}
	}

	// ? Get file extension from fileBytes
	fileExtension, err := utils.GetImageFileExtension(fileBytes)
	if err != nil {
		return &classDto.CreateClassRes{}, err
	}

	// ? Generate a unique file name
	fileName := utils.GenerateFileName(16)

	// ? Upload file to GCS
	err = u.gcsClient.UploadFile(fileName, fileHeader)
	if err != nil {
		return &classDto.CreateClassRes{}, err
	}

	return u.classRepo.CreateClass(createClassReq, fileName+fileExtension)
}

func (u *classUsecase) GetAllClasses(class_tier, keyword string, class_level *int, class_category string, page int, limit int) (*[]models.Class, int64, error) {
	// เรียก Repo พร้อมส่งพารามิเตอร์ที่รองรับ
	classes, total, err := u.classRepo.GetAllClasses(class_tier, keyword, class_level, class_category, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// ? Loop through each class and update the CoverImage with a signed URL
	for i, class := range *classes {
		// ? Call GetSignedURL to retrieve a temporary signed URL
		signedUrl, err := u.gcsClient.Download(class.CoverImage)
		if err != nil {
			// ? Log the error and skip to the next iteration
			fmt.Printf("Failed to get signed URL for CoverImage: %v\n", err)
			continue
		}
		// ? Update the CoverImage of the class with the signed URL
		(*classes)[i].CoverImage = signedUrl
	}

	return classes, total, nil
}

func (u *classUsecase) GetClassById(classId string) (*models.Class, error) {
	selectedClass, err := u.classRepo.GetClassById(classId)
	if err != nil {
		return &models.Class{}, nil
	}

	signedUrl, err := u.gcsClient.Download(selectedClass.CoverImage)
	if err != nil {
		return selectedClass, nil
	}

	// ? Update the IconURL of the AppGroup with the signed URL
	selectedClass.CoverImage = signedUrl
	return selectedClass, nil
}

func (u *classUsecase) ToggleClassEnableQuestion(classID string) (bool, error) {
	// เรียกใช้ Repository เพื่อสลับค่า EnableQuestion
	newState, err := u.classRepo.ToggleClassEnableQuestion(classID)
	if err != nil {
		return false, fmt.Errorf("failed to toggle EnableQuestion for class ID %s: %w", classID, err)
	}

	// ส่งสถานะใหม่กลับไป
	return newState, nil
}

func (u *classUsecase) UpdateClassDetails(classID string, title *string, description *string, classCategoryName *string, classTier *string, classLevel *int) error {
	tx := u.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updates := make(map[string]interface{})

	if title != nil {
		updates["title"] = *title
	}

	if description != nil {
		updates["description"] = *description
	}

	if classCategoryName != nil {
		var classCategory models.ClassCategory
		err := tx.Where("category_name = ?", *classCategoryName).First(&classCategory).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("class category name does not exist")
		}
		updates["class_category_id"] = classCategory.ID
	}

	if classTier != nil {
		updates["class_tier"] = *classTier
	}

	if classLevel != nil {
		if *classLevel <= 0 {
			tx.Rollback()
			return fmt.Errorf("class level must be a positive number")
		}
		updates["class_level"] = *classLevel
	}

	// อัปเดตข้อมูลทั้งหมด
	err := u.classRepo.UpdateClassDetailsWithTransaction(tx, classID, updates)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update class details: %w", err)
	}

	// Commit Transaction
	tx.Commit()

	return nil
}

func (u *classUsecase) UpdateClassCoverImage(classID string, fileBytes []byte, fileHeader *multipart.FileHeader) error {

	_, err := u.classRepo.GetClassById(classID)
	if err != nil {
		return errors.New("class not found")
	}

	// ? Get file extension from fileBytes
	fileExtension, err := utils.GetImageFileExtension(fileBytes)
	if err != nil {
		return err
	}

	// ? Generate a unique file name
	fileName := utils.GenerateFileName(16)

	// ? Upload file to GCS
	err = u.gcsClient.UploadFile(fileName, fileHeader)
	if err != nil {
		return err
	}

	err = u.classRepo.UpdateClassCoverImage(classID, fileName+fileExtension)
	if err != nil {
		return err
	}

	return nil
}

func (u *classUsecase) SoftDeleteClass(classID string) error {

	_, err := u.classRepo.GetClassById(classID)
	if err != nil {
		return errors.New("class not found")
	}

	err = u.classRepo.SoftDeleteClass(classID)
	if err != nil {
		return fmt.Errorf("failed to soft delete class: %w", err)
	}

	return nil
}
