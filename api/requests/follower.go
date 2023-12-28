package requests

type Follower struct {
	Username string `validate:"required,notblank,min=5,max=255"`
}
