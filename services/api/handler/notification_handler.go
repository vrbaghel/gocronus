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
	"github.com/robfig/cron/v3"
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
	sendNotificationNow := false
	if err := gCtx.ShouldBindJSON(&reqPayload); err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Unable to bind request body with GenerateImgResponse %s", err.Error()))
		gCtx.JSON(http.StatusBadRequest, types.APIError{
			Code:    http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Message: "invalid payload",
		})
		return
	}
	isIos := reqPayload.Device == models.NotificationNDeviceIos

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
	}
	if reqPayload.Timezone != "" {
		notification.NTimezone = reqPayload.Timezone
	} else {
		notification.NTimezone = models.NotificationNTimezoneIST
	}

	var isISTZone bool = notification.NTimezone == models.NotificationNTimezoneIST

	if reqPayload.ScheduledFor != "" {
		var scheduledFor time.Time
		var err error
		if isISTZone {
			scheduledFor, err = time.ParseInLocation(types.NOTIFICATION_TIMESTAMP_FORMAT, reqPayload.ScheduledFor, types.IST_TIMEZONE)
		} else {
			scheduledFor, err = time.ParseInLocation(types.NOTIFICATION_TIMESTAMP_FORMAT, reqPayload.ScheduledFor, types.CST_TIMEZONE)
		}
		if err != nil {
			h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Failed to parse timestamp for scheduling notification %s\t%s", reqPayload.ScheduledFor, err.Error()))
			gCtx.JSON(http.StatusBadRequest, types.APIError{
				Code:    http.StatusBadRequest,
				Status:  http.StatusText(http.StatusBadRequest),
				Message: "invalid value for scheduled_for",
			})
			return
		}
		notification.NStatus = models.NotificationNStatusScheduled
		notification.NTimestamp = scheduledFor
	} else {
		notification.NStatus = models.NotificationNStatusCompleted
		notification.NTimestamp = time.Now().In(types.IST_TIMEZONE)
		sendNotificationNow = true
		h.logger.Info(fmt.Sprintf("NotificationHandler : SendNotification :: scheduledNow %+v", notification.NTimestamp))
	}

	if err := h.store.NotificationStore.Insert(txCtx, tx, &notification); err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : SendNotification :: Failed to insert notification for request %s\t%s", gCtx.Request.URL.String(), err.Error()))
		h.InternalServerError(gCtx)
		return
	}

	if notification.ID > 0 && (reqPayload.Category != nil || reqPayload.Navigation != nil) {
		nData := models.NotificationDatum{
			NDSource:      0,
			NDClickAction: null.StringFrom(reqPayload.ClickAction),
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

		if reqPayload.Category != nil && reqPayload.Category.Data != nil {
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

		if reqPayload.Navigation != nil && reqPayload.Navigation.Data != nil {
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
		// h.logger.Info(fmt.Sprintf("NotificationHandler : SendNotification :: Scheduled for now? %t%d", sendNotificationNow, notification.ID))
		if err := h.RequestNotification(utils.NModelToNotificationReq(&notification, isIos)); err != nil {
			h.logger.Error(err.Error())
			h.InternalServerError(gCtx)
			return
		}
	} else {
		go h.ScheduleNotification(&notification, isISTZone)
	}

	gCtx.Status(http.StatusOK)
}

func (h *Handler) RequestNotification(nPayload types.RequestNotificationPayload) error {
	client := &http.Client{
		Timeout: types.TIMEOUT_TRANSACTION_SHORT,
	}
	reqPayload, err := json.Marshal(nPayload)
	if err != nil {
		return fmt.Errorf("NotificationHandler : RequestNotification :: Unable to marshall json for notification request %s", err.Error())
	}
	// h.logger.Info(fmt.Sprintf("NotificationHandler : RequestNotification :: Req payload %+v", string(reqPayload)))
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

func (h *Handler) ScheduleNotification(notification *models.Notification, isISTZone bool) {
	nTimeStamp := notification.NTimestamp
	nID := notification.ID
	day := nTimeStamp.Day()
	month := nTimeStamp.Month()
	hour := nTimeStamp.Hour()
	minute := nTimeStamp.Minute()
	// h.logger.Info(fmt.Sprintf("NotificationHandler : ScheduleNotification :: Scheduled for %d-%d %d:%d:%d", day, month, hour, minute, second))

	txCtx, cancel := context.WithTimeout(context.Background(), types.TIMEOUT_TRANSACTION_SHORT)
	defer cancel()
	tx, err := h.mySql.Client.BeginTx(context.Background(), nil)
	if err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : ScheduleNotification :: Unable to begin sql transaction for notification cron job %d\t%s", nID, err.Error()))
	}

	var scheduler *cron.Cron
	if isISTZone {
		scheduler = h.cron.IST
	} else {
		scheduler = h.cron.CST
	}
	cronExpression := fmt.Sprintf("%d %d %d %d *", minute, hour, day, month)
	jobId, err := scheduler.AddFunc(cronExpression, func() {
		h.ScheduleNotificationHandler(notification)
	})
	if err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : ScheduleNotification :: Failed to schedule notification for cron job scheduled at %s for notification %d\t%s", cronExpression, nID, err.Error()))
	}

	notification.CronJobID = null.IntFrom(int(jobId))
	notification.NStatus = models.NotificationNStatusScheduled

	if err := h.store.NotificationStore.Update(txCtx, tx, notification); err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : ScheduleNotification :: Failed to update notification %d with cron job ID %d\t%s", nID, jobId, err.Error()))
	}

	if err := tx.Commit(); err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : ScheduleNotification :: Unable to commit SQL transaction for notification cron job %d\t%s", nID, err.Error()))
	}
}

func (h *Handler) ScheduleNotificationHandler(notification *models.Notification) {
	nID := notification.ID
	if err := h.RequestNotification(utils.NModelToNotificationReq(notification, notification.NDevice == models.NotificationNDeviceIos)); err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : ScheduleNotificationHandler :: Failed to request scheduled notification %d with cron job ID %d\t%s", nID, notification.CronJobID.Int, err.Error()))
		// remove cron jobs if failed
		h.cron.CST.Remove(cron.EntryID(notification.CronJobID.Int))
		h.cron.IST.Remove(cron.EntryID(notification.CronJobID.Int))
		return
	}
	// remove cron jobs once completed
	h.cron.CST.Remove(cron.EntryID(notification.CronJobID.Int))
	h.cron.IST.Remove(cron.EntryID(notification.CronJobID.Int))

	txCtx, cancel := context.WithTimeout(context.Background(), types.TIMEOUT_TRANSACTION_SHORT)
	defer cancel()
	tx, err := h.mySql.Client.BeginTx(context.Background(), nil)
	if err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : ScheduleNotificationHandler :: Unable to begin sql transaction for notification cron job %d\t%s", nID, err.Error()))
		return
	}
	notification.NStatus = models.NotificationNStatusCompleted
	if err := h.store.NotificationStore.Update(txCtx, tx, notification); err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : ScheduleNotificationHandler :: Failed to update notification %d with cron job ID %d\t%s", nID, notification.CronJobID.Int, err.Error()))
	}
	if err := tx.Commit(); err != nil {
		h.logger.Error(fmt.Sprintf("NotificationHandler : ScheduleNotificationHandler :: Unable to commit SQL transaction for notification cron job %d\t%s", nID, err.Error()))
	}
}
