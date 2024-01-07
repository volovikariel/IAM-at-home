package v1

import (
	"net/http"
	"strings"

	"github.com/volovikariel/IdentityManager/internal/server/gateway/users"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/users/sessions"
)

type v1Handler struct {
	userStore    users.UserStore
	sessionStore sessions.SessionStore
}

func NewHandler(userStore users.UserStore, sessionStore sessions.SessionStore) http.Handler {
	return v1Handler{
		userStore:    userStore,
		sessionStore: sessionStore,
	}
}

func (h v1Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remainingPath := strings.TrimPrefix(r.RequestURI, "/v1")
	pathParameters := strings.Split(remainingPath, "/")

	if len(pathParameters) >= 1 && pathParameters[0] == "users" {
		// /v1/users[/...]
		users.NewHandler(h.userStore, h.sessionStore).Handle(w, r)
	}

	http.NotFound(w, r)
}
