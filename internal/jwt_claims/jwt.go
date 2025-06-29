package jwt

import (
	"crypto/rsa"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

var ErrInvalidSigningMethod = errors.New("invalid signing method")

//go:embed keys/private_key.pem
var privateKey []byte

//go:embed keys/public_key.pem
var publicKey []byte

func New() *Claims {
	tokenTime := time.Now().Add(time.Hour * 24)

	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func ReadPrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA Private Key")
	}

	return privateKey, nil
}

func ReadPublicKey() (*rsa.PublicKey, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA Public Key")
	}

	return publicKey, nil
}

func (c *Claims) GenerateToken(secret *rsa.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)

	resultToken, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return resultToken, nil
}

func (c *Claims) ValidateToken(resultToken string, secret *rsa.PublicKey) error {
	token, err := jwt.ParseWithClaims(resultToken, c, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodRSA); !ok || method != jwt.SigningMethodRS256 {
			return nil, ErrInvalidSigningMethod
		}

		return secret, nil
	})
	if err != nil {
		return fmt.Errorf("failed to parse token")
	}

	if !token.Valid {
		return errors.New("token is not valid")
	}

	return nil
}

func (c *Claims) GetPublicKey() (*rsa.PublicKey, error) {
	publicKey, err := ReadPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get public key")
	}

	return publicKey, nil
}
