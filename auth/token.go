package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shiibs/go-garden-planner/model"
)


type CustomClaims struct {
	Email string
	UserID uint

	jwt.RegisteredClaims
}

var secret = os.Getenv("jwtSecret")

func GenerateToken( user model.User) (string, error) {
	claims := CustomClaims{
		user.Email,
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(30 * 24 * time.Hour)),
		},

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Println("Error in token singing.", err)
		return "", err
	}

	return t, nil
}


// Validate Token
func ValidateToken(clientToken string) (claims *CustomClaims, msg string) {
	token, err := jwt.ParseWithClaims(clientToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*CustomClaims)

	if !ok {
		msg = "Error in claims"
		return
	}

	if !ok {
        msg = "Error in claims or invalid token"
        return
    }

	// Check if the token is expired
    if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
        msg = "Token is expired"
        return
    }

	return claims, msg
}