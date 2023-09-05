package main

import (
	"github.com/avag-sargsyan/golang-clean-arch/internal/adapter/controller"
	"github.com/avag-sargsyan/golang-clean-arch/internal/adapter/repository"
	"github.com/avag-sargsyan/golang-clean-arch/internal/usecase/usecase"
	"net/http"
)

func main() {
	repo := repository.NewInMemoryUserRepository()
	service := usecase.NewUserService(repo)
	handler := controller.NewUserHandler(service)

	// Can be used Google Wire for DI if we have a lot of handlers
	http.HandleFunc("/users", handler.GetUsers)
	http.HandleFunc("/user/create", handler.CreateUser)
	http.ListenAndServe(":8080", nil)
}
