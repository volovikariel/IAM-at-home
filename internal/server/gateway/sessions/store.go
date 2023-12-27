package sessions

type SessionStore interface {
	Add(username string, token string) error
}
