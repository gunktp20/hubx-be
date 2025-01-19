package di

import (
	choiceRepo "github.com/gunktp20/digital-hubx-be/internal/modules/choice/choiceRepository"
	classRepo "github.com/gunktp20/digital-hubx-be/internal/modules/class/classRepository"
	classCategoryRepo "github.com/gunktp20/digital-hubx-be/internal/modules/classCategory/classCategoryRepository"
	classRegistrationRepo "github.com/gunktp20/digital-hubx-be/internal/modules/classRegistration/classRegistrationRepository"
	classSessionRepo "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionRepository"
	questionRepo "github.com/gunktp20/digital-hubx-be/internal/modules/question/questionRepository"
	subQuestionRepo "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestion/subQuestionRepository"
	subQuestionChoiceRepo "github.com/gunktp20/digital-hubx-be/internal/modules/subQuestionChoice/subQuestionChoiceRepository"
	userRepo "github.com/gunktp20/digital-hubx-be/internal/modules/user/userRepository"
	userQuestionAnswerRepo "github.com/gunktp20/digital-hubx-be/internal/modules/userQuestionAnswer/userQuestionAnswerRepository"
	userSubQuestionAnswerRepo "github.com/gunktp20/digital-hubx-be/internal/modules/userSubQuestionAnswer/userSubQuestionAnswerRepository"
	"github.com/gunktp20/digital-hubx-be/pkg/config"
	"github.com/gunktp20/digital-hubx-be/pkg/database"
)

type Repositories struct {
	ClassRepo                 classRepo.ClassRepositoryService
	ChoiceRepo                choiceRepo.ChoiceRepositoryService
	ClassCategory             classCategoryRepo.ClassCategoryRepositoryService
	ClassRegistrationRepo     classRegistrationRepo.ClassRegistrationRepositoryService
	ClassSessionRepo          classSessionRepo.ClassSessionRepositoryService
	QuestionRepo              questionRepo.QuestionRepositoryService
	SubQuestionRepo           subQuestionRepo.SubQuestionRepositoryService
	SubQuestionChoiceRepo     subQuestionChoiceRepo.SubQuestionChoiceRepositoryService
	UserRepo                  userRepo.UserRepositoryService
	UserQuestionAnswerRepo    userQuestionAnswerRepo.UserQuestionAnswerRepositoryService
	UserSubQuestionAnswerRepo userSubQuestionAnswerRepo.UserSubQuestionAnswerRepositoryService
}

type Container struct {
	Repositories Repositories
}

func NewContainer(conf *config.Config, db database.Database) *Container {
	return &Container{
		Repositories: Repositories{
			ClassRepo:                 classRepo.NewClassGormRepository(db.GetDb()),
			ChoiceRepo:                choiceRepo.NewChoiceGormRepository(db.GetDb()),
			ClassCategory:             classCategoryRepo.NewClassCategoryGormRepository(db.GetDb()),
			ClassRegistrationRepo:     classRegistrationRepo.NewClassRegistrationGormRepository(db.GetDb()),
			ClassSessionRepo:          classSessionRepo.NewClassSessionGormRepository(db.GetDb()),
			QuestionRepo:              questionRepo.NewQuestionGormRepository(db.GetDb()),
			SubQuestionRepo:           subQuestionRepo.NewSubQuestionGormRepository(db.GetDb()),
			SubQuestionChoiceRepo:     subQuestionChoiceRepo.NewSubQuestionChoiceGormRepository(db.GetDb()),
			UserRepo:                  userRepo.NewUserGormRepository(db.GetDb()),
			UserQuestionAnswerRepo:    userQuestionAnswerRepo.NewUserQuestionAnswerGormRepository(db.GetDb()),
			UserSubQuestionAnswerRepo: userSubQuestionAnswerRepo.NewUserSubQuestionAnswerGormRepository(db.GetDb()),
			// Initialize other repositories here
		},
	}
}
