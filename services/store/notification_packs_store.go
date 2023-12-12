package store

import (
	"context"
	"ncronus/database/mysql/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type NotificationPackStore struct{}

func NewPacksStore() *NotificationPackStore {
	return &NotificationPackStore{}
}

func (n *NotificationPackStore) Insert(ctx context.Context, exec boil.ContextExecutor, nPack *models.NotificationPack) error {
	if err := nPack.Insert(ctx, exec, boil.Infer()); err != nil {
		return err
	}
	return nil
}
