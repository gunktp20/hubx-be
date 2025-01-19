package repository

import (
	"errors"

	userDto "github.com/gunktp20/digital-hubx-be/internal/modules/user/userDto"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
	"gorm.io/gorm"
)

type userGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) UserRepositoryService {
	return &userGormRepository{db}
}

func (r *userGormRepository) IsUniqueUser(email string) bool {
	return true
}

func (r *userGormRepository) CreateOneUser(createUserReq *userDto.CreateUserReq) (*userDto.CreateUserRes, error) {

	user := &models.User{
		Email: createUserReq.Email,
	}

	if err := r.db.Create(user).Error; err != nil {
		return &userDto.CreateUserRes{}, err
	}

	return &userDto.CreateUserRes{
		Email: user.Email,
	}, nil
}

func (r *userGormRepository) GetOneUserByEmail(email string) (*models.User, error) {

	var user models.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &models.User{}, err
		}
		return &models.User{}, err
	}

	return &models.User{
		Email: user.Email,
	}, nil
}

func (r *userGormRepository) GetOrCreateUser(email string) (*models.User, error) {
	var user models.User

	// ตรวจสอบว่าผู้ใช้มีอยู่ในระบบหรือไม่
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ถ้าผู้ใช้ไม่มีในระบบ สร้างผู้ใช้ใหม่
			user = models.User{
				Email: email,
			}

			// สร้าง Record ผู้ใช้ใหม่
			if createErr := r.db.Create(&user).Error; createErr != nil {
				return nil, createErr
			}

			return &models.User{
				Email: user.Email,
			}, nil
		}
		return nil, err // กรณีอื่นๆ เช่น DB Error
	}

	// Return ผู้ใช้ที่มีอยู่
	return &models.User{
		Email: user.Email,
	}, nil
}
