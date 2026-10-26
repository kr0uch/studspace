package postgres

import (
	"database/sql"
	stderrors "errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/kr0uch/studspace/pkg/errors"
)

const (
	pgExclusionViolation = "23P01"
	pgUniqueViolation    = "23505"
)

func wrapErr(err error) error {
	if err == nil {
		return nil
	}
	if stderrors.Is(err, sql.ErrNoRows) {
		return errors.ErrNotFound
	}

	var httpErr *errors.HttpError
	if stderrors.As(err, &httpErr) {
		return err
	}

	var pgErr *pgconn.PgError
	if stderrors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgExclusionViolation:
			return errors.ErrSlotConflict.Wrap(err)
		case pgUniqueViolation:
			return errors.ErrInternal.Wrap(err)
		}
	}

	return errors.ErrInternal.Wrap(err)
}
