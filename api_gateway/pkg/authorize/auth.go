package authorize

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

// Payload is the structure of the JWT payload.
type Payload struct {
	UserID  uint64
	IsAdmin bool
	jwt.StandardClaims
}

// GenerateJwt creates a new JWT.
func GenerateJwt(userID uint64, isAdmin bool, secret []byte) (string, error) {
	// Set the expiration time for the token.
	expirationTime := time.Now().Add(48 * time.Hour)

	// Create the JWT claims, which includes the user information and expiry time.
	claims := &Payload{
		UserID:  userID,
		IsAdmin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
       log.Println(claims.Id)
	// Declare the token with the algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates the JWT.
func ValidateToken(tokenString string, secret []byte) (map[string]interface{}, error) {
	// Initialize a new instance of `Payload`
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(t *jwt.Token) (interface{}, error) {

		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("invalid token")
		}

		return secret, nil

	})

	if err != nil {
		return nil, err
	}

	if token == nil || !token.Valid {
		return nil, fmt.Errorf("token is not valid or its empty")
	}

	cliams, ok := token.Claims.(*Payload)

	if !ok {
		return nil, fmt.Errorf("cannot parse claims")
	}

	cred := map[string]interface{}{
		"userID":  cliams.UserID,
		"isadmin": cliams.IsAdmin,
		
	}

	if cliams.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return cred, nil

}