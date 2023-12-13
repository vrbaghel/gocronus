package types

import (
	"ncronus/database/mysql/models"
	"time"
)

const NOTIFICATION_TIMESTAMP_FORMAT = "2006-01-02 15:04:05"

var IST_TIMEZONE = time.FixedZone("IST", int(5.5*time.Hour.Hours()*time.Hour.Seconds()))
var CST_TIMEZONE = time.FixedZone("CST", int(-6*time.Hour.Hours()*time.Hour.Seconds()))

var NDCATEGORY_TO_CATEGORY_MAP map[string]int = map[string]int{
	models.NotificationDataNDCategoryText:           0,
	models.NotificationDataNDCategoryCarousel:       1,
	models.NotificationDataNDCategoryThumbnailImage: 2,
	models.NotificationDataNDCategoryGif:            3,
}

var NDNAVTYPE_TO_NAVTYPE_MAP map[string]int = map[string]int{
	models.NotificationDataNDNavtypeAiTool:     0,
	models.NotificationDataNDNavtypeAiFilter:   1,
	models.NotificationDataNDNavtypeAiPhoto:    2,
	models.NotificationDataNDNavtypeProfile:    3,
	models.NotificationDataNDNavtypePackDetail: 4,
}

const NOTIFICATION_ID_PARAM = "notificationID"
