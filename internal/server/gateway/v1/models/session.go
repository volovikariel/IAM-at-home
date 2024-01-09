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

type SessionStore interface {
	// TODO: Patch to update a session, instead of overloading Add to mean refresh?
	Add(username string, token string) error
	// Returns an error if session doesn't exist
	Delete(username string, token string) error
}
