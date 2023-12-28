package requests

type Login struct {
	User struct {
		Email    *string `json:"email" validate:"required,notblank,min=5,max=255"`
		Password *string `json:"password" validate:"required,notblank,min=8,max=72"`
	} `json:"user" validate:"required"`
}
