package models

type UserStore interface {
	Add(username string, password string) error
	// Returns an error if user doesn't exist
	Get(username string) (*User, error)
}
