package controller

import (
	"encoding/json"
	"github.com/avag-sargsyan/golang-clean-arch/internal/domain/dto"
	"github.com/avag-sargsyan/golang-clean-arch/internal/usecase/usecase"
	"io/ioutil"
	"net/http"
)

type Auth interface {
	SignIn(w http.ResponseWriter, r *http.Request)
	GetNonce(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	service usecase.AuthService
}

func NewAuthHandler(s usecase.AuthService) Auth {
	return &authHandler{service: s}
}

func (h *authHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	//from := "0x41491b64Ed1E61580E18AE93FF1cD83f4533f876"
	//sigHex := "0xd15743cb446c29b86b15f6ec02b38023c75ddb2e5cb19b3f458e9c01a7fd0d880ba97e26e94fbda7efdea71ae4aca4a5015e03cf3f478a5cc712df687e54e6171b"
	//expected := []byte("Your unique challenge message here.")

	body, err := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var payload *dto.SignInRequest
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	isValid := h.service.Verify(payload)

	var response *dto.SignInResponse

	if isValid {
		response = &dto.SignInResponse{
			Message: "User verified successfully",
			Success: true,
		}
	} else {
		response = &dto.SignInResponse{
			Message: "User not verified",
			Success: false,
		}
	}

	json.NewEncoder(w).Encode(response)
}

func (h *authHandler) GetNonce(w http.ResponseWriter, r *http.Request) {
	nonce, err := h.service.GetNonce()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Message: "Error"})
		return
	}
	json.NewEncoder(w).Encode(&dto.NonceResponse{Nonce: nonce})
}
