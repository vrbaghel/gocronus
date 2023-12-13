package handler

import (
	"ncronus/database/mysql"
	"ncronus/pkg/auth"
	"ncronus/services/store"
	"ncronus/services/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type HandlerConfig struct {
	API          APIConfig
	Notification NotificationConfig
}

type APIConfig struct {
	Token auth.AuthConfig
}

type NotificationConfig struct {
	BaseURL string
	Actions []string
	AuthKey string
}

type HandlerParams struct {
	Config *HandlerConfig
	MySql  *mysql.MySQL
	Logger *zap.Logger
	Cron   *CronConfig
}

type CronConfig struct {
	CST *cron.Cron
	IST *cron.Cron
}

type Handler struct {
	config  *HandlerConfig
	context map[string]string
	mySql   *mysql.MySQL
	store   *store.Store
	logger  *zap.Logger
	cron    *CronConfig
}

func NewHandler(handlerParams HandlerParams, store *store.Store) *Handler {
	return &Handler{
		config:  handlerParams.Config,
		mySql:   handlerParams.MySql,
		logger:  handlerParams.Logger,
		cron:    handlerParams.Cron,
		context: map[string]string{},
		store:   store,
	}
}

func (h *Handler) ServeBaseRequest(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"request_url": ctx.Request.RequestURI,
		"status":      "OK",
	})
}

func (h *Handler) InternalServerError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, types.APIError{
		Code:    http.StatusInternalServerError,
		Status:  http.StatusText(http.StatusInternalServerError),
		Message: "an unknown error has occured, please try again later",
	})
}
