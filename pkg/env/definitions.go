package env

type ENV struct {
	SERVER_LOG_LEVEL                int
	MYSQL_CONNECTION_STRING         string
	MYSQL_CONNECTION_MAX_OPEN       int
	MYSQL_CONNECTION_MAX_IDLE       int
	SERVER_HTTP_PORT                int
	SERVER_TIMEOUT_READ_SECONDS     int
	SERVER_TIMEOUT_WRITE_SECONDS    int
	SERVER_TIMEOUT_SHUTDOWN_SECONDS int
	ACCESS_TOKEN_SECRET             string
	REFRESH_TOKEN_SECRET            string
	ACCESS_TOKEN_EXPIRATION         int
	REFRESH_TOKEN_EXPIRATION        int
	INFERENCE_BASE_URL              string
	MALE_MODELS_ADULT               []string
	FEMALE_MODELS_ADULT             []string
	MALE_MODELS_KID                 []string
	FEMALE_MODELS_KID               []string
}
