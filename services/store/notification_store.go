package store

import (
	"context"
	"ncronus/database/mysql/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type NotificationStore struct{}

func NewNotificationStore() *NotificationStore {
	return &NotificationStore{}
}

func (n *NotificationStore) Insert(ctx context.Context, exec boil.ContextExecutor, notification *models.Notification) error {
	if err := notification.Insert(ctx, exec, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (n *NotificationStore) Count(ctx context.Context, exec boil.ContextExecutor, queryMods ...qm.QueryMod) (int64, error) {
	count, err := models.Notifications(queryMods...).Count(ctx, exec)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (n *NotificationStore) GetByID(ctx context.Context, exec boil.ContextExecutor, nID int) (*models.Notification, error) {
	notification, err := models.FindNotification(ctx, exec, nID)
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (n *NotificationStore) QueryOne(ctx context.Context, exec boil.ContextExecutor, nID int, queryMods ...qm.QueryMod) (*models.Notification, error) {
	notification, err := models.Notifications(queryMods...).One(ctx, exec)
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (n *NotificationStore) QueryAll(ctx context.Context, exec boil.ContextExecutor, queryMods ...qm.QueryMod) (models.NotificationSlice, error) {
	notification, err := models.Notifications(queryMods...).All(ctx, exec)
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (n *NotificationStore) GetAll(ctx context.Context, exec boil.ContextExecutor, queryMods ...qm.QueryMod) (models.NotificationSlice, error) {
	notifications, err := models.Notifications(queryMods...).All(ctx, exec)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (n *NotificationStore) Update(ctx context.Context, exec boil.ContextExecutor, notification *models.Notification) error {
	if _, err := notification.Update(ctx, exec, boil.Infer()); err != nil {
		return err
	}
	return nil
}
