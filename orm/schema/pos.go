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

type Pos struct {
	ent.Schema
}

func (Pos) Fields() []ent.Field {
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

		field.String("pos").
			MaxLen(64).
			Immutable().
			NotEmpty().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(64)", // Override MySQL.
				dialect.Postgres: "varchar(64)", // Override Postgres.
				dialect.SQLite:   "varchar(64)", // Override Postgres.
			}),

		field.String("floor").
			MaxLen(64).Optional().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(64)", // Override MySQL.
				dialect.Postgres: "varchar(64)", // Override Postgres.
				dialect.SQLite:   "varchar(64)", // Override Postgres.
			}),

		field.String("building").
			MaxLen(64).Optional().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(64)", // Override MySQL.
				dialect.Postgres: "varchar(64)", // Override Postgres.
				dialect.SQLite:   "varchar(64)", // Override Postgres.
			}),

		field.String("area").
			MaxLen(64).Optional().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(64)", // Override MySQL.
				dialect.Postgres: "varchar(64)", // Override Postgres.
				dialect.SQLite:   "varchar(64)", // Override Postgres.
			}),
	}
}

func (Pos) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Pos) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sn"),
	}
}

func (Pos) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "camera_setup"},
	}
}
