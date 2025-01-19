package dto

type (
	CreateUserReq struct {
		Email string `json:"email" form:"email" validate:"required,email,max=255"`
	}

	CreateUserRes struct {
		Email string `json:"email" form:"email" validate:"required,email,max=255"`
	}
)
