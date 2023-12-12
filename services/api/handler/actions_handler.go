package handler

import (
	"net/http"

	"ncronus/services/types"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetNotiticationActions(gCtx *gin.Context) {
	gCtx.JSON(http.StatusOK, types.NotificationActions{
		Actions: h.config.Notification.Actions,
	})
}
