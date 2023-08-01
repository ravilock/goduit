package requests

type GetArticle struct {
	Slug string `validate:"required,notblank,min=5,max=255"`
}
