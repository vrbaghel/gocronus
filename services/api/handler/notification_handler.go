package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"ncronus/database/mysql/models"
	"ncronus/services/api/utils"
	"ncronus/services/types"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (h *Handler) GetNotifications(gCtx *gin.Context) {
	notifications := []types.Notification{}
	totalPages := 1

	txCtx, cancel := context.WithTimeout(context.Background(), types.TIMEOUT_TRANSACTION_SHORT)
	defer cancel()
	tx, err := h.mySql.Client.BeginTx(gCtx, nil)
	if err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : GetNotifications :: Unable to begin sql transaction for request URL %s\t%s", gCtx.Request.URL.String(), err.Error()))
		h.InternalServerError(gCtx)
		return
	}

	qMods := []qm.QueryMod{qm.Load(models.NotificationRels.IDNotificationDatum), qm.Limit(types.NOTIFICATIONS_PER_REQ)}

	// find count of total pages
	totalNotifications, err := h.store.NotificationStore.Count(txCtx, tx, qMods...)
	if err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : GetNotifications :: Unable to get notifications count %s", err.Error()))
		totalNotifications = 0
	}

	totalPages = int(math.Ceil(float64(totalNotifications) / float64(types.NOTIFICATIONS_PER_REQ)))
	if totalNotifications == 0 {
		totalPages = 1
	}

	// do pagination if page param exists
	pageParam := gCtx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : GetNotifications :: Query param parsing failed for page value %s\t%s", pageParam, err.Error()))
		page = 1
	} else if page < types.PAGE_MIN {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	qMods = append(qMods, qm.Offset((page-1)*types.NOTIFICATIONS_PER_REQ))

	// get notifications
	notificationSlice, err := h.store.NotificationStore.GetAll(txCtx, tx, qMods...)
	if err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : GetNotifications :: Unable to get all notifications %s", err.Error()))
		h.InternalServerError(gCtx)
		return
	}

	for _, nModel := range notificationSlice {
		notifications = append(notifications, utils.NModelToNotification(nModel))
	}

	if err := tx.Commit(); err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : GetNotifications :: Unable to commit SQL transaction for request %s\t%s", gCtx.Request.URL.String(), err.Error()))
		h.InternalServerError(gCtx)
		return
	}

	gCtx.JSON(http.StatusOK, types.GetNotificationsResponsePayload{
		Notifications: notifications,
		Pagination: types.PaginationData{
			CurrentPage: page,
			TotalPages:  totalPages,
		},
	})
}

