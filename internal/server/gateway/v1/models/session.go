package models

import "fmt"

type Session struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type UserSessionResponse struct {
	Name    string `json:"username"`
	Session string `json:"session"`
	Rel     Rel    `json:"rel"`
}

type UpdateUserRequest struct {
	Password string `json:"password"`
	Session  string `json:"session"`
}

func (u *UpdateUserRequest) Validate() error {
	if u.Password == "" {
		return fmt.Errorf("Password is required")
	}
	if u.Session == "" {
		return fmt.Errorf("Session is required")
	}
	return nil
}

type SessionStore interface {
	Get(username string) (*Session, error)
	// TODO: Patch to update a session, instead of overloading Add to mean refresh?
	Add(username string, token string) error
	// Returns an error if session doesn't exist
	Delete(username string, token string) error
}

type InMemorySessionStore struct {
	SessionStore

	sessions []Session
}

func (i *InMemorySessionStore) Add(username string, token string) error {
	i.sessions = append(i.sessions, Session{Username: username, Token: token})
	return nil
}

func (i *InMemorySessionStore) Get(username string) (*Session, error) {
	for _, session := range i.sessions {
		if session.Username == username {
			return &session, nil
		}
	}
	return nil, fmt.Errorf("Session %q not found", username)
}

func (i *InMemorySessionStore) Delete(username string, token string) error {
	for sessionIdx, session := range i.sessions {
		if session.Username == username && session.Token == token {
			i.sessions = append(i.sessions[:sessionIdx], i.sessions[sessionIdx+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Session %q not found", username)
}
