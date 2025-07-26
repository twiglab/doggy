package serv

import (
	"database/sql"
	"log"

	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/kv"
	"github.com/twiglab/doggy/taosdb"
)

func MustOpenTaosDB(conf AppConf) *sql.DB {
	drv, dsn := taosdb.TaosDSN(
		conf.BackendConf.TaosDBConf.Username,
		conf.BackendConf.TaosDBConf.Password,
		conf.BackendConf.TaosDBConf.Addr,
		conf.BackendConf.TaosDBConf.Port,
		conf.BackendConf.TaosDBConf.DBName,
	)
	db, err := sql.Open(drv, dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func MustSubReq(req holo.SubscriptionReq, err error) holo.SubscriptionReq {
	if err != nil {
		log.Fatal(err)
	}
	return req
}

func MustOpenKV(conf AppConf) *kv.Handle {
	h, err := kv.FromURLs(conf.EtcdConf.URLs)
	if err != nil {
		log.Fatal(err)
	}
	return h
}
