package postgres

import (
	"context"
	"time"

	"github.com/kr0uch/studspace/pkg/vk"
)

const notifyTimeout = 8 * time.Second

type NotifyRepository struct {
	client *vk.VKClient
}

func NewNotifyRepository(client *vk.VKClient) *NotifyRepository {
	return &NotifyRepository{client: client}
}

func (r *NotifyRepository) Notify(ctx context.Context, vkUserIDs []int64, message string) ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), notifyTimeout)
	defer cancel()
	err := r.client.Send(ctx, vkUserIDs, message)
	return vk.DisabledVKUserIDs(err), err
}
