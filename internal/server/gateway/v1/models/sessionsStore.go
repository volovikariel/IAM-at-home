package models

type SessionStore interface {
	Add(username string, token string) error
	Delete(username string, token string) error
}
