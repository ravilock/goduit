package responses

type User struct {
	User struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
		Token    *string `json:"token"`
		Bio      *string `json:"bio"`
		Image    *string `json:"image"`
	} `json:"user"`
}
