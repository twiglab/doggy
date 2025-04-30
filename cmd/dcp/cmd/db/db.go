package db

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/orm/ent/migrate"
)

// DbCmd represents the serv command
var DbCmd = &cobra.Command{
	Use:   "db",
	Short: "数据库操作",
	Long:  `数据库操作`,
	Run: func(cmd *cobra.Command, args []string) {
		dbMigrate()
	},
	Example: "dcp db --addr 1.2.3.4",
}

var (
	dsn    string
	dbname string
)

func init() {
	DbCmd.Flags().StringVar(&dbname, "driver", "sqlite3", "数据库名称")
	DbCmd.Flags().StringVar(&dsn, "datasource", "dcp.db?_fk=1", "数据库dsn")
}

func dbMigrate() {
	client, err := orm.OpenClient(dbname, dsn)
	if err!=nil{
		log.Fatal(err)
	}

	ctx := context.Background()
	err = client.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
