package models

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

type SessionStore interface {
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
