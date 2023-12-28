package requests

type CreateArticle struct {
	Article struct {
		Title       string   `json:"title" validate:"required,notblank,min=5,max=255"`
		Description string   `json:"description" validate:"required,notblank,min=5,max=255"`
		Body        string   `json:"body" validate:"required,notblank"`
		TagList     []string `json:"tagList" validate:"required,unique,gt=0,dive,min=3,max=30"`
	} `json:"article" validate:"required"`
}
