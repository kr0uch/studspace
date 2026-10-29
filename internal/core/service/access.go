package service

import (
	"github.com/kr0uch/studspace/internal/core/models/entities"
	"github.com/kr0uch/studspace/pkg/errors"
)

func requireAdmin(actor entities.User) error {
	if !actor.IsAdmin() {
		return errors.ErrForbidden
	}
	return nil
}
