package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type Setup struct {
	ent.Schema
}

func (Setup) Fields() []ent.Field {
	return []ent.Field{

		field.String("sn").
			MaxLen(36).
			NotEmpty().
			Unique().
			Immutable().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(36)", // Override MySQL.
				dialect.Postgres: "varchar(36)", // Override Postgres.
				dialect.SQLite:   "varchar(36)", // Override Postgres.
			}),

		field.String("ip").
			MaxLen(64).NotEmpty().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(64)", // Override MySQL.
				dialect.Postgres: "varchar(64)", // Override Postgres.
				dialect.SQLite:   "varchar(64)", // Override Postgres.
			}),

		field.Time("last_time").
			Default(time.Now).
			UpdateDefault(time.Now),

		field.String("user").
			MaxLen(64).Optional().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(64)", // Override MySQL.
				dialect.Postgres: "varchar(64)", // Override Postgres.
				dialect.SQLite:   "varchar(64)", // Override Postgres.
			}),

		field.String("pwd").
			MaxLen(64).Optional().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(64)", // Override MySQL.
				dialect.Postgres: "varchar(64)", // Override Postgres.
				dialect.SQLite:   "varchar(64)", // Override Postgres.
			}),
		field.String("uuid1").
			MaxLen(36).
			NotEmpty().
			Unique().
			SchemaType(map[string]string{
				dialect.MySQL:    "char(36)", // Override MySQL.
				dialect.Postgres: "char(36)", // Override Postgres.
				dialect.SQLite:   "char(36)", // Override Postgres.
			}),

		field.String("uuid2").
			MaxLen(36).
			NotEmpty().
			Unique().
			SchemaType(map[string]string{
				dialect.MySQL:    "char(36)", // Override MySQL.
				dialect.Postgres: "char(36)", // Override Postgres.
				dialect.SQLite:   "char(36)", // Override Postgres.
			}),
	}
}

func (Setup) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Setup) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sn"),
	}
}

func (Setup) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "camera_upload"},
	}
}
