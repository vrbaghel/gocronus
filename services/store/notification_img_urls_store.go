package store

import (
	"context"
	"ncronus/database/mysql/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type NotificationImgUrlsStore struct{}

func NewImgUrlsStore() *NotificationImgUrlsStore {
	return &NotificationImgUrlsStore{}
}

func (n *NotificationImgUrlsStore) Insert(ctx context.Context, exec boil.ContextExecutor, nImgUrl *models.NotificationImgURL) error {
	if err := nImgUrl.Insert(ctx, exec, boil.Infer()); err != nil {
		return err
	}
	return nil
}
