package main

import (
	"log"
	"ncronus/database/mysql"
	"ncronus/pkg/auth"
	"ncronus/pkg/env"
	"ncronus/pkg/logger"
	"ncronus/services/api/handler"
	"ncronus/services/api/server"
	"ncronus/services/types"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var process *Process

type Process struct {
	env    *env.ENV
	logger *zap.Logger
}

func NewProcess(env *env.ENV) *Process {
	return &Process{
		env: env,
	}
}

func init() {
	env, err := env.LoadENV()
	if err != nil {
		log.Fatal(err.Error())
	}
	process = NewProcess(env)
}

// initialize logger
func (p *Process) initLogger() *zap.Logger {
	p.logger = logger.NewLogger(p.env.SERVER_LOG_LEVEL)
	defer p.logger.Sync()
	return p.logger
}

// initialize mysql client instance
func (p *Process) initMySQL() *mysql.MySQL {
	mySql, err := mysql.NewMySQL(p.env.MYSQL_CONNECTION_STRING)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err = mySql.Ping(); err != nil {
		log.Fatal(err.Error())
	}
	return mySql
}

func (p *Process) closeSQLConnection(client *mysql.MySQL) {
	if client != nil {
		if err := client.CloseConnection(); err != nil {
			log.Fatal(err.Error())
		}
	}
	log.Println("closed sql connection")
}

func (p *Process) initCron() *handler.CronConfig {
	return &handler.CronConfig{
		CST: cron.New(cron.WithLocation(types.CST_TIMEZONE)),
		IST: cron.New(cron.WithLocation(types.IST_TIMEZONE)),
	}
}

func (p *Process) startCron(cronConfig *handler.CronConfig) {
	cronConfig.CST.Start()
	cronConfig.IST.Start()
	log.Println("cron started")
}

func (p *Process) destroyCron(cronConfig *handler.CronConfig) {
	cronConfig.CST.Stop()
	cronConfig.IST.Stop()
	log.Println("cron stopped")
}

// // initialize http server
func (p *Process) initHTTPServer(handlerParams handler.HandlerParams) (func(), func()) {
	serverConfig := server.ServerConfig{
		ServerHTTPPort:            p.env.SERVER_HTTP_PORT,
		ServerHTTPReadTimeout:     time.Duration(p.env.SERVER_TIMEOUT_READ_SECONDS),
		ServerHTTPWriteTimeout:    time.Duration(p.env.SERVER_TIMEOUT_WRITE_SECONDS),
		ServerHTTPShutdownTimeout: time.Duration(p.env.SERVER_TIMEOUT_SHUTDOWN_SECONDS),
	}
	httpServer := server.NewServer(serverConfig)
	server.RegisterEndpoints(httpServer.RoutingEngine, handlerParams)
	return httpServer.StartServer, httpServer.StopServer
}

func main() {
	handlerParams := handler.HandlerParams{}
	// logger config
	handlerParams.Logger = process.initLogger()
	// cron config
	handlerParams.Cron = process.initCron()
	// service config
	handlerParams.Config = &handler.HandlerConfig{
		API: handler.APIConfig{
			Token: auth.AuthConfig{
				AccessTokenSecret:      []byte(process.env.ACCESS_TOKEN_SECRET),
				AccessTokenExpiration:  time.Duration(process.env.ACCESS_TOKEN_EXPIRATION),
				RefreshTokenSecret:     []byte(process.env.REFRESH_TOKEN_SECRET),
				RefreshTokenExpiration: time.Duration(process.env.REFRESH_TOKEN_EXPIRATION),
			},
		},
		Notification: handler.NotificationConfig{
			BaseURL: process.env.FCM_NOTIFICATION_URL,
			Actions: process.env.FCM_NOTIFICATION_ACTIONS,
			AuthKey: process.env.FCM_NOTIFICATION_AUTH_KEY,
		},
	}
	// sql client config
	sqlClient := process.initMySQL()
	handlerParams.MySql = sqlClient
	startHTTPServer, stopHTTPServer := process.initHTTPServer(handlerParams)

	defer process.closeSQLConnection(sqlClient)
	defer process.destroyCron(handlerParams.Cron)
	defer stopHTTPServer()
	process.startCron(handlerParams.Cron)
	startHTTPServer()
}
