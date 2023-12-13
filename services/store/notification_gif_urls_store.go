package store

import (
	"context"
	"ncronus/database/mysql/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type NotificationGifUrlsStore struct{}

func NewGifUrlsStore() *NotificationGifUrlsStore {
	return &NotificationGifUrlsStore{}
}

func (n *NotificationGifUrlsStore) Insert(ctx context.Context, exec boil.ContextExecutor, nGifUrl *models.NotificationGifURL) error {
	if err := nGifUrl.Insert(ctx, exec, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (n *NotificationGifUrlsStore) GetByNID(ctx context.Context, exec boil.ContextExecutor, nID int) (models.NotificationGifURLSlice, error) {
	if gifUrls, err := models.NotificationGifUrls(qm.Where("nd_id = ?", nID)).All(ctx, exec); err != nil {
		return nil, err
	} else {
		return gifUrls, nil
	}
}
