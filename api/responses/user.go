package responses

type User struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Token    string `json:"token,omitempty"`
		Bio      string `json:"bio,omitempty"`
		Image    string `json:"image,omitempty"`
	} `json:"user"`
}
