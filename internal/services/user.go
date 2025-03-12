package services

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/guna/url-shortener/internal/models"
	"github.com/guna/url-shortener/internal/utils"
	"github.com/sirupsen/logrus"
)

type IUserService interface {
	AuthenticateUser(ctx *fiber.Ctx, email, password string) (*models.User, error)
	CreateUser(ctx *fiber.Ctx, user *models.User) error
	GetUserByID(ctx *fiber.Ctx, userID uint64) (*models.User, error)
	UpdateUser(ctx *fiber.Ctx, user *models.UpdateUserRequest) error
	LogoutUser(ctx *fiber.Ctx, userID uint64) error
}

func (svc *Service) AuthenticateUser(ctx *fiber.Ctx, email, password string) (*models.User, error) {
	user, err := svc.Storage.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(models.ErrInvalidCredentials)
		}
		return nil, err
	}

	if !utils.ComparePasswords(user.Password, password) {
		logrus.Error(models.ErrPasswordDoestMatch)
		return nil, errors.New(models.ErrPasswordDoestMatch)
	}

	return user, nil
}

func (svc *Service) CreateUser(ctx *fiber.Ctx, user *models.User) error {
	_, err := svc.Storage.GetUserByEmail(user.Email)
	if err == nil {
		return errors.New(models.ErrUserWithEmailExists)
	}
	_, err = svc.Storage.GetUserByName(user.Name)
	if err == nil {
		return errors.New(models.ErrUserNameExists)
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		logrus.Error(err)
		return err
	}

	user.Password = hashedPassword
	return svc.Storage.CreateUser(user)
}

func (svc *Service) GetUserByID(ctx *fiber.Ctx, userID uint64) (*models.User, error) {
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

func (svc *Service) UpdateUser(ctx *fiber.Ctx, updateReq *models.UpdateUserRequest) error {
	user, err := svc.Storage.GetUserByID(updateReq.UserID)
	if err != nil {
		return err
	}
	if !utils.ComparePasswords(user.Password, updateReq.OldPassword) {
		logrus.Error(models.ErrInvalidCredentials)
		return errors.New(models.ErrInvalidCredentials)
	}
	if updateReq.UpdateName {
		usr, err := svc.Storage.GetUserByName(updateReq.Name)
		if err == nil && usr.ID != updateReq.UserID {
			return errors.New(models.ErrUserNameExists)
		}
	}
	user.Name = updateReq.Name
	hashedPassword, err := utils.HashPassword(updateReq.NewPassword)
	if err != nil {
		logrus.Error(err)
		return err
	}
	user.Password = hashedPassword
	return svc.Storage.UpdateUser(user, updateReq.UpdateName)
}

func (svc *Service) LogoutUser(ctx *fiber.Ctx, userID uint64) error {
	return svc.Storage.LogoutUser(userID)
}
