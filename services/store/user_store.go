package store

import (
	"context"
	"database/sql"
	"ncronus/database/mysql/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type UserStore struct{}

func NewUserStore() *UserStore {
	return &UserStore{}
}

func (u *UserStore) Insert(ctx context.Context, exec boil.ContextExecutor, userModel *models.User) error {
	if err := userModel.Insert(ctx, exec, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (u *UserStore) Get(ctx context.Context, exec boil.ContextExecutor, userID int) (*models.User, error) {
	user, err := models.FindUser(ctx, exec, userID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return user, nil
}

func (u *UserStore) GetByUsername(ctx context.Context, exec boil.ContextExecutor, username string) (*models.User, error) {
	user, err := models.Users(qm.Where("username = ?", username)).One(ctx, exec)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserStore) Exists(ctx context.Context, exec boil.ContextExecutor, userID int) (bool, error) {
	userExists, err := models.UserExists(ctx, exec, userID)
	if err != nil {
		return false, err
	}
	return userExists, nil
}

func (u *UserStore) Update(ctx context.Context, exec boil.ContextExecutor, userModel *models.User) error {
	_, err := userModel.Update(ctx, exec, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}
