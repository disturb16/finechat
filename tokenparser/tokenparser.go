package tokenparser

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AuthClaims is the JWT claims for the auth token.
type AuthClaims struct {
	*jwt.StandardClaims
	Email string `json:"email"`
	Name  string `json:"name"`
}

const expireMinutes time.Duration = 20

var (
	signKey                *rsa.PrivateKey
	ErrTokenNotValid       error = errors.New("Token is not valid")
	ErrTokenClaimsNotValid error = errors.New("Token claims is not valid")
)

func init() {
	var err error
	signKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
}

// CreateAuthToken returns a new auth token.
func CreateAuthToken(email, name string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	t.Claims = &AuthClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * expireMinutes).Unix(),
		},
		email,
		name,
	}

	return t.SignedString(signKey)
}

// VerifyAuthToken verifies the auth token.
func VerifyAuthToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, parseCallback)
	if err != nil {
		return err
	}

	if !token.Valid {
		return ErrTokenNotValid
	}

	return nil
}

// parseCallback is the callback function for the jwt.Parse.
func parseCallback(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
	}

	return signKey.Public(), nil
}
