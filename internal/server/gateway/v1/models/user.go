package models

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Name string `json:"username"`
	Rel  Rel    `json:"rel"`
}

type Rel struct {
	Self string `json:"self"`
}

type UserStore interface {
	// Returns an error if user already exists
	Add(username string, password string) error
	// Returns an error if user doesn't exist
	Get(username string) (*User, error)
}
