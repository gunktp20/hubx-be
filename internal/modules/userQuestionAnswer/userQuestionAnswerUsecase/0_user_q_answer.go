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
	userQuestionAnswerRepo userQuestionAnswerRepository.UserQuestionAnswerRepositoryService, questionRepo questionRepository.QuestionRepositoryService,
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

// func (u *userQuestionAnswerUsecase) CreateUserQuestionAnswer(createUserQuestionAnswerReq *userQuestionAnswerDto.CreateUserQuestionAnswerReq, email string) (*userQuestionAnswerDto.CreateUserQuestionAnswerRes, error) {

// 	_, err := u.questionRepo.GetQuestionById(createUserQuestionAnswerReq.QuestionID)
// 	if err != nil {
// 		return &userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, err
// 	}
// 	// ? check class id request equal actual class id of question
// 	// if question.ClassID != createUserQuestionAnswerReq.ClassID {
// 	// 	return &userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, errors.New("actual class id of question doesn't match with class that you provided")
// 	// }

// 	// choice, err := u.choiceRepo.GetChoiceById(createUserQuestionAnswerReq.ChoiceID)
// 	if err != nil {
// 		return &userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, err
// 	}
// 	// ? check question id request equal actual question id of choice
// 	// if choice.QuestionID != createUserQuestionAnswerReq.QuestionID {
// 	// 	return &userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, errors.New("actual question id of choice doesn't match with question that you provided")
// 	// }

// 	isAnswered, err := u.userQuestionAnswerRepo.IsUserAnsweredThisQuestion(email, createUserQuestionAnswerReq.QuestionID)
// 	if err != nil {
// 		return &userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, err
// 	}

// 	if isAnswered {
// 		return &userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, errors.New("user has already answered this question")
// 	}

// 	return u.userQuestionAnswerRepo.CreateUserQuestionAnswer(createUserQuestionAnswerReq, email)
// }

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

	// ! Question from All Questions of Selected Class Display
	// if len(*questionsOfClass) > 0 {
	// 	for i, req := range *questionsOfClass { // Dereference the pointer res
	// 		data, err := json.MarshalIndent(req, "", "  ")
	// 		if err != nil {
	// 			log.Fatalf("Error marshalling JSON: %v", err)
	// 		}
	// 		var question models.Question
	// 		if err := json.Unmarshal(data, &question); err != nil {
	// 			log.Fatalf("Error unmarshalling JSON: %v", err)
	// 		}

	// 		log.Print(i+1, ".")
	// 		log.Println("   Question ID : ", question.ID)
	// 		log.Println("   Question Description : ", question.Description)
	// 		if question.QuestionType == "choice" {
	// 			log.Println("   ❌  Question Type : ", question.QuestionType)
	// 		} else {
	// 			log.Println("   ✏️  Question Type : ", question.QuestionType)
	// 		}

	// 		// ! Choices
	// 		if len(req.Choices) > 0 {
	// 			log.Print("         ")
	// 			log.Println("     Choices :")
	// 		}
	// 		if len(req.Choices) > 0 {
	// 			for _, req_choice := range req.Choices {
	// 				data, err := json.MarshalIndent(req_choice, "", "  ")
	// 				if err != nil {
	// 					log.Fatalf("Error marshalling JSON: %v", err)
	// 				}
	// 				var choice models.Choice
	// 				if err := json.Unmarshal(data, &choice); err != nil {
	// 					log.Fatalf("Error unmarshalling JSON: %v", err)
	// 				}
	// 				log.Println("       Choice :")
	// 				log.Println("      	  -    Choice ID : ", choice.ID)
	// 				log.Println("      	  -    Choice Description : ", choice.Description)
	// 				log.Println("       ")
	// 				// ! Sub Questions
	// 				if len(req_choice.SubQuestions) > 0 {
	// 					log.Println("               Sub Questions ของ Choice  :", choice.Description)
	// 				}
	// 				if len(req_choice.SubQuestions) > 0 {
	// 					for _, req_choice_sub_question := range req_choice.SubQuestions {

	// 						data, err := json.MarshalIndent(req_choice_sub_question, "", "  ")
	// 						if err != nil {
	// 							log.Fatalf("Error marshalling JSON: %v", err)
	// 						}
	// 						var choice models.Choice
	// 						if err := json.Unmarshal(data, &choice); err != nil {
	// 							log.Fatalf("Error unmarshalling JSON: %v", err)
	// 						}
	// 						log.Println("      	      -    Sub Question ID : ", req_choice_sub_question.ID)
	// 						log.Println("      	      -    Sub Question Description : ", req_choice_sub_question.Description)
	// 						if req_choice_sub_question.QuestionType == "choice" {
	// 							log.Println("      	      - ❌    Sub Question Type : ", req_choice_sub_question.QuestionType)
	// 						} else {
	// 							log.Println("      	      - ✏️    Sub Question Type : ", req_choice_sub_question.QuestionType)
	// 						}

	// 						// ! Choices SQ
	// 						if len(req_choice_sub_question.SubQuestionChoices) > 0 {
	// 							log.Println("      	 ")
	// 							log.Println("                        Choices SQ ของ SQ :", choice.Description)
	// 						}
	// 						if len(req_choice_sub_question.SubQuestionChoices) > 0 {
	// 							for _, req_choice_sub_question_choice := range req_choice_sub_question.SubQuestionChoices {

	// 								data, err := json.MarshalIndent(req_choice_sub_question_choice, "", "  ")
	// 								if err != nil {
	// 									log.Fatalf("Error marshalling JSON: %v", err)
	// 								}

	// 								var choice models.Choice
	// 								if err := json.Unmarshal(data, &choice); err != nil {
	// 									log.Fatalf("Error unmarshalling JSON: %v", err)
	// 								}
	// 								log.Println("      	               -    🔹 Choice ID SQ : ", req_choice_sub_question_choice.ID)
	// 								log.Println("      	               -    🔹 Choice Description SQ : ", req_choice_sub_question_choice.Description)
	// 								log.Println("      	 ")
	// 							}
	// 						}
	// 						log.Println("      	 ")
	// 					}
	// 				}
	// 				log.Println("      	 ")
	// 			}
	// 		}
	// 		log.Println("___________________________________________________________________________________________")
	// 	}
	// }

	// เริ่ม Transaction
	tx := u.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Rollback หากมี Panic
		}
	}()

	// TODO
	if len(*questionsOfClass) > 0 {
		for i, question := range *questionsOfClass {
			// ? Check if all provided question IDs exist for the given class
			// ตรวจสอบว่า Question ID ตรงกับที่ส่งมา
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
					ChoiceID:   createUserQuestionAnswerReqs[i].SelectedChoiceID,
					AnswerText: createUserQuestionAnswerReqs[i].AnswerText,
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
					ChoiceID:   "",
					AnswerText: createUserQuestionAnswerReqs[i].AnswerText,
				})
				if err != nil {
					tx.Rollback()
					return []userQuestionAnswerDto.CreateUserQuestionAnswerRes{}, nil
				}
			}
		}
	}

	return createUserQuestionAnswerRes, nil
}
