package responses

type User struct {
	User struct {
		Username string `json:"username,omitempty"`
		Email    string `json:"email"`
		Bio      string `json:"bio,omitempty"`
		Image    string `json:"image,omitempty"`
	} `json:"user"`
}
