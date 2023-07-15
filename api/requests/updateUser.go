package requests

type UpdateUser struct {
	User struct {
		Username string `json:"username" validate:"omitempty,notblank,min=5,max=255"`
		Email    string `json:"email" validate:"required,notblank,min=5,max=255,email"`
		Password string `json:"password" validate:"omitempty,notblank,min=8,max=72"`
		Bio      string `json:"bio" validate:"omitempty,notblank,min=1,max=255"`
		Image    string `json:"image" validate:"omitempty,notblank,min=1,max=65000,http_url|base64"`
	} `json:"user" validate:"required"`
}
