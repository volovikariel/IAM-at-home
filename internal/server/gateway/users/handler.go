package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/volovikariel/IdentityManager/internal/server/gateway"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/users/sessions"
)

// TODO: Maybe add "location" in the handler (for future testing)
type UserHandler struct {
	userStore    UserStore
	sessionStore sessions.SessionStore
}

// TODO: Update this function
func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type header must be set to application/json", http.StatusUnsupportedMediaType)
		return
	}
	var user User
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&user); err != nil {
		http.Error(w, "Could not parse request body", http.StatusBadRequest)
		return
	}
	serverConfig := gateway.ServerConfig
	if len(user.Name) < serverConfig.MinUsernameLen || len(user.Name) > serverConfig.MaxUsernameLen {
		http.Error(w, fmt.Sprintf("Username must be between %d and %d characters long", serverConfig.MinUsernameLen, serverConfig.MaxUsernameLen), http.StatusBadRequest)
		return
	}
	if err := u.userStore.Get(user.Name); err == nil {
		http.Error(w, fmt.Sprintf("User with name %q already exists", user.Name), http.StatusConflict)
		return
	}
	u.userStore.Add(user.Name, user.Password)
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, username string) {
	serverConfig := gateway.ServerConfig
	if err := ValidateUsernameLength(username, serverConfig.MinUsernameLen, serverConfig.MaxUsernameLen); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: Check if user exists
	userExists := true
	if !userExists {
		// TODO: user doesn't exist
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// TODO: extract session
	userResponse := UserSessionResponse{
		Name:    username,
		Session: "",
		Rel:     rel{Self: "/v1/users/" + username},
	}
	ur, err := json.Marshal(userResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ur)
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, username string) {
	serverConfig := gateway.ServerConfig
	if err := ValidateUsernameLength(username, serverConfig.MinUsernameLen, serverConfig.MaxUsernameLen); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: Check if user exists
	userExists := true
	if !userExists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// TODO: extract password & session (as we'll be updating the user's password)
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	// TODO: check if body contains ONLY the password & session
	containsBoth := true
	if !containsBoth {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// If this user account is protected for some reason
	accountEditable := true
	if !accountEditable {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// TODO: verify session is associated with account
	sessionTokenOk := true
	if !sessionTokenOk {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("WWW_Authenticate: Bearer charset=\"UTF-8\""))
		return
	}
	// TODO: update user's password
	userResponse := UserResponse{
		Name: username,
		Rel:  rel{Self: "/v1/users/" + username},
	}
	ur, err := json.Marshal(userResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ur)
}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, username string) {
	serverConfig := gateway.ServerConfig
	if err := ValidateUsernameLength(username, serverConfig.MinUsernameLen, serverConfig.MaxUsernameLen); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: check if the header params contains ONLY the session-token
	containsBoth := true
	if !containsBoth {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// TODO: Check if user exists
	userExists := true
	if !userExists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// If this user account is protected for some reason
	accountEditable := true
	if !accountEditable {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// TODO: verify session is associated with account
	sessionTokenOk := true
	if !sessionTokenOk {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("WWW_Authenticate: Bearer charset=\"UTF-8\""))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (u *UserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	remainingPath := strings.TrimPrefix(r.RequestURI, "/v1/users")
	pathParameters := strings.Split(remainingPath, "/")

	if len(pathParameters) == 0 {
		// /v1/users
		switch r.Method {
		case http.MethodPost:
			u.CreateUser(w, r)
		default:
			http.Error(w, fmt.Sprintf("Unsupported method: %s for /v1/users", r.Method), http.StatusMethodNotAllowed)
		}
	} else if len(pathParameters) == 1 {
		// /v1/users/{username}
		username := pathParameters[0]
		switch r.Method {
		case http.MethodGet:
			u.GetUser(w, r, username)
		case http.MethodPatch:
			u.UpdateUser(w, r, username)
		case http.MethodDelete:
			u.DeleteUser(w, r, username)
		default:
			http.Error(w, fmt.Sprintf("Unsupported method: %s for /v1/users/{username}", r.Method), http.StatusMethodNotAllowed)
		}

	} else if len(pathParameters) == 2 && pathParameters[1] == "sessions" {
		// /v1/users/sessions/{username}
		sessions.NewHandler(u.sessionStore, u.userStore).Handle(w, r)
	}

	http.NotFound(w, r)
}

func NewHandler(userStore UserStore, sessionStore sessions.SessionStore) *UserHandler {
	return &UserHandler{userStore: userStore, sessionStore: sessionStore}
}
