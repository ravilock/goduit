package responses

type ProfileResponse struct {
	Profile Profile `json:"profile"`
}

type Profile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio,omitempty"`
	Image     string `json:"image,omitempty"`
	Following bool   `json:"following"`
}
