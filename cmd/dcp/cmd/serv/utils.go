package serv

import (
	"database/sql"
	"log"

	"github.com/spf13/viper"
	"github.com/twiglab/doggy/be/taosdb"
	"github.com/twiglab/doggy/holo"
	"github.com/twiglab/doggy/kv"
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


func subviper(prefix string, v *viper.Viper) *viper.Viper {
	if len(prefix) == 0 {
		return v
	}
	return v.Sub("backend." + prefix)
}
