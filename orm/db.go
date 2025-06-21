package orm

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	_ "github.com/twiglab/doggy/orm/ent/runtime"

	"github.com/twiglab/doggy/orm/ent"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

func pgxClient(in *sql.DB, opts ...ent.Option) *ent.Client {
	drv := entsql.OpenDB(dialect.Postgres, in)
	return ent.NewClient(append(opts, ent.Driver(drv))...)
}

func pgxDB(ctx context.Context, url string, ops ...stdlib.OptionOpenDB) (*sql.DB, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	return stdlib.OpenDBFromPool(pool, ops...), nil
}

func OpenClient(driveName string, dataSourceName string, opts ...ent.Option) (*ent.Client, error) {
	if driveName == "pgx" {
		db, err := pgxDB(context.Background(), dataSourceName)
		if err != nil {
			return nil, err
		}
		return pgxClient(db, opts...), nil
	}

	return ent.Open(driveName, dataSourceName, opts...)
}
