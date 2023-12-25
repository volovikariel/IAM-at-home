package users

type UserStore interface {
	Add(username string, password string) error
	Get(username string) error
}
