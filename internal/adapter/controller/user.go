package controller

import (
	"encoding/json"
	"github.com/avag-sargsyan/golang-clean-arch/internal/usecase/usecase"
	"net/http"
)

type User interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	service usecase.UserService
}

func NewUserHandler(s usecase.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, _ := h.service.GetUsers()
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	h.service.CreateUser(name)
	w.Write([]byte("User created successfully"))
}
