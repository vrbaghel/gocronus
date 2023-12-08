package env

import (
	"os"
	"strconv"
	"strings"
)

func readENV() (*ENV, error) {
	var errSlice []error

	serverLogLevel, errServerLogLevel := strconv.Atoi(os.Getenv("SERVER_LOG_LEVEL"))
	mysqlConnectionString := os.Getenv("MYSQL_CONNECTION_STRING")
	mysqlConnectionMaxOpen, errMySQLConnectionMaxOpen := strconv.Atoi(os.Getenv("MYSQL_CONNECTION_MAX_OPEN"))
	mysqlConnectionMaxIdle, errMySQLConnectionMaxIdle := strconv.Atoi(os.Getenv("MYSQL_CONNECTION_MAX_IDLE"))

	httpServerPort, errHttpServerPort := strconv.Atoi(os.Getenv("SERVER_HTTP_PORT"))
	serverTimeoutReadSeconds, errServerTimeoutReadSeconds := strconv.Atoi(os.Getenv("SERVER_TIMEOUT_READ_SECONDS"))
	serverTimeoutWriteSeconds, errServerTimeoutWriteSeconds := strconv.Atoi(os.Getenv("SERVER_TIMEOUT_WRITE_SECONDS"))
	serverTimeoutShutdownSeconds, errServerTimeoutShutdownSeconds := strconv.Atoi(os.Getenv("SERVER_TIMEOUT_SHUTDOWN_SECONDS"))

	atSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	atExpiration, errAtExpiration := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRATION"))
	rtSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	rtExpiration, errRtExpiration := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRATION"))

	infBaseURL := os.Getenv("INFERENCE_BASE_URL")
	maleModelsA := strings.Split(os.Getenv("MALE_MODELS_ADULT"), ",")
	femaleModelsA := strings.Split(os.Getenv("FEMALE_MODELS_ADULT"), ",")
	maleModelsK := strings.Split(os.Getenv("MALE_MODELS_KID"), ",")
	femaleModelsK := strings.Split(os.Getenv("FEMALE_MODELS_KID"), ",")

	errSlice = append(errSlice, errServerLogLevel, errMySQLConnectionMaxOpen, errMySQLConnectionMaxIdle, errHttpServerPort, errServerTimeoutReadSeconds, errServerTimeoutWriteSeconds, errServerTimeoutShutdownSeconds, errAtExpiration, errRtExpiration)

	if err := checkENVErrors(errSlice); err != nil {
		return nil, err
	}

	return &ENV{
		SERVER_LOG_LEVEL:                serverLogLevel,
		MYSQL_CONNECTION_STRING:         mysqlConnectionString,
		MYSQL_CONNECTION_MAX_OPEN:       mysqlConnectionMaxOpen,
		MYSQL_CONNECTION_MAX_IDLE:       mysqlConnectionMaxIdle,
		SERVER_HTTP_PORT:                httpServerPort,
		SERVER_TIMEOUT_READ_SECONDS:     serverTimeoutReadSeconds,
		SERVER_TIMEOUT_WRITE_SECONDS:    serverTimeoutWriteSeconds,
		SERVER_TIMEOUT_SHUTDOWN_SECONDS: serverTimeoutShutdownSeconds,
		ACCESS_TOKEN_SECRET:             atSecret,
		ACCESS_TOKEN_EXPIRATION:         atExpiration,
		REFRESH_TOKEN_SECRET:            rtSecret,
		REFRESH_TOKEN_EXPIRATION:        rtExpiration,
		INFERENCE_BASE_URL:              infBaseURL,
		MALE_MODELS_ADULT:               maleModelsA,
		FEMALE_MODELS_ADULT:             femaleModelsA,
		MALE_MODELS_KID:                 maleModelsK,
		FEMALE_MODELS_KID:               femaleModelsK,
	}, nil
}
