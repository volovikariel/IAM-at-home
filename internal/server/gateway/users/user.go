package users

import "fmt"

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type UserSessionResponse struct {
	Name    string `json:"username"`
	Session string `json:"session"`
	Rel     rel    `json:"rel"`
}

type UserResponse struct {
	Name string `json:"username"`
	Rel  rel    `json:"rel"`
}

type rel struct {
	Self string `json:"self"`
}

func ValidateUsernameLength(username string, minLength int, maxLength int) error {
	if len(username) < minLength || len(username) > maxLength {
		return fmt.Errorf("Username must be between %d and %d characters long", minLength, maxLength)
	}
	return nil
}
