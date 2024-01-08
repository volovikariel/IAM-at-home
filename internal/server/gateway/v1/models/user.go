package models

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type UserSessionResponse struct {
	Name    string `json:"username"`
	Session string `json:"session"`
	Rel     Rel    `json:"rel"`
}

type UserResponse struct {
	Name string `json:"username"`
	Rel  Rel    `json:"rel"`
}

type Rel struct {
	Self string `json:"self"`
}
