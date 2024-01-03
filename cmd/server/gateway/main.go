package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/volovikariel/IdentityManager/internal/paths"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/sessions"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/users"
)

const DEFAULT_ENV_NAME = "default.env"

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

	currFile := paths.GetCurrFile()
	currFileDir := filepath.Dir(currFile)
	defaultEnv, err := godotenv.Read(path.Join(currFileDir, DEFAULT_ENV_NAME))
	if err != nil {
		log.Fatalf("Error loading %s file\n", DEFAULT_ENV_NAME)
	}
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = defaultEnv["DEFAULT_PORT"]
		log.Printf("PORT not set, defaulting to %q\n", defaultEnv["DEFAULT_PORT"])
	}
	log.Printf("Listening on port %v\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), mux))
}
