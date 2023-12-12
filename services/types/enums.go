package types

import "time"

const (
	TIMEOUT_TRANSACTION_SHORT time.Duration = 1 * time.Minute
	TIMEOUT_TRANSACTION_LONG  time.Duration = 2 * time.Minute
)

const (
	NOTIFICATIONS_PER_REQ int = 20
	PAGE_MIN              int = 1
)

const (
	NOTIFICATION_STATUS_SCHEDULED  = "scheduled"
	NOTIFICATION_STATUS_RUNNING    = "running"
	NOTIFICATION_STATUS_COMPLETED  = "completed"
	NOTIFICATION_STATUS_TERMINATED = "terminated"
)

const (
	NOTIFICATION_TIMEZONE_IST = "IST"
	NOTIFICATION_TIMEZONE_GMT = "GMT"
)

const (
	DEVICE_IOS     = "ios"
	DEVICE_ANDROID = "android"
)

const (
	CATEGORY_TEXT            = "text"
	CATEGORY_CAROUSEL        = "carousel"
	CATEGORY_THUMBNAIL_IMAGE = "thumbnail_image"
	CATEGORY_GIF             = "gif"
)

const (
	NAVTYPE_AITOOL      = "ai_tool"
	NAVTYPE_AIFILTER    = "ai_filter"
	NAVTYPE_AIPHOTO     = "ai_photo"
	NAVTYPE_PROFILE     = "profile"
	NAVTYPE_PACK_DETAIL = "pack_detail"
)
