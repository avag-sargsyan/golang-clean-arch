package usecase

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type AuthService interface {
	VerifySignature(from, sigHex string, expectedMsg []byte) bool
}

type authService struct {
}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) VerifySignature(from, sigHex string, expectedMsg []byte) bool {
	sig := hexutil.MustDecode(sigHex)

	expectedMsg = accounts.TextHash(expectedMsg)
	sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	recovered, err := crypto.SigToPub(expectedMsg, sig)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)

	return from == recoveredAddr.Hex()
}
