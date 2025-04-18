package orm

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	_ "github.com/twiglab/doggy/orm/ent/runtime"

	"github.com/twiglab/doggy/orm/ent"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

func OpenClient(in *sql.DB) *ent.Client {
	drv := entsql.OpenDB(dialect.Postgres, in)
	return ent.NewClient(ent.Driver(drv))
}

func FromURL(ctx context.Context, url string, ops ...stdlib.OptionOpenDB) (*sql.DB, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	return stdlib.OpenDBFromPool(pool, ops...), nil
}
