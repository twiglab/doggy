package sqlite

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"github.com/twiglab/doggy/orm"
	"github.com/twiglab/doggy/orm/ent/migrate"
)

// SqliteCmd represents the serv command
var SqliteCmd = &cobra.Command{
	Use:   "sqlite",
	Short: "sqlite",
	Long:  `sqlite`,
	Run: func(cmd *cobra.Command, args []string) {
		dbMigrate(cmd, args)
	},
	Example: "dcp sqlite --db dcp",
}

var (
	name string
)

func init() {
	SqliteCmd.Flags().StringVar(&name, "db", "dcp", "数据库名称")
}

func dbMigrate(_ *cobra.Command, _ []string) {

	dbname := name + ".db?_fk=1"

	client, err := orm.OpenClient("sqlite3", dbname)

	ctx := context.Background()
	// Run migration.
	err = client.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)
	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
