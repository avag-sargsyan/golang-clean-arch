package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/avag-sargsyan/golang-clean-arch/internal/domain/dto"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"regexp"
	"strings"
)

type AuthService interface {
	Verify(req *dto.SignInRequest) bool
	GetNonce() (string, error)
}

type authService struct {
	nonce []byte
}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Verify(req *dto.SignInRequest) bool {
	// TODO: make struct for parsed message
	parsedMessage, err := parseMessage(req.Message)
	if err != nil {
		fmt.Println("Error parsing message: ", err)
		return false
	}
	return s.verifySignature(parsedMessage["Address"], req.Signature)
}

func (s *authService) verifySignature(address, signature string) bool {
	sig := hexutil.MustDecode(signature)

	expectedMsg := accounts.TextHash(s.nonce)
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

func (s *authService) GetNonce() (string, error) {
	nonce, err := generateNonce()
	if err != nil {
		fmt.Println("Error generating nonce: ", err)
		return "", err
	}
	s.nonce = []byte(nonce)
	return nonce, nil
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

	result := make(map[string]string)

	fieldExtractors := map[string]*regexp.Regexp{
		"Address":   regexp.MustCompile(`Address: ([0-9a-zA-Z]+)`),
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
