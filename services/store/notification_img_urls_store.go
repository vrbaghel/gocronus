package store

import (
	"context"
	"ncronus/database/mysql/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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

func (n *NotificationImgUrlsStore) GetByNID(ctx context.Context, exec boil.ContextExecutor, nID int) (models.NotificationImgURLSlice, error) {
	if imgUrls, err := models.NotificationImgUrls(qm.Where("nd_id = ?", nID)).All(ctx, exec); err != nil {
		return nil, err
	} else {
		return imgUrls, nil
	}
}
