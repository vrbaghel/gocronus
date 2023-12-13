package utils

import (
	"fmt"
	"ncronus/database/mysql/models"
	"ncronus/services/types"
	"strings"

	"github.com/volatiletech/null/v8"
)

func NModelToNotification(nModel *models.Notification) types.Notification {
	isIST := nModel.NTimezone == models.NotificationNTimezoneIST
	notification := types.Notification{
		ID:       nModel.ID,
		Title:    nModel.R.IDNotificationDatum.NDTitle.String,
		Body:     nModel.R.IDNotificationDatum.NDBody.String,
		Action:   nModel.NAction,
		Timezone: nModel.NTimezone,
		Device:   nModel.NDevice,
		Status:   nModel.NStatus,
	}
	if isIST {
		notification.ScheduledOn = nModel.NTimestamp.In(types.IST_TIMEZONE).Format(types.NOTIFICATION_TIMESTAMP_FORMAT)
	} else {
		notification.ScheduledOn = nModel.NTimestamp.In(types.CST_TIMEZONE).Format(types.NOTIFICATION_TIMESTAMP_FORMAT)
	}
	return notification
}

func NotificationImgUrlsToNIUModel(imgUrls []string) models.NotificationImgURLSlice {
	nImgUrls := models.NotificationImgURLSlice{}
	for _, url := range imgUrls {
		nImgUrls = append(nImgUrls, &models.NotificationImgURL{
			NiuURL: null.StringFrom(url),
		})
	}
	return nImgUrls
}

func NotificationGifUrlsToNGUModel(gifUrls []string) models.NotificationGifURLSlice {
	nGifUrls := models.NotificationGifURLSlice{}
	for _, url := range gifUrls {
		nGifUrls = append(nGifUrls, &models.NotificationGifURL{
			NguURL: null.StringFrom(url),
		})
	}
	return nGifUrls
}

func NModelToNotificationReq(nModel *models.Notification, isIos bool) types.RequestNotificationPayload {
	reqPayload := types.RequestNotificationPayload{
		To:             fmt.Sprintf("/topics/%s", nModel.NAction),
		MutableContent: true,
		Data: types.RequestNotificationDataPayload{
			Id:          nModel.R.IDNotificationDatum.NDUUID,
			Title:       nModel.R.IDNotificationDatum.NDTitle.String,
			Body:        nModel.R.IDNotificationDatum.NDBody.String,
			Source:      nModel.R.IDNotificationDatum.NDSource,
			Category:    types.NDCATEGORY_TO_CATEGORY_MAP[nModel.R.IDNotificationDatum.NDCategory],
			NavType:     types.NDNAVTYPE_TO_NAVTYPE_MAP[nModel.R.IDNotificationDatum.NDNavtype],
			ImageUrls:   strings.Join(NDImgUrlsToUrls(nModel.R.IDNotificationDatum.R.NDNotificationImgUrls)[:], ","),
			GifUrls:     strings.Join(NDGifUrlsToUrls(nModel.R.IDNotificationDatum.R.NDNotificationGifUrls)[:], ","),
			PackageId:   nModel.R.IDNotificationDatum.R.IDNotificationPack.NPID.String,
			PackageName: nModel.R.IDNotificationDatum.R.IDNotificationPack.NPName.String,
			OrderId:     nModel.R.IDNotificationDatum.R.IDNotificationPack.NPOrderID.String,
			FilterId:    nModel.R.IDNotificationDatum.R.IDNotificationPack.NPFilterID.String,
			ToolId:      nModel.R.IDNotificationDatum.R.IDNotificationPack.NPToolID.String,
		},
	}
	if isIos {
		reqPayload.Notification = types.RequestNotificationAdditionalPayload{
			Title:       nModel.R.IDNotificationDatum.NDTitle.String,
			Body:        nModel.R.IDNotificationDatum.NDBody.String,
			ClickAction: nModel.R.IDNotificationDatum.NDClickAction.String,
		}
	}
	return reqPayload
}

func NDImgUrlsToUrls(ndImgUrls models.NotificationImgURLSlice) []string {
	imgUrls := []string{}
	for _, url := range ndImgUrls {
		imgUrls = append(imgUrls, url.NiuURL.String)
	}
	return imgUrls
}

func NDGifUrlsToUrls(ndGifUrls models.NotificationGifURLSlice) []string {
	imgUrls := []string{}
	for _, url := range ndGifUrls {
		imgUrls = append(imgUrls, url.NguURL.String)
	}
	return imgUrls
}
