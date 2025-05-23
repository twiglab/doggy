// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/twiglab/doggy/orm/ent/using"
)

// Using is the model entity for the Using schema.
type Using struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Sn holds the value of the "sn" field.
	Sn string `json:"sn,omitempty"`
	// UUID holds the value of the "uuid" field.
	UUID string `json:"uuid,omitempty"`
	// DeviceID holds the value of the "device_id" field.
	DeviceID string `json:"device_id,omitempty"`
	// Alg holds the value of the "alg" field.
	Alg string `json:"alg,omitempty"`
	// Name holds the value of the "name" field.
	Name         string `json:"name,omitempty"`
	selectValues sql.SelectValues
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Using) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case using.FieldID:
			values[i] = new(sql.NullInt64)
		case using.FieldSn, using.FieldUUID, using.FieldDeviceID, using.FieldAlg, using.FieldName:
			values[i] = new(sql.NullString)
		case using.FieldCreateTime, using.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Using fields.
func (u *Using) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case using.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			u.ID = int(value.Int64)
		case using.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				u.CreateTime = value.Time
			}
		case using.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				u.UpdateTime = value.Time
			}
		case using.FieldSn:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field sn", values[i])
			} else if value.Valid {
				u.Sn = value.String
			}
		case using.FieldUUID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field uuid", values[i])
			} else if value.Valid {
				u.UUID = value.String
			}
		case using.FieldDeviceID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field device_id", values[i])
			} else if value.Valid {
				u.DeviceID = value.String
			}
		case using.FieldAlg:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field alg", values[i])
			} else if value.Valid {
				u.Alg = value.String
			}
		case using.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				u.Name = value.String
			}
		default:
			u.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Using.
// This includes values selected through modifiers, order, etc.
func (u *Using) Value(name string) (ent.Value, error) {
	return u.selectValues.Get(name)
}

// Update returns a builder for updating this Using.
// Note that you need to call Using.Unwrap() before calling this method if this Using
// was returned from a transaction, and the transaction was committed or rolled back.
func (u *Using) Update() *UsingUpdateOne {
	return NewUsingClient(u.config).UpdateOne(u)
}

// Unwrap unwraps the Using entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (u *Using) Unwrap() *Using {
	_tx, ok := u.config.driver.(*txDriver)
	if !ok {
		panic("ent: Using is not a transactional entity")
	}
	u.config.driver = _tx.drv
	return u
}

// String implements the fmt.Stringer.
func (u *Using) String() string {
	var builder strings.Builder
	builder.WriteString("Using(")
	builder.WriteString(fmt.Sprintf("id=%v, ", u.ID))
	builder.WriteString("create_time=")
	builder.WriteString(u.CreateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("update_time=")
	builder.WriteString(u.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("sn=")
	builder.WriteString(u.Sn)
	builder.WriteString(", ")
	builder.WriteString("uuid=")
	builder.WriteString(u.UUID)
	builder.WriteString(", ")
	builder.WriteString("device_id=")
	builder.WriteString(u.DeviceID)
	builder.WriteString(", ")
	builder.WriteString("alg=")
	builder.WriteString(u.Alg)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(u.Name)
	builder.WriteByte(')')
	return builder.String()
}

// Usings is a parsable slice of Using.
type Usings []*Using
