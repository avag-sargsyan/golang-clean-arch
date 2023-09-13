package dto

type NonceResponse struct {
	Nonce string `json:"nonce"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SignInResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type SignInRequest struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}
