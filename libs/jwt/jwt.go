package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Header is json web token header
type Header struct {
	Algorithm Algorithm `json:"alg"`
	Type      string    `json:"typ"`
}

// Payload is json web token payload
type Payload map[string]interface{}

// Token interface
type Token interface {
	SetSubject(sub string)
	GetSubject() (string, error)
}

// Algorithm type
type Algorithm string

// Algorithms
const (
	HS256 Algorithm = "HS256"
)

// Registered Claim Names
const (
	ClaimSubject = "sub"
)

var (
	ErrClaimNotFound         = errors.New("claim not found")
	ErrInvalidClaimType      = errors.New("invalid claim type")
	ErrInvalidTokenSignature = errors.New("invalid token signature")
	ErrUnsupportedAlgorithm  = errors.New("unsupported algorithm")
)

type token struct {
	header  Header
	payload Payload
}

func (t *token) SetSubject(sub string) {
	t.payload[ClaimSubject] = sub
}

func (t *token) GetSubject() (string, error) {
	value, ok := t.payload[ClaimSubject]
	if !ok {
		return "", ErrClaimNotFound
	}

	sub, ok := value.(string)
	if !ok {
		return "", ErrInvalidClaimType
	}

	return sub, nil
}

// New returns new json web token
func New() Token {
	return &token{
		header: Header{
			Algorithm: HS256,
			Type:      "JWT",
		},
		payload: map[string]interface{}{},
	}
}

// Sign the token with secret key
func Sign(t Token, key []byte) (string, error) {
	h := t.(interface{ GetHeader() Header }).GetHeader()
	header, err := json.Marshal(h)
	if err != nil {
		return "", fmt.Errorf("error on marshal header: %s", err.Error())
	}

	payload, err := json.Marshal(t.(interface{ GetPayload() Payload }).GetPayload())
	if err != nil {
		return "", fmt.Errorf("error on marshal payload: %s", err.Error())
	}

	unsignedToken := fmt.Sprintf("%s.%s", base64.RawURLEncoding.EncodeToString(header), base64.RawURLEncoding.EncodeToString(payload))

	switch h.Algorithm {
	case HS256:
		mac := hmac.New(sha256.New, key)
		_, _ = mac.Write([]byte(unsignedToken))
		return fmt.Sprintf("%s.%s", unsignedToken, base64.RawURLEncoding.EncodeToString(mac.Sum(nil))), nil
	default:
		return "", ErrUnsupportedAlgorithm
	}
}

// Verify token string with secret key
func Verify(t string, key []byte) error {
	arr := strings.Split(t, ".")
	if len(arr) != 3 {
		return errors.New("invalid token provided")
	}

	tok := token{}

	header, err := base64.RawURLEncoding.DecodeString(arr[0])
	if err != nil {
		return fmt.Errorf("invalid token header encoding: %s", err.Error())
	}

	err = json.Unmarshal(header, &tok.header)
	if err != nil {
		return fmt.Errorf("invalid token header: %s", err.Error())
	}

	if typ := tok.header.Type; typ != "JWT" {
		return fmt.Errorf("unsupported token type: %s", typ)
	}

	switch tok.header.Algorithm {
	case HS256:
		mac := hmac.New(sha256.New, key)
		_, _ = mac.Write([]byte(fmt.Sprintf("%s.%s", arr[0], arr[1])))
		sig, err := base64.RawURLEncoding.DecodeString(arr[2])
		if err != nil {
			return fmt.Errorf("invalid token signature encoding: %s", err.Error())
		}
		if !hmac.Equal(mac.Sum(nil), sig) {
			return ErrInvalidTokenSignature
		}

		return nil
	default:
		return ErrUnsupportedAlgorithm
	}
}

// Parse token string without verifying
func Parse(t string) (Token, error) {
	arr := strings.Split(t, ".")
	if len(arr) != 3 {
		return nil, errors.New("invalid token provided")
	}

	tok := token{}

	header, err := base64.RawURLEncoding.DecodeString(arr[0])
	if err != nil {
		return nil, fmt.Errorf("invalid token header encoding: %s", err.Error())
	}

	err = json.Unmarshal(header, &tok.header)
	if err != nil {
		return nil, fmt.Errorf("invalid token header: %s", err.Error())
	}

	payload, err := base64.RawURLEncoding.DecodeString(arr[1])
	if err != nil {
		return nil, fmt.Errorf("invalid token payload encoding: %s", err.Error())
	}

	err = json.Unmarshal(payload, &tok.payload)
	if err != nil {
		return nil, fmt.Errorf("invalid token payload: %s", err.Error())
	}

	return &tok, nil
}
