package postgresrepo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	pb "github.com/awakair/awakair_todo_bot/api/todo-service"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{pool: pool}
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
		_, err := pr.pool.Exec(ctx, queryEmptyUser, user.GetId())

		return err
	}

	if user.GetLanguageCode() != nil && user.GetUtcOffset() != nil {
		_, err := pr.pool.Exec(
			ctx, queryFullUser, user.GetId(),
			user.GetLanguageCode().GetValue(), user.GetUtcOffset().GetValue(),
		)

		return err
	}

	if user.GetLanguageCode() != nil {
		_, err := pr.pool.Exec(
			ctx, queryLanguageCode, user.GetId(), user.GetLanguageCode().GetValue(),
		)

		return err
	}

	_, err := pr.pool.Exec(
		ctx, queryUtcOffset,
		user.GetId(),
		user.GetUtcOffset().GetValue(),
	)

	return err
}
