package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kr0uch/studspace/internal/core/models/api"
	"github.com/kr0uch/studspace/internal/core/models/dto"
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/kr0uch/studspace/pkg/web"
)

type UserService interface {
	UpsertUserByVK(ctx context.Context, in dto.UpsertUserByVK) (*entities.User, error)
	UpdateUser(ctx context.Context, in dto.UpdateUser) (*entities.User, error)
	UpdateUserNotifications(ctx context.Context, in dto.UpdateUserNotifications) (*entities.User, error)
}

type UserController struct {
	userService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{userService: userService}
}

// SyncHandler godoc
// @Summary Синхронизация пользователя из VK
// @Description Создаёт или обновляет профиль по vk_user_id из подписи. Имя берётся из тела, не из users.get. VK ID из admins.txt получает роль admin, остальные — student.
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security VKAuth
// @Param request body api.SyncUserRequest true "Имя, фамилия, фото"
// @Success 200 {object} entities.User
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /auth/sync [post]
func (c *UserController) SyncHandler(w http.ResponseWriter, r *http.Request) error {
	vkUserID, err := web.GetVKUserIDFromCtx(r.Context())
	if err != nil {
		return err
	}

	var in api.SyncUserRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return errors.ErrBadRequest
	}

	user, err := c.userService.UpsertUserByVK(r.Context(), dto.UpsertUserByVK{
		VKUserID:  vkUserID,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		PhotoURL:  in.PhotoURL,
	})
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(user)
}

// GetMeHandler godoc
// @Summary Текущий пользователь
// @Description Профиль из БД по vk_user_id из подписи.
// @Tags Пользователь
// @Produce json
// @Security VKAuth
// @Success 200 {object} entities.User
// @Failure 401 {object} errors.ErrorResponse
// @Failure 403 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /me [get]
func (c *UserController) GetMeHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(user)
}

// PatchMeHandler godoc
// @Summary Редактирование профиля
// @Description Отчество и должность/статус. Имя и фамилия только из VK (POST /auth/sync).
// @Tags Пользователь
// @Accept json
// @Produce json
// @Security VKAuth
// @Param request body api.PatchMeRequest true "профиль"
// @Success 200 {object} entities.User
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /me [patch]
func (c *UserController) PatchMeHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	var in api.PatchMeRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return errors.ErrBadRequest
	}

	updated, err := c.userService.UpdateUser(r.Context(), dto.UpdateUser{
		UserID:     user.UserID,
		Patronymic: in.Patronymic,
		Occupation: in.Occupation,
	})
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(updated)
}

// PatchNotificationsHandler godoc
// @Summary Флаг уведомлений
// @Description Только notifications_allowed в БД, без отправки в VK.
// @Tags Пользователь
// @Accept json
// @Security VKAuth
// @Param request body api.PatchNotificationsRequest true "allowed"
// @Success 204
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /me/notifications [patch]
func (c *UserController) PatchNotificationsHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}

	var in api.PatchNotificationsRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return errors.ErrBadRequest
	}

	if _, err := c.userService.UpdateUserNotifications(r.Context(), dto.UpdateUserNotifications{
		UserID:  user.UserID,
		Allowed: in.Allowed,
	}); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
