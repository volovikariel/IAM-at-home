package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/volovikariel/IdentityManager/internal/server/gateway/sessions"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/users"
)

type inMemoryUserStore struct {
	users []users.User
}

func (i *inMemoryUserStore) Add(username string, password string) error {
	i.users = append(i.users, users.User{Name: username, Password: password})
	return nil
}

func (u *inMemoryUserStore) Get(username string) error {
	for _, user := range u.users {
		if user.Name == username {
			return nil
		}
	}
	return fmt.Errorf("User %q found", username)
}

type inMemorySessionStore struct {
	sessions []sessions.Session
}

func (i *inMemorySessionStore) Add(username string, token string) error {
	i.sessions = append(i.sessions, sessions.Session{Username: username, Token: token})
	return nil
}

func main() {
	startingUsers := []users.User{}
	imus := &inMemoryUserStore{users: startingUsers}
	userHandler := users.NewUserHandler(imus)

	startingSessions := []sessions.Session{}
	imss := &inMemorySessionStore{sessions: startingSessions}
	sessionHandler := sessions.NewSessionHandler(imss, imus)

	mux := http.NewServeMux()
	mux.Handle("/v1/users", userHandler)
	mux.Handle("/v1/users/sessions", sessionHandler)
	// handle all the other paths by returning a 404 Not Found error
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	const PORT int = 10000
	log.Printf("Listening on port %v\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux))
}
