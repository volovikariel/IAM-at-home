package users

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: Maybe add "location" in the handler (for future testing)
type UserHandler struct {
	Store UserStore
}

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
	if err := u.Store.Get(user.Name); err == nil {
		http.Error(w, fmt.Sprintf("User with name %q already exists", user.Name), http.StatusConflict)
		return
	}
	u.Store.Add(user.Name, user.Password)
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func NewHandler(store UserStore) *UserHandler {
	return &UserHandler{Store: store}
}
