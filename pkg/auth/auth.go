package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func ValidateToken(token string, isRefresh bool, config AuthConfig) (*JWTMetadata, *AuthError) {
	var secret []byte
	var accessLevel int

	parsedToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if isRefresh {
			secret = config.RefreshTokenSecret
		} else {
			secret = config.AccessTokenSecret
		}
		accessLevelF, ok := t.Header["departmentid"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid token code")
		}
		accessLevel = int(accessLevelF)
		return secret, nil
	})

	if err != nil {
		var message string
		if err == jwt.ErrSignatureInvalid {
			message = "err : invalid signature"
		} else {
			message = "err : invalid token"
		}
		return nil, &AuthError{
			Code:    http.StatusUnauthorized,
			Status:  http.StatusText(http.StatusUnauthorized),
			Message: message,
		}
	}

	if !parsedToken.Valid {
		return nil, &AuthError{
			Code:    http.StatusUnauthorized,
			Status:  http.StatusText(http.StatusUnauthorized),
			Message: "invalid token",
		}
	}

	return &JWTMetadata{
		AccessLevel: accessLevel,
	}, nil
}

func VerifyPassword(plainText []byte, hash []byte) *AuthError {
	if err := bcrypt.CompareHashAndPassword(hash, plainText); err != nil {
		return &AuthError{
			Code:    http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Message: "password is invalid",
		}
	}
	return nil
}

func GenerateTokens(id int, accessLevel int, config AuthConfig) (string, string, *AuthError) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["exp"] = time.Now().Add(config.AccessTokenExpiration * time.Hour).Unix()
	accessClaims["sub"] = id
	accessClaims["acl"] = accessLevel
	accessTokenString, err := accessToken.SignedString(config.AccessTokenSecret)
	if err != nil {
		return "", "", &AuthError{
			Code:    http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: "err : failed to generate tokens",
		}
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["exp"] = time.Now().Add(config.RefreshTokenExpiration * time.Hour).Unix()
	refreshClaims["sub"] = id
	refreshTokenString, err := refreshToken.SignedString(config.RefreshTokenSecret)
	if err != nil {
		return "", "", &AuthError{
			Code:    http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: "err : failed to generate tokens",
		}
	}

	return accessTokenString, refreshTokenString, nil
}
