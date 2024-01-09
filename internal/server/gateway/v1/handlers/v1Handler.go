package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/models"
)

type v1Handler struct {
	userStore    models.UserStore
	sessionStore models.SessionStore
}

func NewV1Handler(userStore models.UserStore, sessionStore models.SessionStore) http.Handler {
	return v1Handler{
		userStore:    userStore,
		sessionStore: sessionStore,
	}
}

func (h v1Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remainingPath := strings.TrimPrefix(r.RequestURI, "/v1/")
	pathParameters := strings.Split(remainingPath, "/")

	log.Printf("/v1/ hit with full path %q\n", r.RequestURI)
	if len(pathParameters) >= 1 && pathParameters[0] == "users" {
		// /v1/users[/...]
		NewUsersHandler(h.userStore, h.sessionStore).Handle(w, r)
	} else {
		http.NotFound(w, r)
	}
}
