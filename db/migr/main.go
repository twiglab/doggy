package main

/*
import (
	"context"
	"log"

	"github.com/twiglab/crm/member/orm"
	"github.com/twiglab/crm/member/orm/ent/migrate"
)

const DATABASE_URL = "user=crm password=cRm9ijn)OKM host=pipi.dev port=5432 database=crm sslmode=disable"

func main() {
	db, err := orm.FromURL(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}

	client := orm.OpenClient(db)

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
*/
