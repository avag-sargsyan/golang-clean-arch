package controller

import (
	"github.com/avag-sargsyan/golang-clean-arch/internal/usecase/usecase"
	"net/http"
)

type Auth interface {
	Verify(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	service usecase.AuthService
}

func NewAuthHandler(s usecase.AuthService) Auth {
	return &authHandler{service: s}
}

func (h *authHandler) Verify(w http.ResponseWriter, r *http.Request) {
	from := "asd"
	sigHex := "0x34850b7e36e635783df0563c7202c3ac776df59db5015d2" +
		"b6f0add33955bb5c43ce35efb5ce695a243bc4c5dc4298db4" +
		"0cd765f3ea5612d2d57da1e4933b2f201b"
	expected := []byte("Example `personal_sign` message")
	isValid := h.service.VerifySignature(from, sigHex, expected)
	if isValid {
		w.Write([]byte("User verified successfully"))
	} else {
		w.Write([]byte("User not verified"))
	}
}
