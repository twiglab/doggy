// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/twiglab/doggy/orm/ent/using"
)

// UsingCreate is the builder for creating a Using entity.
type UsingCreate struct {
	config
	mutation *UsingMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreateTime sets the "create_time" field.
func (uc *UsingCreate) SetCreateTime(t time.Time) *UsingCreate {
	uc.mutation.SetCreateTime(t)
	return uc
}

// SetNillableCreateTime sets the "create_time" field if the given value is not nil.
func (uc *UsingCreate) SetNillableCreateTime(t *time.Time) *UsingCreate {
	if t != nil {
		uc.SetCreateTime(*t)
	}
	return uc
}

// SetUpdateTime sets the "update_time" field.
func (uc *UsingCreate) SetUpdateTime(t time.Time) *UsingCreate {
	uc.mutation.SetUpdateTime(t)
	return uc
}

// SetNillableUpdateTime sets the "update_time" field if the given value is not nil.
func (uc *UsingCreate) SetNillableUpdateTime(t *time.Time) *UsingCreate {
	if t != nil {
		uc.SetUpdateTime(*t)
	}
	return uc
}

// SetSn sets the "sn" field.
func (uc *UsingCreate) SetSn(s string) *UsingCreate {
	uc.mutation.SetSn(s)
	return uc
}

// SetUUID sets the "uuid" field.
func (uc *UsingCreate) SetUUID(s string) *UsingCreate {
	uc.mutation.SetUUID(s)
	return uc
}

// SetDeviceID sets the "device_id" field.
func (uc *UsingCreate) SetDeviceID(s string) *UsingCreate {
	uc.mutation.SetDeviceID(s)
	return uc
}

// SetAlg sets the "alg" field.
func (uc *UsingCreate) SetAlg(s string) *UsingCreate {
	uc.mutation.SetAlg(s)
	return uc
}

// SetName sets the "name" field.
func (uc *UsingCreate) SetName(s string) *UsingCreate {
	uc.mutation.SetName(s)
	return uc
}

// SetNillableName sets the "name" field if the given value is not nil.
func (uc *UsingCreate) SetNillableName(s *string) *UsingCreate {
	if s != nil {
		uc.SetName(*s)
	}
	return uc
}

// Mutation returns the UsingMutation object of the builder.
func (uc *UsingCreate) Mutation() *UsingMutation {
	return uc.mutation
}

// Save creates the Using in the database.
func (uc *UsingCreate) Save(ctx context.Context) (*Using, error) {
	uc.defaults()
	return withHooks(ctx, uc.sqlSave, uc.mutation, uc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (uc *UsingCreate) SaveX(ctx context.Context) *Using {
	v, err := uc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (uc *UsingCreate) Exec(ctx context.Context) error {
	_, err := uc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uc *UsingCreate) ExecX(ctx context.Context) {
	if err := uc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uc *UsingCreate) defaults() {
	if _, ok := uc.mutation.CreateTime(); !ok {
		v := using.DefaultCreateTime()
		uc.mutation.SetCreateTime(v)
	}
	if _, ok := uc.mutation.UpdateTime(); !ok {
		v := using.DefaultUpdateTime()
		uc.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uc *UsingCreate) check() error {
	if _, ok := uc.mutation.CreateTime(); !ok {
		return &ValidationError{Name: "create_time", err: errors.New(`ent: missing required field "Using.create_time"`)}
	}
	if _, ok := uc.mutation.UpdateTime(); !ok {
		return &ValidationError{Name: "update_time", err: errors.New(`ent: missing required field "Using.update_time"`)}
	}
	if _, ok := uc.mutation.Sn(); !ok {
		return &ValidationError{Name: "sn", err: errors.New(`ent: missing required field "Using.sn"`)}
	}
	if v, ok := uc.mutation.Sn(); ok {
		if err := using.SnValidator(v); err != nil {
			return &ValidationError{Name: "sn", err: fmt.Errorf(`ent: validator failed for field "Using.sn": %w`, err)}
		}
	}
	if _, ok := uc.mutation.UUID(); !ok {
		return &ValidationError{Name: "uuid", err: errors.New(`ent: missing required field "Using.uuid"`)}
	}
	if v, ok := uc.mutation.UUID(); ok {
		if err := using.UUIDValidator(v); err != nil {
			return &ValidationError{Name: "uuid", err: fmt.Errorf(`ent: validator failed for field "Using.uuid": %w`, err)}
		}
	}
	if _, ok := uc.mutation.DeviceID(); !ok {
		return &ValidationError{Name: "device_id", err: errors.New(`ent: missing required field "Using.device_id"`)}
	}
	if v, ok := uc.mutation.DeviceID(); ok {
		if err := using.DeviceIDValidator(v); err != nil {
			return &ValidationError{Name: "device_id", err: fmt.Errorf(`ent: validator failed for field "Using.device_id": %w`, err)}
		}
	}
	if _, ok := uc.mutation.Alg(); !ok {
		return &ValidationError{Name: "alg", err: errors.New(`ent: missing required field "Using.alg"`)}
	}
	if v, ok := uc.mutation.Alg(); ok {
		if err := using.AlgValidator(v); err != nil {
			return &ValidationError{Name: "alg", err: fmt.Errorf(`ent: validator failed for field "Using.alg": %w`, err)}
		}
	}
	if v, ok := uc.mutation.Name(); ok {
		if err := using.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Using.name": %w`, err)}
		}
	}
	return nil
}

func (uc *UsingCreate) sqlSave(ctx context.Context) (*Using, error) {
	if err := uc.check(); err != nil {
		return nil, err
	}
	_node, _spec := uc.createSpec()
	if err := sqlgraph.CreateNode(ctx, uc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	uc.mutation.id = &_node.ID
	uc.mutation.done = true
	return _node, nil
}

func (uc *UsingCreate) createSpec() (*Using, *sqlgraph.CreateSpec) {
	var (
		_node = &Using{config: uc.config}
		_spec = sqlgraph.NewCreateSpec(using.Table, sqlgraph.NewFieldSpec(using.FieldID, field.TypeInt))
	)
	_spec.OnConflict = uc.conflict
	if value, ok := uc.mutation.CreateTime(); ok {
		_spec.SetField(using.FieldCreateTime, field.TypeTime, value)
		_node.CreateTime = value
	}
	if value, ok := uc.mutation.UpdateTime(); ok {
		_spec.SetField(using.FieldUpdateTime, field.TypeTime, value)
		_node.UpdateTime = value
	}
	if value, ok := uc.mutation.Sn(); ok {
		_spec.SetField(using.FieldSn, field.TypeString, value)
		_node.Sn = value
	}
	if value, ok := uc.mutation.UUID(); ok {
		_spec.SetField(using.FieldUUID, field.TypeString, value)
		_node.UUID = value
	}
	if value, ok := uc.mutation.DeviceID(); ok {
		_spec.SetField(using.FieldDeviceID, field.TypeString, value)
		_node.DeviceID = value
	}
	if value, ok := uc.mutation.Alg(); ok {
		_spec.SetField(using.FieldAlg, field.TypeString, value)
		_node.Alg = value
	}
	if value, ok := uc.mutation.Name(); ok {
		_spec.SetField(using.FieldName, field.TypeString, value)
		_node.Name = value
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Using.Create().
//		SetCreateTime(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.UsingUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (uc *UsingCreate) OnConflict(opts ...sql.ConflictOption) *UsingUpsertOne {
	uc.conflict = opts
	return &UsingUpsertOne{
		create: uc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Using.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (uc *UsingCreate) OnConflictColumns(columns ...string) *UsingUpsertOne {
	uc.conflict = append(uc.conflict, sql.ConflictColumns(columns...))
	return &UsingUpsertOne{
		create: uc,
	}
}

type (
	// UsingUpsertOne is the builder for "upsert"-ing
	//  one Using node.
	UsingUpsertOne struct {
		create *UsingCreate
	}

	// UsingUpsert is the "OnConflict" setter.
	UsingUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdateTime sets the "update_time" field.
func (u *UsingUpsert) SetUpdateTime(v time.Time) *UsingUpsert {
	u.Set(using.FieldUpdateTime, v)
	return u
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *UsingUpsert) UpdateUpdateTime() *UsingUpsert {
	u.SetExcluded(using.FieldUpdateTime)
	return u
}

// SetUUID sets the "uuid" field.
func (u *UsingUpsert) SetUUID(v string) *UsingUpsert {
	u.Set(using.FieldUUID, v)
	return u
}

// UpdateUUID sets the "uuid" field to the value that was provided on create.
func (u *UsingUpsert) UpdateUUID() *UsingUpsert {
	u.SetExcluded(using.FieldUUID)
	return u
}

// SetDeviceID sets the "device_id" field.
func (u *UsingUpsert) SetDeviceID(v string) *UsingUpsert {
	u.Set(using.FieldDeviceID, v)
	return u
}

// UpdateDeviceID sets the "device_id" field to the value that was provided on create.
func (u *UsingUpsert) UpdateDeviceID() *UsingUpsert {
	u.SetExcluded(using.FieldDeviceID)
	return u
}

// SetAlg sets the "alg" field.
func (u *UsingUpsert) SetAlg(v string) *UsingUpsert {
	u.Set(using.FieldAlg, v)
	return u
}

// UpdateAlg sets the "alg" field to the value that was provided on create.
func (u *UsingUpsert) UpdateAlg() *UsingUpsert {
	u.SetExcluded(using.FieldAlg)
	return u
}

// SetName sets the "name" field.
func (u *UsingUpsert) SetName(v string) *UsingUpsert {
	u.Set(using.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *UsingUpsert) UpdateName() *UsingUpsert {
	u.SetExcluded(using.FieldName)
	return u
}

// ClearName clears the value of the "name" field.
func (u *UsingUpsert) ClearName() *UsingUpsert {
	u.SetNull(using.FieldName)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.Using.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *UsingUpsertOne) UpdateNewValues() *UsingUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.CreateTime(); exists {
			s.SetIgnore(using.FieldCreateTime)
		}
		if _, exists := u.create.mutation.Sn(); exists {
			s.SetIgnore(using.FieldSn)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Using.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *UsingUpsertOne) Ignore() *UsingUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *UsingUpsertOne) DoNothing() *UsingUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the UsingCreate.OnConflict
// documentation for more info.
func (u *UsingUpsertOne) Update(set func(*UsingUpsert)) *UsingUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&UsingUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *UsingUpsertOne) SetUpdateTime(v time.Time) *UsingUpsertOne {
	return u.Update(func(s *UsingUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *UsingUpsertOne) UpdateUpdateTime() *UsingUpsertOne {
	return u.Update(func(s *UsingUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetUUID sets the "uuid" field.
func (u *UsingUpsertOne) SetUUID(v string) *UsingUpsertOne {
	return u.Update(func(s *UsingUpsert) {
		s.SetUUID(v)
	})
}

// UpdateUUID sets the "uuid" field to the value that was provided on create.
func (u *UsingUpsertOne) UpdateUUID() *UsingUpsertOne {
	return u.Update(func(s *UsingUpsert) {
		s.UpdateUUID()
	})
}

// SetDeviceID sets the "device_id" field.
func (u *UsingUpsertOne) SetDeviceID(v string) *UsingUpsertOne {
	return u.Update(func(s *UsingUpsert) {
		s.SetDeviceID(v)
	})
}

// UpdateDeviceID sets the "device_id" field to the value that was provided on create.
func (u *UsingUpsertOne) UpdateDeviceID() *UsingUpsertOne {
	return u.Update(func(s *UsingUpsert) {
		s.UpdateDeviceID()
	})
}

// SetAlg sets the "alg" field.
func (u *UsingUpsertOne) SetAlg(v string) *UsingUpsertOne {
	return u.Update(func(s *UsingUpsert) {
		s.SetAlg(v)
	})
}

// UpdateAlg sets the "alg" field to the value that was provided on create.
func (u *UsingUpsertOne) UpdateAlg() *UsingUpsertOne {
	return u.Update(func(s *UsingUpsert) {
		s.UpdateAlg()
	})
}

// SetName sets the "name" field.
func (u *UsingUpsertOne) SetName(v string) *UsingUpsertOne {
	return u.Update(func(s *UsingUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *UsingUpsertOne) UpdateName() *UsingUpsertOne {
	return u.Update(func(s *UsingUpsert) {
		s.UpdateName()
	})
}

// ClearName clears the value of the "name" field.
func (u *UsingUpsertOne) ClearName() *UsingUpsertOne {
	return u.Update(func(s *UsingUpsert) {
		s.ClearName()
	})
}

// Exec executes the query.
func (u *UsingUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for UsingCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *UsingUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *UsingUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *UsingUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// UsingCreateBulk is the builder for creating many Using entities in bulk.
type UsingCreateBulk struct {
	config
	err      error
	builders []*UsingCreate
	conflict []sql.ConflictOption
}

// Save creates the Using entities in the database.
func (ucb *UsingCreateBulk) Save(ctx context.Context) ([]*Using, error) {
	if ucb.err != nil {
		return nil, ucb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(ucb.builders))
	nodes := make([]*Using, len(ucb.builders))
	mutators := make([]Mutator, len(ucb.builders))
	for i := range ucb.builders {
		func(i int, root context.Context) {
			builder := ucb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*UsingMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, ucb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = ucb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ucb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, ucb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ucb *UsingCreateBulk) SaveX(ctx context.Context) []*Using {
	v, err := ucb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ucb *UsingCreateBulk) Exec(ctx context.Context) error {
	_, err := ucb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ucb *UsingCreateBulk) ExecX(ctx context.Context) {
	if err := ucb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Using.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.UsingUpsert) {
//			SetCreateTime(v+v).
//		}).
//		Exec(ctx)
func (ucb *UsingCreateBulk) OnConflict(opts ...sql.ConflictOption) *UsingUpsertBulk {
	ucb.conflict = opts
	return &UsingUpsertBulk{
		create: ucb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Using.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (ucb *UsingCreateBulk) OnConflictColumns(columns ...string) *UsingUpsertBulk {
	ucb.conflict = append(ucb.conflict, sql.ConflictColumns(columns...))
	return &UsingUpsertBulk{
		create: ucb,
	}
}

// UsingUpsertBulk is the builder for "upsert"-ing
// a bulk of Using nodes.
type UsingUpsertBulk struct {
	create *UsingCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Using.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *UsingUpsertBulk) UpdateNewValues() *UsingUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.CreateTime(); exists {
				s.SetIgnore(using.FieldCreateTime)
			}
			if _, exists := b.mutation.Sn(); exists {
				s.SetIgnore(using.FieldSn)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Using.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *UsingUpsertBulk) Ignore() *UsingUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *UsingUpsertBulk) DoNothing() *UsingUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the UsingCreateBulk.OnConflict
// documentation for more info.
func (u *UsingUpsertBulk) Update(set func(*UsingUpsert)) *UsingUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&UsingUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdateTime sets the "update_time" field.
func (u *UsingUpsertBulk) SetUpdateTime(v time.Time) *UsingUpsertBulk {
	return u.Update(func(s *UsingUpsert) {
		s.SetUpdateTime(v)
	})
}

// UpdateUpdateTime sets the "update_time" field to the value that was provided on create.
func (u *UsingUpsertBulk) UpdateUpdateTime() *UsingUpsertBulk {
	return u.Update(func(s *UsingUpsert) {
		s.UpdateUpdateTime()
	})
}

// SetUUID sets the "uuid" field.
func (u *UsingUpsertBulk) SetUUID(v string) *UsingUpsertBulk {
	return u.Update(func(s *UsingUpsert) {
		s.SetUUID(v)
	})
}

// UpdateUUID sets the "uuid" field to the value that was provided on create.
func (u *UsingUpsertBulk) UpdateUUID() *UsingUpsertBulk {
	return u.Update(func(s *UsingUpsert) {
		s.UpdateUUID()
	})
}

// SetDeviceID sets the "device_id" field.
func (u *UsingUpsertBulk) SetDeviceID(v string) *UsingUpsertBulk {
	return u.Update(func(s *UsingUpsert) {
		s.SetDeviceID(v)
	})
}

// UpdateDeviceID sets the "device_id" field to the value that was provided on create.
func (u *UsingUpsertBulk) UpdateDeviceID() *UsingUpsertBulk {
	return u.Update(func(s *UsingUpsert) {
		s.UpdateDeviceID()
	})
}

// SetAlg sets the "alg" field.
func (u *UsingUpsertBulk) SetAlg(v string) *UsingUpsertBulk {
	return u.Update(func(s *UsingUpsert) {
		s.SetAlg(v)
	})
}

// UpdateAlg sets the "alg" field to the value that was provided on create.
func (u *UsingUpsertBulk) UpdateAlg() *UsingUpsertBulk {
	return u.Update(func(s *UsingUpsert) {
		s.UpdateAlg()
	})
}

// SetName sets the "name" field.
func (u *UsingUpsertBulk) SetName(v string) *UsingUpsertBulk {
	return u.Update(func(s *UsingUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *UsingUpsertBulk) UpdateName() *UsingUpsertBulk {
	return u.Update(func(s *UsingUpsert) {
		s.UpdateName()
	})
}

// ClearName clears the value of the "name" field.
func (u *UsingUpsertBulk) ClearName() *UsingUpsertBulk {
	return u.Update(func(s *UsingUpsert) {
		s.ClearName()
	})
}

// Exec executes the query.
func (u *UsingUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the UsingCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for UsingCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *UsingUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
