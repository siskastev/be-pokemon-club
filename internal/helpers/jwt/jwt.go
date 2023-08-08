package jwt

import (
	"be-pokemon-club/internal/models"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte(os.Getenv("JWT_PRIVATE_KEY"))

// Claims represents the custom claims for the JWT token
type Claims struct {
	Data models.UserResponse `json:"data"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token for a given email
func GenerateJWT(user models.UserResponse) (string, error) {
	claims := Claims{
		Data: user,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Email,                           // Use the user Email as the subject
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func VerifyAndExtractUserFromJWT(tokenString string) (*models.UserResponse, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, errors.New("failed to parse JWT token: " + err.Error())
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &claims.Data, nil
	}

	return nil, errors.New("unauthorized")
}
