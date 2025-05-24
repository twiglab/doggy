package taosdb

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/taosdata/driver-go/v3/taosWS"
	"github.com/twiglab/doggy/pkg/oc"
	// _ "github.com/taosdata/driver-go/v3/taosSql"
)

func OpenDB(conf Config) (*sql.DB, error) {
	driverName, dsn := db(conf)

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxLifetime(1 * time.Minute)

	return db, err
}

type OutS struct {
	DB *sql.DB
}

func (o *OutS) CollectOf(ctx context.Context, in *oc.AreaArg, out *oc.Reply) error {
	out.ValueA = 1
	out.ValueB = 2
	return nil
}

func (o *OutS) SumOf(ctx context.Context, in *oc.AreaArg, out *oc.Reply) error {
	out.ValueA = 1
	out.ValueB = 3
	return nil
}
