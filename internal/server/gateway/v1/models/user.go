package models

import "fmt"

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

func ValidateStringLength(str, strName string, minLength, maxLength int) error {
	if len(str) < minLength || len(str) > maxLength {
		return fmt.Errorf("%s must be between %d and %d characters long", strName, minLength, maxLength)
	}
	return nil
}
