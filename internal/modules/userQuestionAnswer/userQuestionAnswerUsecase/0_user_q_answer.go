package usecase

import (
	"errors"

	choiceRepository "github.com/gunktp20/digital-hubx-be/internal/modules/choice/choiceRepository"
	questionRepository "github.com/gunktp20/digital-hubx-be/internal/modules/question/questionRepository"
	userQuestionAnswerDto "github.com/gunktp20/digital-hubx-be/internal/modules/userQuestionAnswer/userQuestionAnswerDto"
	userQuestionAnswerRepository "github.com/gunktp20/digital-hubx-be/internal/modules/userQuestionAnswer/userQuestionAnswerRepository"
	userSubQuestionAnswerDto "github.com/gunktp20/digital-hubx-be/internal/modules/userSubQuestionAnswer/userSubQuestionAnswerDto"
	userSubQuestionAnswerRepository "github.com/gunktp20/digital-hubx-be/internal/modules/userSubQuestionAnswer/userSubQuestionAnswerRepository"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
	"gorm.io/gorm"
)

type (
	UserQuestionAnswerUsecaseService interface {
		// CreateUserQuestionAnswer(createUserQuestionAnswerReq *userQuestionAnswerDto.CreateUserQuestionAnswerReq, emailStr string) (*userQuestionAnswerDto.CreateUserQuestionAnswerRes, error)
		GetUserQuestionAnswersWithClassId(email, classID string, page int, limit int) (*[]userQuestionAnswerDto.GetUserQuestionAnswerRes, int64, error)
		CreateMultipleUserQuestionAnswers(createUserQuestionAnswerReqs []userQuestionAnswerDto.CreateUserQuestionAnswerReq, classID, email string) ([]userQuestionAnswerDto.CreateUserQuestionAnswerRes, error)
	}

	userQuestionAnswerUsecase struct {
		userQuestionAnswerRepo    userQuestionAnswerRepository.UserQuestionAnswerRepositoryService
		userSubQuestionAnswerRepo userSubQuestionAnswerRepository.UserSubQuestionAnswerRepositoryService
		questionRepo              questionRepository.QuestionRepositoryService
		choiceRepo                choiceRepository.ChoiceRepositoryService
		db                        *gorm.DB
		// bucketClient gcs.GcsClientService
	}
)

func NewUserQuestionAnswerUsecase(
	userQuestionAnswerRepo userQuestionAnswerRepository.UserQuestionAnswerRepositoryService,
	questionRepo questionRepository.QuestionRepositoryService,
	choiceRepo choiceRepository.ChoiceRepositoryService,
	userSubQuestionAnswerRepo userSubQuestionAnswerRepository.UserSubQuestionAnswerRepositoryService,
	db *gorm.DB,
) UserQuestionAnswerUsecaseService {
	return &userQuestionAnswerUsecase{
		userQuestionAnswerRepo:    userQuestionAnswerRepo,
		questionRepo:              questionRepo,
		choiceRepo:                choiceRepo,
		userSubQuestionAnswerRepo: userSubQuestionAnswerRepo,
		db:                        db,
	}
}

