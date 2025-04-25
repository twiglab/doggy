package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type Using struct {
	ent.Schema
}

func (Using) Fields() []ent.Field {
	return []ent.Field{

		field.String("sn").
			MaxLen(36).
			NotEmpty().
			Immutable().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(36)", // Override MySQL.
				dialect.Postgres: "varchar(36)", // Override Postgres.
				dialect.SQLite:   "varchar(36)", // Override Postgres.
			}),

		field.String("uuid").
			MaxLen(36).
			NotEmpty().
			SchemaType(map[string]string{
				dialect.MySQL:    "char(36)", // Override MySQL.
				dialect.Postgres: "char(36)", // Override Postgres.
				dialect.SQLite:   "char(36)", // Override Postgres.
			}),

		field.String("device_id").
			MaxLen(64).
			Optional().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(64)", // Override MySQL.
				dialect.Postgres: "varchar(64)", // Override Postgres.
				dialect.SQLite:   "varchar(64)", // Override Postgres.
			}),

		field.String("alg").
			MaxLen(16).NotEmpty().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(16)", // Override MySQL.
				dialect.Postgres: "varchar(16)", // Override Postgres.
				dialect.SQLite:   "varchar(16)", // Override Postgres.
			}),

		field.String("name").
			MaxLen(64).Optional().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(64)", // Override MySQL.
				dialect.Postgres: "varchar(64)", // Override Postgres.
				dialect.SQLite:   "varchar(64)", // Override Postgres.
			}),

		field.String("memo").
			MaxLen(64).Optional().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(64)", // Override MySQL.
				dialect.Postgres: "varchar(64)", // Override Postgres.
				dialect.SQLite:   "varchar(64)", // Override Postgres.
			}),

		field.String("bk").
			MaxLen(64).NotEmpty().Unique().Immutable().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(64)", // Override MySQL.
				dialect.Postgres: "varchar(64)", // Override Postgres.
				dialect.SQLite:   "varchar(64)", // Override Postgres.
			}),
	}
}

func (Using) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Using) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("uuid"),
	}
}

func (Using) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "camera_using"},
	}
}
