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
)

type (
	ClassUsecaseService interface {
		CreateClass(createClassReq *classDto.CreateClassReq, fileBytes []byte, fileHeader *multipart.FileHeader) (*classDto.CreateClassRes, error)
		// CreateClass(createClassReq *classDto.CreateClassReq, fileBytes []byte) (*classDto.CreateClassRes, error)
		GetAllClasses(class_tier, keyword string, class_level *int, class_category string, page int, limit int) (*[]models.Class, int64, error)
		GetClassById(classId string) (*models.Class, error)
		ToggleClassEnableQuestion(classID string) (bool, error)
	}

	classUsecase struct {
		classRepo         classRepository.ClassRepositoryService
		classCategoryRepo classCategoryRepository.ClassCategoryRepositoryService
		gcsClient         gcs.GcsClientService
	}
)

func NewClassUsecase(classRepo classRepository.ClassRepositoryService, classCategoryRepo classCategoryRepository.ClassCategoryRepositoryService, gcsClient gcs.GcsClientService) ClassUsecaseService {
	return &classUsecase{
		classRepo:         classRepo,
		classCategoryRepo: classCategoryRepo,
		gcsClient:         gcsClient,
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