func (u *userQuestionAnswerUsecase) GetUserQuestionAnswersWithClassId(email, classID string, page int, limit int) (*[]userQuestionAnswerDto.GetUserQuestionAnswerRes, int64, error) {

	userQuestionAnswers, total, err := u.userQuestionAnswerRepo.GetUserQuestionAnswersWithClassId(email, classID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return userQuestionAnswers, total, nil
}

func (u *userQuestionAnswerUsecase) CreateMultipleUserQuestionAnswers(createUserQuestionAnswerReqs []userQuestionAnswerDto.CreateUserQuestionAnswerReq, classID, email string) ([]userQuestionAnswerDto.CreateUserQuestionAnswerRes, error) {

	var createUserQuestionAnswerRes []userQuestionAnswerDto.CreateUserQuestionAnswerRes

	questionsOfClass, _, err := u.questionRepo.GetQuestionsByClassID(classID, 1, 100)
	if err != nil {
		return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, err
	}

	// ? init transaction
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Rollback if there is panic
		}
	}()

	// TODO
	if len(*questionsOfClass) > 0 {
		for i, question := range *questionsOfClass {

			// ? Check if all provided question IDs exist for the given class
			if question.ID != createUserQuestionAnswerReqs[i].QuestionID {
				tx.Rollback()
				return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{},
					errors.New("invalid question ID provided for the given class")
			}

			// ? Separate 2 operations for choice and text question types
			if question.QuestionType == "choice" {

				// ? Reject the request if no selected choice is provided
				if createUserQuestionAnswerReqs[i].SelectedChoiceID == "" {
					tx.Rollback()
					return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{},
						errors.New("no selected choice provided for the choice of question")
				}

				// ? Check if the selected choice exists in the question's choices
				choiceExists := false
				var selectedChoice models.Choice
				for _, choice := range question.Choices {
					if choice.ID == createUserQuestionAnswerReqs[i].SelectedChoiceID {
						choiceExists = true
						selectedChoice = choice
						break
					}
				}

				if !choiceExists {
					tx.Rollback()
					return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, errors.New("the selected choice ID does not exist for the given question")
				}

				// ? If has no any errors insert into database
				_, err := u.userQuestionAnswerRepo.CreateUserQuestionAnswer(tx, &models.UserQuestionAnswer{
					UserEmail:  email,
					QuestionID: question.ID,
					ClassID:    question.ClassID,
					ChoiceID: func() *string {
						if createUserQuestionAnswerReqs[i].SelectedChoiceID == "" {
							return nil
						}
						return &createUserQuestionAnswerReqs[i].SelectedChoiceID // แปลงเป็น pointer
					}(),
					AnswerText: func() *string {
						if createUserQuestionAnswerReqs[i].AnswerText == "" {
							return nil
						}
						return &createUserQuestionAnswerReqs[i].AnswerText // แปลงเป็น pointer
					}(),
				})

				if err != nil {
					tx.Rollback()
					return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, nil
				}

				// ? Check if the selected choice has any sub-questions ?
				if len(selectedChoice.SubQuestions) > 0 || selectedChoice.SubQuestions != nil {

					for ii, subQuestion := range selectedChoice.SubQuestions {

						if subQuestion.ID != createUserQuestionAnswerReqs[i].SubQuestionsAnswers[ii].SubQuestionID {
							tx.Rollback()
							return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{},
								errors.New("invalid sub question ID provided for the given sub question")
						}

						if subQuestion.QuestionType == "choice" {
							// ? Reject the request if no selected choice is provided
							if createUserQuestionAnswerReqs[i].SubQuestionsAnswers[ii].SelectedSubQuestionChoiceID == "" {
								tx.Rollback()
								return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{},
									errors.New("no selected sub question choice provided for the choice of sub question")
							}

							// ? Check if the selected sub choice exists in the sub question's choices
							subQChoiceExists := false
							var selectedSubQChoice models.SubQuestionChoice
							for _, subChoice := range subQuestion.SubQuestionChoices {
								if subChoice.ID == createUserQuestionAnswerReqs[i].SubQuestionsAnswers[ii].SelectedSubQuestionChoiceID {
									subQChoiceExists = true
									selectedSubQChoice = subChoice
									break
								}
							}

							if !subQChoiceExists {
								tx.Rollback()
								return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, errors.New("the selected sub question choice ID does not exist for the given sub question")
							}

							// ? If has no any errors insert into database
							_, err := u.userSubQuestionAnswerRepo.CreateUserSubQuestionAnswer(tx, &userSubQuestionAnswerDto.CreateUserSubQuestionAnswerReq{
								SubQuestionChoiceID: createUserQuestionAnswerReqs[i].SubQuestionsAnswers[ii].SelectedSubQuestionChoiceID,
								SubQuestionID:       selectedSubQChoice.SubQuestionID,
								ClassID:             question.ClassID,
								AnswerText:          createUserQuestionAnswerReqs[i].SubQuestionsAnswers[ii].AnswerText,
							}, email)
							if err != nil {
								tx.Rollback()
								return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, nil
							}

						} else if subQuestion.QuestionType == "text" {

							if createUserQuestionAnswerReqs[i].SubQuestionsAnswers[ii].AnswerText == "" {
								tx.Rollback()
								return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{},
									errors.New("no answer text provided for the sub question")
							}

							// ? If has no any errors insert into database
							_, err := u.userSubQuestionAnswerRepo.CreateUserSubQuestionAnswer(tx, &userSubQuestionAnswerDto.CreateUserSubQuestionAnswerReq{
								SubQuestionChoiceID: "",
								SubQuestionID:       createUserQuestionAnswerReqs[i].SubQuestionsAnswers[ii].SubQuestionID,
								ClassID:             question.ClassID,
								AnswerText:          createUserQuestionAnswerReqs[i].SubQuestionsAnswers[ii].AnswerText,
							}, email)
							if err != nil {
								tx.Rollback()
								return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, nil
							}

						}
					}
				}
			} else if question.QuestionType == "text" {
				if createUserQuestionAnswerReqs[i].AnswerText == "" {
					tx.Rollback()
					return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{},
						errors.New("no answer text provided for the question")
				}

				// ? If has no any errors insert into database
				_, err := u.userQuestionAnswerRepo.CreateUserQuestionAnswer(tx, &models.UserQuestionAnswer{
					UserEmail:  email,
					QuestionID: question.ID,
					ClassID:    question.ClassID,
					ChoiceID:   nil,
					AnswerText: func() *string {
						if createUserQuestionAnswerReqs[i].AnswerText == "" {
							return nil
						}
						return &createUserQuestionAnswerReqs[i].AnswerText
					}(),
				})

				if err != nil {
					tx.Rollback()
					return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, nil
				}
			}
		}
	}

	// ? Commit transaction
	tx.Commit()

	return createUserQuestionAnswerRes, nil
}
