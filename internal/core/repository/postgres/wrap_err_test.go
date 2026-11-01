package postgres

import (
	"database/sql"
	stderrors "errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kr0uch/studspace/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestWrapErr(t *testing.T) {
	tests := []struct {
		name     string
		in       error
		expected error
	}{
		{name: "nil", in: nil, expected: nil},
		{name: "not found", in: sql.ErrNoRows, expected: errors.ErrNotFound},
		{name: "http passthrough", in: errors.ErrForbidden, expected: errors.ErrForbidden},
		{name: "exclusion", in: &pgconn.PgError{Code: pgExclusionViolation}, expected: errors.ErrSlotConflict},
		{name: "unique", in: &pgconn.PgError{Code: pgUniqueViolation}, expected: errors.ErrInternal},
		{name: "other", in: stderrors.New("boom"), expected: errors.ErrInternal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapErr(tt.in)
			if tt.expected == nil {
				assert.NoError(t, got)
				return
			}
			assert.ErrorIs(t, got, tt.expected)
		})
	}
}
