package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// decode -> payload -> signature
type Service interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type JwtService struct {
}

func NewService() *JwtService {
	return &JwtService{}
}

// ini seharusnya di simpan di env
var SECRET_KREY = []byte("S3CR3T_K3Y_JWT")

// generate token
func (s *JwtService) GenerateToken(userID int) (string, error) {

	claim := jwt.MapClaims{}
	// userid -> payload
	claim["user_id"] = userID

	// token berubah menjadi jwt.SigningMethodHS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// signature seperti ttd yang tidak bisa diubah
	signerToken, err := token.SignedString(SECRET_KREY)
	if err != nil {
		return signerToken, err
	}

	return signerToken, nil

}

func (s *JwtService) ValidateToken(token string) (*jwt.Token, error) {

	tokenValidate, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(SECRET_KREY), nil

	})

	if err != nil {
		return tokenValidate, err
	}

	return tokenValidate, nil
}
