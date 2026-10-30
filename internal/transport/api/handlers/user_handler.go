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
	UpdateUserNotifications(ctx context.Context, in dto.UpdateUserNotifications) (*entities.User, error)
}

type UserController struct {
	userService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{userService: userService}
}

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

func (c *UserController) GetMeHandler(w http.ResponseWriter, r *http.Request) error {
	user, err := web.GetUserFromCtx(r.Context())
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(user)
}

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
