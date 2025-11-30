package helpers

//TODO: Configure

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ishu17077/code_runner_backend/worker-node/constants"
)

type SignedDetails struct {
	Username string
	jwt.RegisteredClaims
}

func GenerateTokens(username string) (signedToken, refreshToken string, err error) {
	claims := &SignedDetails{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MrX",
			Subject:   "Admin Authentication",
			IssuedAt:  jwt.NewNumericDate(time.Now().Local()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 1)),
		},
	}

	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(constants.SecretKey))
	if err != nil {
		panic(err)

	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(constants.SecretKey))
	if err != nil {
		panic(err)
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, err error) {
	token, err := checkToken(signedToken)

	if err != nil {
		// if errors.Is(err, jwt.ErrTokenExpired) {
		// 	// refreshToken, err := checkToken(signedRefreshToken)
		// 	// if err != nil || token == nil || !token.Valid {
		// 	// 	return nil, fmt.Errorf("The token has expired, please sign in again")
		// 	// }
		// 	// refreshClaims, ok := token.Claims.(*SignedDetails)
		// 	// if !ok {
		// 	// 	return nil, fmt.Errorf("The token is invalid")
		// 	// }
		// 	return nil, fmt.Errorf("The token is invalid")
		// } else if errors.Is(err, jwt.ErrTokenNotValidYet) {
		// 	return nil, fmt.Errorf("The token is not valid yet")
		// }
		return nil, err
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok || !token.Valid {
		msg := "The token is invalid"
		return nil, fmt.Errorf(msg)
	}
	return claims, err
}

func checkToken(signedToken string) (token *jwt.Token, err error) {
	token, err = jwt.ParseWithClaims(signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (any, error) {
			if token.Method.Alg() != jwt.SigningMethodES256.Alg() {
				return nil, fmt.Errorf("Unexpected Signing Method")
			}
			return []byte(constants.SecretKey), nil
		})
	return
}
