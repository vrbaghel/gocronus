package types

import "time"

const (
	TIMEOUT_TRANSACTION_SHORT time.Duration = 1 * time.Minute
	TIMEOUT_TRANSACTION_LONG  time.Duration = 2 * time.Minute
)

const (
	NOTIFICATIONS_PER_REQ int = 10
	PAGE_MIN              int = 1
)

const (
	ACCESS_SUPER   = 1
	ACCESS_LIMITED = 5
)

const (
	NOTIFICATION_ENV_PROD = 0
	NOTIFICATION_ENV_DEV  = 1
)
