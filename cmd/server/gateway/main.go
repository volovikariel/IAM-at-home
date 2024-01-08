package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/volovikariel/IdentityManager/internal/server/gateway"
	v1 "github.com/volovikariel/IdentityManager/internal/server/gateway/v1"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/models"
)

type inMemoryUserStore struct {
	models.UserStore

	users []models.User
}

func (i *inMemoryUserStore) Add(username string, password string) error {
	i.users = append(i.users, models.User{Name: username, Password: password})
	return nil
}

func (u *inMemoryUserStore) Get(username string) (*models.User, error) {
	for _, user := range u.users {
		if user.Name == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("User %q not found", username)
}

type inMemorySessionStore struct {
	models.SessionStore

	sessions []models.Session
}

func (i *inMemorySessionStore) Add(username string, token string) error {
	i.sessions = append(i.sessions, models.Session{Username: username, Token: token})
	return nil
}

func main() {
	serverConfig := gateway.ServerConfig
	port := flag.String("p", "", "Port to listen on")
	host := flag.String("h", "", "Host to listen on")
	flag.Parse()
	if port == nil || *port == "" {
		log.Printf("Port not set, defaulting to %q\n", serverConfig.Port)
	} else {
		serverConfig.Port = *port
	}
	if host == nil || *host == "" {
		log.Printf("Host not set, defaulting to %q\n", serverConfig.Host)
	} else {
		serverConfig.Host = *host
	}

	v1Handler := v1.NewHandler(&inMemoryUserStore{}, &inMemorySessionStore{})
	http.Handle("/v1/", v1Handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	serverUrl := fmt.Sprintf("%s:%s", serverConfig.Host, serverConfig.Port)
	log.Printf("Listening on %s\n", serverUrl)
	log.Fatal(http.ListenAndServe(serverUrl, nil))
}
