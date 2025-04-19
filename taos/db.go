package taos

import (
	"database/sql"

	_ "github.com/taosdata/driver-go/v3/taosSql"
	_ "github.com/taosdata/driver-go/v3/taosWS"
)

type TaoSConfig struct {
	DriverName string
	DataSource string
}

func Open(config TaoSConfig) (*sql.DB, error) {
	db, err := sql.Open(config.DriverName, config.DataSource)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)
	db.SetConnMaxIdleTime(30)
	db.SetConnMaxLifetime(60)

	return db, nil
}
