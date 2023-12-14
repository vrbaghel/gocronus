package handler

import (
	"fmt"
	"ncronus/pkg/auth"
	"ncronus/services/types"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthMiddleware(ctx *gin.Context) {
	metadata, err := h.ValidateToken(ctx.Request)
	if err != nil {
		h.logger.Error(fmt.Sprintf("auth handler : token validation failed :: unable to validate authorization header %s", err.Message))
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}
	if metadata.AccessLevel != types.ACCESS_SUPER && metadata.AccessLevel != types.ACCESS_LIMITED {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, &auth.AuthError{
			Code:    http.StatusUnauthorized,
			Status:  http.StatusText(http.StatusUnauthorized),
			Message: "Access Denied",
		})
	}
	ctx.Next()
}

func (h *Handler) ValidateToken(request *http.Request) (*auth.JWTMetadata, *auth.AuthError) {
	clientToken := request.Header.Get("Authorization")
	if clientToken == "" {
		return nil, &auth.AuthError{
			Code:    http.StatusUnauthorized,
			Status:  http.StatusText(http.StatusUnauthorized),
			Message: "error missing token in header",
		}
	}
	extractedToken := strings.Split(clientToken, "Bearer ")
	if len(extractedToken) == 2 {
		clientToken = strings.TrimSpace(extractedToken[1])
	} else {
		return nil, &auth.AuthError{
			Code:    http.StatusUnauthorized,
			Status:  http.StatusText(http.StatusUnauthorized),
			Message: "error invalid Authorization header",
		}
	}
	meta, err := auth.ValidateToken(clientToken, false, h.config.API.Token)
	if err != nil {
		return nil, &auth.AuthError{
			Code:    err.Code,
			Status:  err.Status,
			Message: err.Message,
		}
	}
	return meta, err
}
