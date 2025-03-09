package services

import (
	"context"
	"errors"

	"github.com/guna/url-shortener/internal/models"
	"github.com/guna/url-shortener/internal/utils"
)

type IUserService interface {
	AuthenticateUser(ctx context.Context, email, password string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID uint64) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	LogoutUser(ctx context.Context, userID uint64) error
}

func (svc *Service) AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := svc.Storage.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !utils.ComparePasswords(user.Password, password) {
		return nil, errors.New(models.ErrInvalidCredentials)
	}

	return user, nil
}

func (svc *Service) CreateUser(ctx context.Context, user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return svc.Storage.CreateUser(user)
}

func (svc *Service) GetUserByID(ctx context.Context, userID uint64) (*models.User, error) {
	return svc.Storage.GetUserByID(userID)
}

func (svc *Service) UpdateUser(ctx context.Context, user *models.User) error {
	return svc.Storage.UpdateUser(user)
}

func (svc *Service) LogoutUser(ctx context.Context, userID uint64) error {
	return svc.Storage.LogoutUser(userID)
}
