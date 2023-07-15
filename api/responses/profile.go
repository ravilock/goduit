package responses

type Profile struct {
	Profile struct {
		Username  string `json:"username"`
		Bio       string `json:"bio,omitempty"`
		Image     string `json:"image,omitempty"`
		Following bool   `json:"following"`
	} `json:"profile"`
}
