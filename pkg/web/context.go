package web

import (
	"context"

	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
)

type ctxKey int

const (
	ctxHttpErrorKey ctxKey = iota
	ctxNotifyErrorKey
	ctxUserKey
	ctxVKUserIDKey
)

func ContextWithUser(ctx context.Context, user *entities.User) context.Context {
	return context.WithValue(ctx, ctxUserKey, user)
}

func ContextWithVKUserID(ctx context.Context, vkUserID int64) context.Context {
	return context.WithValue(ctx, ctxVKUserIDKey, vkUserID)
}

func GetUserFromCtx(ctx context.Context) (*entities.User, error) {
	user, ok := ctx.Value(ctxUserKey).(*entities.User)
	if !ok || user == nil {
		return nil, errors.ErrForbidden
	}
	return user, nil
}

func GetVKUserIDFromCtx(ctx context.Context) (int64, error) {
	vkUserID, ok := ctx.Value(ctxVKUserIDKey).(int64)
	if !ok || vkUserID < 1 {
		return 0, errors.ErrInvalidSign
	}
	return vkUserID, nil
}

func GetNotifyErrorFromCtx(ctx context.Context) *errors.NotifyError {
	val := ctx.Value(ctxNotifyErrorKey)
	if val == nil {
		return nil
	}
	err, ok := val.(*errors.NotifyError)
	if !ok {
		return nil
	}
	return err
}

func GetHttpErrorFromCtx(ctx context.Context) *errors.HttpError {
	val := ctx.Value(ctxHttpErrorKey)
	if val == nil {
		return nil
	}

	err, ok := val.(*errors.HttpError)
	if !ok {
		return nil
	}
	return err
}
