package server

import (
	"ncronus/services/api/handler"
	"ncronus/services/store"

	"github.com/gin-gonic/gin"
)

func RegisterEndpoints(routineEngine *gin.Engine, handlerParams handler.HandlerParams) {
	store := store.NewStore()
	handler := handler.NewHandler(handlerParams, store)
	routineEngine.GET("/", handler.ServeBaseRequest)

	api := routineEngine.Group("api")
	// login API
	// api.POST("/user/login", handler.Login)

	notification := api.Group("notification")
	// get actions
	notification.GET("/actions", handler.GetNotiticationActions)
	// get notifications
	notification.GET("/all", handler.GetNotifications)
	// send notification
	notification.POST("/send", handler.SendNotification)
}
