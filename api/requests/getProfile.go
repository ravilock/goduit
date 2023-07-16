package requests

type GetProfile struct {
	Username string `validate:"required,notblank,min=5,max=255"`
}
