package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/volovikariel/IdentityManager/internal/server/gateway"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/models"
	"github.com/volovikariel/IdentityManager/internal/server/gateway/v1/users/sessions"
)

type UserHandler struct {
	userStore    models.UserStore
	sessionStore models.SessionStore
}

// TODO: Update this function
func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /v1/users hit")
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var user models.User
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("Could not parse request body: %v\n", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	username := user.Name
	password := user.Password

	serverConfig := gateway.ServerConfig
	if err := models.ValidateStringLength(username, "username", serverConfig.MinUsernameLen, serverConfig.MaxUsernameLen); err != nil {
		// Username is too short or too long
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := u.userStore.Get(username); err == nil {
		// User already exists
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err := models.ValidateStringLength(password, "password", serverConfig.MinPasswordLen, serverConfig.MaxPasswordLen); err != nil {
		// Password is too short or too long
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := u.userStore.Add(username, password); err != nil {
		// Create user failed
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userResponse := models.UserResponse{
		Name: username,
		Rel:  models.Rel{Self: "/v1/users/" + username},
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

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, username string) {
	serverConfig := gateway.ServerConfig
	if err := models.ValidateStringLength(username, "username", serverConfig.MinUsernameLen, serverConfig.MaxUsernameLen); err != nil {
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
	userResponse := models.UserSessionResponse{
		Name:    username,
		Session: "",
		Rel:     models.Rel{Self: "/v1/users/" + username},
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
	if err := models.ValidateStringLength(username, "username", serverConfig.MinUsernameLen, serverConfig.MaxUsernameLen); err != nil {
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
	userResponse := models.UserResponse{
		Name: username,
		Rel:  models.Rel{Self: "/v1/users/" + username},
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
	if err := models.ValidateStringLength(username, "username", serverConfig.MinUsernameLen, serverConfig.MaxUsernameLen); err != nil {
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
	log.Printf("/v1/users hit with full path %q\n", r.RequestURI)
	if remainingPath == "" {
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
	} else {
		http.NotFound(w, r)
	}
}

func NewHandler(userStore models.UserStore, sessionStore models.SessionStore) *UserHandler {
	return &UserHandler{userStore: userStore, sessionStore: sessionStore}
}
