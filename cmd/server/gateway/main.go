package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/volovikariel/IdentityManager/cmd/server/gateway/internal/sessions"
	"github.com/volovikariel/IdentityManager/cmd/server/gateway/internal/users"
)

// Specify output file for the responses to the requests.
const PORT int = 10000

var startingUsers = []users.User{}
var imus = &inMemoryUserStore{users: startingUsers}

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

func handleRequests(out *os.File) {
	userHandler := users.NewUserHandler(imus)

	startingSessions := []sessions.Session{}
	imss := &inMemorySessionStore{sessions: startingSessions}
	sessionHandler := sessions.NewSessionHandler(imss, imus)

	mux := http.NewServeMux()
	mux.Handle("/users", userHandler)
	mux.Handle("/sessions", sessionHandler)
	// handle all the other paths by returning a not found error
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	fmt.Fprintf(out, "Listening on port %v\n", PORT)
	log.Fatal(out, http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux))
}

func main() {
	handleRequests(os.Stdout)
}
