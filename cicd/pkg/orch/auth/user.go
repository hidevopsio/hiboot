package auth

type User struct {
	Id        int32  `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	State     string `json:"state"`
	AvatarUrl string `json:"avatar_url"`
	WebUrl    string `json:"web_url"`
	Email     string `json:"email"`
	// ...
}

func Login(username, password, gitUrl string) (*User, error) {
	return nil, nil
}
