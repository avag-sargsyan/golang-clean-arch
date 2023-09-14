package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/avag-sargsyan/golang-clean-arch/internal/domain/dto"
	"github.com/avag-sargsyan/golang-clean-arch/pkg/uuid"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"regexp"
	"strconv"
	"strings"
)

type AuthService interface {
	Verify(req *dto.SignInRequest) bool
	GetNonce() (string, string, error)
}

type authService struct {
	nonces map[string][]byte
}

func NewAuthService() AuthService {
	n := make(map[string][]byte)
	return &authService{nonces: n}
}

func (s *authService) Verify(req *dto.SignInRequest) bool {
	// TODO: make struct for parsed message
	nonce, ok := s.nonces[req.UserID]
	if !ok {
		fmt.Println("Error getting nonce")
		return false
	}
	message := constructMessage(req.Address, string(nonce), req.ChainID, req.IssuedAt, req.ExpiresAt)
	fmt.Println("Message: ", message)

	return s.verifySignature(req.Address, req.Signature, message)
}

func (s *authService) verifySignature(address, signature, message string) bool {
	sig := hexutil.MustDecode(signature)

	expectedMsg := accounts.TextHash([]byte(message))
	sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	recovered, err := crypto.SigToPub(expectedMsg, sig)
	if err != nil {
		fmt.Println("Error recovering public key: ", err)
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)
	fmt.Println(recoveredAddr.Hex())
	fmt.Println("Address: ", address)
	return address == recoveredAddr.Hex()
}

func (s *authService) GetNonce() (string, string, error) {
	nonce, err := generateNonce()
	if err != nil {
		fmt.Println("Error generating nonce: ", err)
		return "", "", err
	}
	userID, err := uuid.Generate()
	if err != nil {
		fmt.Println("Error generating uuid: ", err)
		return "", "", err
	}
	s.nonces[userID] = []byte(nonce)
	return nonce, userID, nil
}

func generateNonce() (string, error) {
	nonce := make([]byte, 16)
	_, err := rand.Read(nonce)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(nonce), nil
}

func parseMessage(message string) (map[string]string, error) {
	if message == "" {
		return nil, errors.New("message is empty")
	}
	fmt.Println(message)

	result := make(map[string]string)

	fieldExtractors := map[string]*regexp.Regexp{
		"Address":   regexp.MustCompile(`wants you to sign in with your Ethereum account: ([0-9a-zA-Z]+)`),
		"Nonce":     regexp.MustCompile(`Nonce: ([0-9a-zA-Z]+)`),
		"ChainID":   regexp.MustCompile(`Chain ID: ([0-9]+)`),
		"IssuedAt":  regexp.MustCompile(`Issued At: ([0-9]+)`),
		"ExpiresAt": regexp.MustCompile(`Expires At: ([0-9]+)`),
	}

	for key, re := range fieldExtractors {
		field := extractField(re, message)
		if field == "" {
			return nil, fmt.Errorf("failed to extract %s", key)
		}
		result[key] = field
	}

	fmt.Println(result)

	return result, nil
}

func extractField(re *regexp.Regexp, message string) string {
	match := re.FindStringSubmatch(message)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func constructMessage(address string, nonce string, chainID int, issuedAt int64, expiresAt int64) string {
	message := fmt.Sprintf(
		"example.com wants to sign in with your Ethereum account:\n%s\n\n"+
			"By signing this message, you agree to the terms of use and privacy policy of example.com.\n"+
			"URI: https://example.com/auth/login\n"+
			"Version: 1\n"+
			"Nonce: %s\n"+
			"Chain ID: %s\n"+
			"Issued At: %s\n"+
			"Expires At: %s\n",
		address,
		nonce,
		strconv.Itoa(chainID),
		strconv.FormatInt(issuedAt, 10),
		strconv.FormatInt(expiresAt, 10),
	)
	return message
}
