package tokenparser

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

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

func VerifyAuthToken(tokenString string) (*AuthClaims, error) {

	token, err := jwt.Parse(tokenString, parseCallback)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrTokenNotValid
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, ErrTokenClaimsNotValid
	}

	return claims, nil
}

func parseCallback(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
	}

	return signKey.Public(), nil
}
