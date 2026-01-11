package ddb

import (
	"context"
	"database/sql"
	"sync"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
	"github.com/twiglab/doggy/pf"
)

type DuckDB struct {
	db   *sql.DB
	from string
	tbl  string
	mu   sync.Mutex
}

func New(from string) (*DuckDB, error) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DuckDB{
		db:   db,
		from: from,
		tbl:  "x",
	}, nil
}

func (d *DuckDB) Load(ctx context.Context) error {

	nextTbl, cr := losdSql(d.tbl, d.from)
	if _, err := d.db.ExecContext(ctx, cr); err != nil {
		return err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.tbl = nextTbl
	return nil
}

func (d *DuckDB) Loop(ctx context.Context) (chan struct{}, chan struct{}, error) {
	if err := d.Load(ctx); err != nil {
		return nil, nil, err
	}

	reLoadCh := make(chan struct{})
	soptCh := make(chan struct{})

	go func(ctx context.Context) {

		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_ = d.Load(ctx)
			case <-reLoadCh:
				_ = d.Load(ctx)
			case <-soptCh:
				return
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	return reLoadCh, soptCh, nil
}

func (d *DuckDB) List(ctx context.Context) (rs []pf.ChannelUserData, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	listsql := listSql(d.tbl)

	var rows *sql.Rows
	if rows, err = d.db.QueryContext(ctx, listsql); err != nil {
		return
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var data pf.ChannelUserData
		if err = rows.Scan(&data.SN, &data.UUID, &data.Code, &data.X, &data.Y, &data.Z); err != nil {
			return
		}
		rs = append(rs, data)
	}
	return
}

func (d *DuckDB) TblName(ctx context.Context) string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.tbl
}

func (d *DuckDB) Get(ctx context.Context, id string) (data pf.ChannelUserData, ok bool, err error) {
	tblname := d.TblName(ctx)
	sql := querySql(tblname)
	row := d.db.QueryRowContext(ctx, sql, id)
	err = row.Scan(&data.SN, &data.UUID, &data.Code, &data.X, &data.Y, &data.Z)
	ok = err == nil
	return
}

func (d *DuckDB) Set(_ context.Context, _ string, _ pf.ChannelUserData) error { return nil }
