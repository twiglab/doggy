package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type AutoReg struct {
	ent.Schema
}

func (AutoReg) Fields() []ent.Field {
	return []ent.Field{

		field.String("sn").
			MaxLen(36).
			NotEmpty().
			Unique().
			Immutable().
			SchemaType(map[string]string{
				dialect.MySQL:    "char(36)", // Override MySQL.
				dialect.Postgres: "char(36)", // Override Postgres.
				dialect.SQLite:   "char(36)", // Override Postgres.
			}),

		field.String("ip").
			MaxLen(64).NotEmpty().
			SchemaType(map[string]string{
				dialect.MySQL:    "varchar(32)", // Override MySQL.
				dialect.Postgres: "varchar(32)", // Override Postgres.
				dialect.SQLite:   "varchar(32)", // Override Postgres.
			}),

		/*
			field.String("nickname").
				MaxLen(64).NotEmpty().
				Unique().Immutable().
				SchemaType(map[string]string{
					dialect.MySQL:    "varchar(64)", // Override MySQL.
					dialect.Postgres: "varchar(64)", // Override Postgres.
					dialect.SQLite:   "varchar(64)", // Override Postgres.
				}),
		*/
	}
}

func (AutoReg) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (AutoReg) Indexes() []ent.Index {
	return []ent.Index{}
}

func (AutoReg) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "auto_reg"},
	}
}
