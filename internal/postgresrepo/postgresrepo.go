package postgresrepo

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"

	pb "github.com/awakair/awakair_todo_bot/api/todo-service"
)

type DbDriver interface {
	// Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	// QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type PostgresRepo struct {
	dbDriver DbDriver
}

func New(dbDriver DbDriver) *PostgresRepo {
	return &PostgresRepo{dbDriver: dbDriver}
}

func (pr PostgresRepo) SetUser(ctx context.Context, user *pb.User) error {
	const (
		queryEmptyUser = `INSERT INTO users (id)
	VALUES ($1::bigint);
	ON CONFLICT (id)
	DO NOTHING;`

		queryFullUser = `INSERT INTO users (id, language_code, utc_offset)
	VALUES ($1::bigint, $2, $3);
	ON CONFLICT (id)
	DO UPDATE SET language_code = $2, utc_offset = $3;`

		queryLanguageCode = `INSERT INTO users (id, language_code)
	VALUES ($1::bigint, $2);
	ON CONFLICT (id)
	DO UPDATE SET language_code = $2;`

		queryUtcOffset = `INSERT INTO users (id, utc_offset)
	VALUES ($1::bigint, $2);
	ON CONFLICT (id)
	DO UPDATE SET utc_offset = $2;`
	)

	if user.GetLanguageCode() == nil && user.GetUtcOffset() == nil {
		_, err := pr.dbDriver.Exec(ctx, queryEmptyUser, user.GetId())

		return err
	}

	if user.GetLanguageCode() != nil && user.GetUtcOffset() != nil {
		_, err := pr.dbDriver.Exec(
			ctx, queryFullUser, user.GetId(),
			user.GetLanguageCode().GetValue(), user.GetUtcOffset().GetValue(),
		)

		return err
	}

	if user.GetLanguageCode() != nil {
		_, err := pr.dbDriver.Exec(
			ctx, queryLanguageCode, user.GetId(), user.GetLanguageCode().GetValue(),
		)

		return err
	}

	_, err := pr.dbDriver.Exec(
		ctx, queryUtcOffset,
		user.GetId(),
		user.GetUtcOffset().GetValue(),
	)

	return err
}
