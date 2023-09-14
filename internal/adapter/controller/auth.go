package controller

import (
	"encoding/json"
	"errors"
	"github.com/avag-sargsyan/golang-clean-arch/internal/domain/dto"
	"github.com/avag-sargsyan/golang-clean-arch/internal/usecase/usecase"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
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

	if err := ValidateSignInRequest(payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	nonce, userID, err := h.service.GetNonce()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&dto.ErrorResponse{Message: "Error"})
		return
	}
	json.NewEncoder(w).Encode(&dto.NonceResponse{Nonce: nonce, UserID: userID})
}

func ValidateSignInRequest(req *dto.SignInRequest) error {
	// Check if address is a valid Ethereum address
	if matched, err := regexp.MatchString("^0x[a-fA-F0-9]{40}$", req.Address); err != nil || !matched {
		return errors.New("invalid Ethereum address")
	}

	// Validate Chain ID
	if req.ChainID <= 0 {
		return errors.New("invalid Chain ID")
	}

	// Validate signature (you may have more complex logic here)
	if len(req.Signature) == 0 {
		return errors.New("signature is required")
	}

	// Validate Issued At and Expires At
	if req.IssuedAt > req.ExpiresAt {
		return errors.New("issuedAt timestamp is later than expiresAt")
	}

	// Check for expiry
	currentTime := time.Now().Unix()
	if currentTime > req.ExpiresAt {
		return errors.New("the request has expired")
	}

	return nil
}
