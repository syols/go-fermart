package authorizer

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/docker/distribution/registry/auth"
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/models"
)

type UserClaims struct {
	jwt.StandardClaims
	Username string
}

type Authorizer struct {
	sign []byte
}

func NewAuthorizer(config config.Config) Authorizer {
	return Authorizer{
		sign: []byte(config.Sign),
	}
}

func (a Authorizer) CreateToken(user models.User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Hour).Unix(),
			Issuer:    user.Username,
		},
		Username: user.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.sign)
}

func (a Authorizer) VerifyToken(token string) (string, error) {
	var claims UserClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrAuthenticationFailure
		}
		return a.sign, nil
	})
	return claims.Username, err
}
