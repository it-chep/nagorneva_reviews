// Package postgres contains database boundaries shared by repositories.
package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nagorneva/nagorneva_reviews/internal/pkg/domain_error"
)

// DB can be implemented by a pool or a transaction.
type DB interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

// Error translates persistence errors before they leave a repository.
func Error(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain_error.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return domain_error.ErrConflict
		case "23503":
			return domain_error.ErrReference
		case "23514", "23502", "22003":
			return domain_error.ErrInvalid
		}
	}
	return err
}
func DeleteResult(tag pgconn.CommandTag, err error) error {
	if err != nil {
		return Error(err)
	}
	if tag.RowsAffected() == 0 {
		return domain_error.ErrNotFound
	}
	return nil
}

// Collect closes rows before another query can acquire a connection.
func Collect[T any](rows pgx.Rows, scan func(pgx.Row) (*T, error)) ([]T, error) {
	defer rows.Close()
	result := make([]T, 0)
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, Error(err)
		}
		result = append(result, *item)
	}
	return result, Error(rows.Err())
}