func (h *Handler) SendNotification(gCtx *gin.Context) {
	var reqPayload types.SendNotificationRequestPayload
	// @todo set it to false when impemented CRON
	sendNotificationNow := true
	if err := gCtx.ShouldBindJSON(&reqPayload); err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Unable to bind request body with GenerateImgResponse %s", err.Error()))
		gCtx.JSON(http.StatusBadRequest, types.APIError{
			Code:    http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Message: "invalid payload",
		})
		return
	}
	isIos := reqPayload.Device == types.DEVICE_IOS

	txCtx, cancel := context.WithTimeout(context.Background(), types.TIMEOUT_TRANSACTION_SHORT)
	defer cancel()
	tx, err := h.mySql.Client.BeginTx(gCtx, nil)
	if err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Unable to begin sql transaction for request URL %s\t%s", gCtx.Request.URL.String(), err.Error()))
		h.InternalServerError(gCtx)
		return
	}

	notification := models.Notification{
		NAction: reqPayload.Action,
		NDevice: reqPayload.Device,
		NStatus: types.NOTIFICATION_STATUS_SCHEDULED,
	}
	if reqPayload.Timezone != "" {
		notification.NTimezone = reqPayload.Timezone
	} else {
		notification.NTimezone = types.NOTIFICATION_TIMEZONE_GMT
	}

	if reqPayload.ScheduledFor != "" {
		scheduledFor, err := time.Parse(types.NOTIFICATION_TIMESTAMP_FORMAT, reqPayload.ScheduledFor)
		if err != nil {
			h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Failed to parse timestamp for scheduling notification %s\t%s", reqPayload.ScheduledFor, err.Error()))
		}
		if notification.NTimezone == types.NOTIFICATION_TIMEZONE_IST {
			istSeconds := int(5.5 * time.Hour.Hours() * time.Hour.Seconds())
			istZone := time.FixedZone("IST Time", istSeconds)
			notification.NTimestamp = scheduledFor.In(istZone)
		} else {
			notification.NTimestamp = scheduledFor.UTC()
		}
	} else {
		notification.NTimestamp = time.Now().UTC()
		sendNotificationNow = true
	}
	// h.logger.Info(fmt.Sprintf("NotificationHandler : SendNotification :: scheduledFor %+v", notification.NTimestamp))

	if err := h.store.NotificationStore.Insert(txCtx, tx, &notification); err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Failed to insert notification for request %s\t%s", gCtx.Request.URL.String(), err.Error()))
		h.InternalServerError(gCtx)
		return
	}

	if notification.ID > 0 && (reqPayload.Category != nil || reqPayload.Navigation != nil) {
		nData := models.NotificationDatum{
			NDSource: 0,
		}
		nUuid, err := strconv.Atoi(fmt.Sprintf("%d%d%d%d", notification.ID, notification.ID+1, notification.ID+2, notification.ID+3))
		if err != nil {
			h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Failed to create UUID for notification data for request %s\t%s", gCtx.Request.URL.String(), err.Error()))
			h.InternalServerError(gCtx)
			return
		}
		nData.NDUUID = nUuid

		if reqPayload.Category != nil {
			nData.NDCategory = reqPayload.Category.Type
			nData.NDTitle = null.StringFrom(reqPayload.Category.Data.Title)
			nData.NDBody = null.StringFrom(reqPayload.Category.Data.Body)
		}

		if reqPayload.Navigation != nil {
			nData.NDNavtype = reqPayload.Navigation.Type
		}

		if err := notification.SetIDNotificationDatum(txCtx, tx, true, &nData); err != nil {
			h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Failed to add notification data in notification for request %s\t%s", gCtx.Request.URL.String(), err.Error()))
			h.InternalServerError(gCtx)
			return
		}

		if reqPayload.Category != nil {
			if reqPayload.Category.Data.ImageURLs != nil {
				nImgUrls := utils.NotificationImgUrlsToNIUModel(reqPayload.Category.Data.ImageURLs)
				if err := nData.AddNDNotificationImgUrls(txCtx, tx, true, nImgUrls...); err != nil {
					h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Failed to insert notification image urls for request %s\t%s", gCtx.Request.URL.String(), err.Error()))
					h.InternalServerError(gCtx)
					return
				}
			}
			if reqPayload.Category.Data.GifURLs != nil {
				nImgUrls := utils.NotificationGifUrlsToNGUModel(reqPayload.Category.Data.GifURLs)
				if err := nData.AddNDNotificationGifUrls(txCtx, tx, true, nImgUrls...); err != nil {
					h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Failed to insert notification gif urls for request %s\t%s", gCtx.Request.URL.String(), err.Error()))
					h.InternalServerError(gCtx)
					return
				}
			}
		}

		if reqPayload.Navigation != nil {
			nPackData := models.NotificationPack{
				NPID:       null.StringFrom(reqPayload.Navigation.Data.PackageID),
				NPName:     null.StringFrom(reqPayload.Navigation.Data.PackageName),
				NPOrderID:  null.StringFrom(reqPayload.Navigation.Data.OrderID),
				NPFilterID: null.StringFrom(reqPayload.Navigation.Data.FilterID),
				NPToolID:   null.StringFrom(reqPayload.Navigation.Data.ToolID),
			}

			if err := nData.SetIDNotificationPack(txCtx, tx, true, &nPackData); err != nil {
				h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Failed to insert notification pack data for request %s\t%s", gCtx.Request.URL.String(), err.Error()))
				h.InternalServerError(gCtx)
				return
			}
		}

		if err := h.store.NotificationDataStore.Update(txCtx, tx, &nData); err != nil {
			h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Failed to upsert notification for request %s\t%s", gCtx.Request.URL.String(), err.Error()))
			h.InternalServerError(gCtx)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Unable to commit SQL transaction for request %s\t%s", gCtx.Request.URL.String(), err.Error()))
		h.InternalServerError(gCtx)
		return
	}

	if sendNotificationNow {
		h.logger.Info(fmt.Sprintf("NotificationHandler : SendNotification :: Scheduled for now? %t%d", sendNotificationNow, notification.ID))
		if err := h.RequestNotification(&notification, isIos); err != nil {
			h.logger.Error(err.Error())
			h.InternalServerError(gCtx)
			return
		}
	}

	gCtx.Status(http.StatusOK)
}

func (h *Handler) RequestNotification(notification *models.Notification, isIos bool) error {
	client := &http.Client{
		Timeout: types.TIMEOUT_TRANSACTION_SHORT,
	}
	reqPayload, err := json.Marshal(utils.NModelToNotificationReq(notification, isIos))
	if err != nil {
		return fmt.Errorf("NotificationHandler : RequestNotification :: Unable to marshall json for notification request %s", err.Error())
	}
	req, err := http.NewRequest(http.MethodPost, h.config.Notification.BaseURL, bytes.NewBuffer(reqPayload))
	if err != nil {
		return fmt.Errorf("NotificationHandler : RequestNotification :: http post request failed for notification request %s", err.Error())
	}

	// add headers
	req.Header.Add("Authorization", fmt.Sprintf("key=%s", h.config.Notification.AuthKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("NotificationHandler : RequestNotification :: http post request failed for notification request %s", err.Error())
	}
	defer res.Body.Close()
	client.CloseIdleConnections()
	return nil
}
