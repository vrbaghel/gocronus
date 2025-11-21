package handler

import (
	"context"
	"fmt"
	"log"
	"ncronus/database/mysql/models"
	"ncronus/services/types"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (h *Handler) RestartCronJobs() {
	txCtx, cancel := context.WithTimeout(context.Background(), types.TIMEOUT_TRANSACTION_SHORT)
	defer cancel()
	tx, err := h.mySql.Client.BeginTx(context.Background(), nil)
	if err != nil {
		h.logger.Fatal("CronHandler : RestartCronJobs :: Failed to begin transaction")
		return
	}
	qMods := []qm.QueryMod{qm.Load(qm.Rels(models.NotificationRels.IDNotificationDatum, models.NotificationDatumRels.IDNotificationPack)), qm.Where("cron_job_id > ?", 0), qm.And("n_status = ?", models.NotificationNStatusScheduled)}
	notificationSlice, err := h.store.NotificationStore.QueryAll(txCtx, tx, qMods...)
	if err != nil {
		h.logger.Fatal("CronHandler : RestartCronJobs :: Failed to fetch notifications with pending cron job")
		return
	}
	log.Printf("CronHandler : RestartCronJobs :: Restarting %d jobs...\n", len(notificationSlice))
	for _, notification := range notificationSlice {
		imgUrls, err := h.store.NotificationImgUrlsStore.GetByNID(txCtx, tx, notification.ID)
		if err != nil {
			h.logger.Fatal(fmt.Sprintf("CronHandler : RestartCronJobs :: Failed to fetch notification img urls for notification %d with pending cron job %s", notification.ID, err))
			return
		}
		gifUrls, err := h.store.NotificationGifUrlsStore.GetByNID(txCtx, tx, notification.ID)
		if err != nil {
			h.logger.Fatal(fmt.Sprintf("CronHandler : RestartCronJobs :: Failed to fetch notification gif urls for notification %d with pending cron job %s", notification.ID, err))
			return
		}
		h.ScheduleNotification(notification, notification.NTimezone == models.NotificationNTimezoneIST, imgUrls, gifUrls)
	}
	log.Println("CronHandler : RestartCronJobs :: All jobs started")
	if err := tx.Commit(); err != nil {
		h.logger.Fatal("CronHandler : RestartCronJobs :: Failed to commit SQL transaction")
		return
	}
}
