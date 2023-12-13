package server

import (
	"ncronus/services/api/handler"
	"ncronus/services/store"

	"github.com/gin-gonic/gin"
)

func RegisterEndpoints(routineEngine *gin.Engine, handlerParams handler.HandlerParams) {
	store := store.NewStore()
	handler := handler.NewHandler(handlerParams, store)
	// restart pending CRON jobs
	handler.RestartCronJobs()

	routineEngine.GET("/", handler.ServeBaseRequest)
	api := routineEngine.Group("api")
	notification := api.Group("notification", handler.AuthMiddleware)
	// get actions
	notification.GET("/actions", handler.GetNotiticationActions)
	// get notifications
	notification.GET("/all", handler.GetNotifications)
	// send notification
	notification.POST("/send", handler.SendNotification)
	// terminate notification
	// notification.DELETE("/")
}
