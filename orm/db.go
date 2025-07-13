package orm

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	_ "github.com/twiglab/doggy/orm/ent/runtime"

	"github.com/twiglab/doggy/orm/ent"
)

func OpenClient(driveName string, dataSourceName string, opts ...ent.Option) (*ent.Client, error) {
	return ent.Open(driveName, dataSourceName, opts...)
}
