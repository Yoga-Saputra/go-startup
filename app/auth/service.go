package auth

import (
	"bwastartup/app/helper"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userId int) (string, error)
}

type jwtService struct{}

func NewService() *jwtService {
	return &jwtService{}
}

var KEY = helper.GetInv("SECRET_KEY")
var SECRET_KEY = []byte(KEY)

func (s *jwtService) GenerateToken(userId int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil

}
