package services

import (
	"context"
	"errors"

	"github.com/guna/url-shortener/internal/models"
	"github.com/guna/url-shortener/internal/utils"
	"github.com/sirupsen/logrus"
)

type IUserService interface {
	AuthenticateUser(ctx context.Context, email, password string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID uint64) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.UpdateUserRequest) error
	LogoutUser(ctx context.Context, userID uint64) error
}

func (svc *Service) AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := svc.Storage.GetUserByEmail(email)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if !utils.ComparePasswords(user.Password, password) {
		logrus.Error(models.ErrPasswordDoestMatch)
		return nil, errors.New(models.ErrPasswordDoestMatch)
	}

	return user, nil
}

func (svc *Service) CreateUser(ctx context.Context, user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		logrus.Error(err)
		return err
	}

	user.Password = hashedPassword
	return svc.Storage.CreateUser(user)
}

func (svc *Service) GetUserByID(ctx context.Context, userID uint64) (*models.User, error) {
	user, err := svc.Storage.GetUserByID(userID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if user.Status == models.USER_INACTIVE {
		logrus.Error(models.ErrUserInactive)
		return nil, errors.New(models.ErrUserInactive)
	}
	return user, nil
}

func (svc *Service) UpdateUser(ctx context.Context, updateReq *models.UpdateUserRequest) error {
	user, err := svc.Storage.GetUserByEmail(updateReq.Email)
	if err != nil {
		logrus.Error(err)
		return err
	}
	if !utils.ComparePasswords(user.Password, updateReq.OldPassword) {
		logrus.Error(models.ErrInvalidCredentials)
		return errors.New(models.ErrInvalidCredentials)
	}
	hashedPassword, err := utils.HashPassword(updateReq.NewPassword)
	if err != nil {
		logrus.Error(err)
		return err
	}
	user.Password = hashedPassword
	return svc.Storage.UpdateUser(user)
}

func (svc *Service) LogoutUser(ctx context.Context, userID uint64) error {
	return svc.Storage.LogoutUser(userID)
}
