package store

import (
	"context"
	"ncronus/database/mysql/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type NotificationDataStore struct{}

func NewDataStore() *NotificationDataStore {
	return &NotificationDataStore{}
}

func (n *NotificationDataStore) Insert(ctx context.Context, exec boil.ContextExecutor, nData *models.NotificationDatum) error {
	if err := nData.Insert(ctx, exec, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (n *NotificationDataStore) Update(ctx context.Context, exec boil.ContextExecutor, nData *models.NotificationDatum) error {
	if _, err := nData.Update(ctx, exec, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (n *NotificationDataStore) QueryOne(ctx context.Context, exec boil.ContextExecutor, nID int, queryMods ...qm.QueryMod) (*models.NotificationDatum, error) {
	nData, err := models.NotificationData(queryMods...).One(ctx, exec)
	if err != nil {
		return nil, err
	}
	return nData, nil
}
