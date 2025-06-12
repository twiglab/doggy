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

const sum_sql = `
select
  sum(human_in)  as in_total ,
  sum(human_out) as out_total
from
  s_count
where
  _ts >= %d and _ts < %d  and
  uuid in %s
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

func OpenTaos(drv, dsn string) (*sql.DB, error) {
	db, err := sql.Open(drv, dsn)
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

func (o *OutS) Sum(ctx context.Context, in *oc.SumArg, out *oc.SumReply) error {
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

	out.InTotal = inTotal
	out.OutTotal = outTotal
	return nil
}

func (o *OutS) MutiSum(ctx context.Context, in *oc.MutiSumArg, out *oc.MutiSumReply) error {
	for _, arg := range in.Args {
		var reply oc.SumReply
		if err := o.Sum(ctx, &arg, &reply); err != nil {
			return err
		}
		out.Replies = append(out.Replies, reply)
	}
	return nil
}
