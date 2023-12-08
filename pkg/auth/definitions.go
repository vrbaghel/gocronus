package auth

import "time"

type JWT struct {
	AccessTokenString  string `json:"token,omitempty"`
	RefreshTokenString string `json:"refresh_token,omitempty"`
}

type JWTMetadata struct {
	UserID      int
	AccessLevel int
}

type AuthError struct {
	Code    int         `json:"code,omitempty"`
	Status  string      `json:"status,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

type AuthConfig struct {
	AccessTokenSecret      []byte
	AccessTokenExpiration  time.Duration
	RefreshTokenSecret     []byte
	RefreshTokenExpiration time.Duration
}
