package taosdb

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"

	_ "github.com/taosdata/driver-go/v3/taosWS"
	"github.com/twiglab/doggy/pkg/oc"
	// _ "github.com/taosdata/driver-go/v3/taosSql"
)

// select sum(human_in) as total  from s_count where _ts >= 1 and _ts < 100000000000000 and device_id in ('AAAAAAAAAA', 'BBBBB');

const sum_sql = `
select
  sum(human_in)  as in_total ,
  sum(human_out) as out_total
from
  s_count
where
  _ts >= %d and _ts < %d  and
  device_id in %s
`

func sqlIn(tables []string) string {
	ss := slices.Clone(tables)
	for i, s := range tables {
		ss[i] = "'" + strings.TrimSpace(s) + "'"
	}
	return "  (" + strings.Join(ss, ",") + ")"
}

func sumSQL(start, end int64, tables []string) string {
	inStr := sqlIn(tables)
	return fmt.Sprintf(sum_sql, start, end, inStr)
}

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

func (o *OutS) SumOf(ctx context.Context, in *oc.AreaArg, out *oc.Reply) error {
	sql := sumSQL(in.Start, in.End, in.IDs)

	rows, err := o.DB.QueryContext(ctx, sql)
	if err != nil {
		return err
	}

	var inTotal int64
	var outTotal int64

	for rows.Next() {
		if err = rows.Scan(&inTotal, &outTotal); err != nil {
			return err
		}
	}

	out.ValueA = inTotal
	out.ValueB = outTotal

	return nil
}

func (o *OutS) CollectOf(ctx context.Context, in *oc.AreaArg, out *oc.Reply) error {
	out.ValueA = 1
	out.ValueB = 2
	return nil
}
