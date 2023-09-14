package dto

type NonceResponse struct {
	Nonce  string `json:"nonce"`
	UserID string `json:"user_id"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SignInResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type SignInRequest struct {
	Signature string `json:"signature"`
	Address   string `json:"address"`
	ChainID   int    `json:"chainId"`
	IssuedAt  int64  `json:"issuedAt"`
	ExpiresAt int64  `json:"expiresAt"`
	UserID    string `json:"user_id"`
}
