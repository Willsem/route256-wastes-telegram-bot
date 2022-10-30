// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent/predicate"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent/user"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent/waste"

	"entgo.io/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeUser  = "User"
	TypeWaste = "Waste"
)

// UserMutation represents an operation that mutates the User nodes in the graph.
type UserMutation struct {
	config
	op             Op
	typ            string
	id             *int64
	first_name     *string
	last_name      *string
	user_name      *string
	waste_limit    *uint64
	addwaste_limit *int64
	clearedFields  map[string]struct{}
	wastes         map[uuid.UUID]struct{}
	removedwastes  map[uuid.UUID]struct{}
	clearedwastes  bool
	done           bool
	oldValue       func(context.Context) (*User, error)
	predicates     []predicate.User
}

var _ ent.Mutation = (*UserMutation)(nil)

// userOption allows management of the mutation configuration using functional options.
type userOption func(*UserMutation)

// newUserMutation creates new mutation for the User entity.
func newUserMutation(c config, op Op, opts ...userOption) *UserMutation {
	m := &UserMutation{
		config:        c,
		op:            op,
		typ:           TypeUser,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withUserID sets the ID field of the mutation.
func withUserID(id int64) userOption {
	return func(m *UserMutation) {
		var (
			err   error
			once  sync.Once
			value *User
		)
		m.oldValue = func(ctx context.Context) (*User, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().User.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withUser sets the old User of the mutation.
func withUser(node *User) userOption {
	return func(m *UserMutation) {
		m.oldValue = func(context.Context) (*User, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m UserMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m UserMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of User entities.
func (m *UserMutation) SetID(id int64) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *UserMutation) ID() (id int64, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *UserMutation) IDs(ctx context.Context) ([]int64, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int64{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().User.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetFirstName sets the "first_name" field.
func (m *UserMutation) SetFirstName(s string) {
	m.first_name = &s
}

// FirstName returns the value of the "first_name" field in the mutation.
func (m *UserMutation) FirstName() (r string, exists bool) {
	v := m.first_name
	if v == nil {
		return
	}
	return *v, true
}

// OldFirstName returns the old "first_name" field's value of the User entity.
// If the User object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *UserMutation) OldFirstName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldFirstName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldFirstName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldFirstName: %w", err)
	}
	return oldValue.FirstName, nil
}

// ResetFirstName resets all changes to the "first_name" field.
func (m *UserMutation) ResetFirstName() {
	m.first_name = nil
}

// SetLastName sets the "last_name" field.
func (m *UserMutation) SetLastName(s string) {
	m.last_name = &s
}

// LastName returns the value of the "last_name" field in the mutation.
func (m *UserMutation) LastName() (r string, exists bool) {
	v := m.last_name
	if v == nil {
		return
	}
	return *v, true
}

// OldLastName returns the old "last_name" field's value of the User entity.
// If the User object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *UserMutation) OldLastName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldLastName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldLastName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldLastName: %w", err)
	}
	return oldValue.LastName, nil
}

// ResetLastName resets all changes to the "last_name" field.
func (m *UserMutation) ResetLastName() {
	m.last_name = nil
}

// SetUserName sets the "user_name" field.
func (m *UserMutation) SetUserName(s string) {
	m.user_name = &s
}

// UserName returns the value of the "user_name" field in the mutation.
func (m *UserMutation) UserName() (r string, exists bool) {
	v := m.user_name
	if v == nil {
		return
	}
	return *v, true
}

// OldUserName returns the old "user_name" field's value of the User entity.
// If the User object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *UserMutation) OldUserName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUserName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUserName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUserName: %w", err)
	}
	return oldValue.UserName, nil
}

// ResetUserName resets all changes to the "user_name" field.
func (m *UserMutation) ResetUserName() {
	m.user_name = nil
}

// SetWasteLimit sets the "waste_limit" field.
func (m *UserMutation) SetWasteLimit(u uint64) {
	m.waste_limit = &u
	m.addwaste_limit = nil
}

// WasteLimit returns the value of the "waste_limit" field in the mutation.
func (m *UserMutation) WasteLimit() (r uint64, exists bool) {
	v := m.waste_limit
	if v == nil {
		return
	}
	return *v, true
}

// OldWasteLimit returns the old "waste_limit" field's value of the User entity.
// If the User object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *UserMutation) OldWasteLimit(ctx context.Context) (v *uint64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldWasteLimit is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldWasteLimit requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldWasteLimit: %w", err)
	}
	return oldValue.WasteLimit, nil
}

// AddWasteLimit adds u to the "waste_limit" field.
func (m *UserMutation) AddWasteLimit(u int64) {
	if m.addwaste_limit != nil {
		*m.addwaste_limit += u
	} else {
		m.addwaste_limit = &u
	}
}

// AddedWasteLimit returns the value that was added to the "waste_limit" field in this mutation.
func (m *UserMutation) AddedWasteLimit() (r int64, exists bool) {
	v := m.addwaste_limit
	if v == nil {
		return
	}
	return *v, true
}

// ClearWasteLimit clears the value of the "waste_limit" field.
func (m *UserMutation) ClearWasteLimit() {
	m.waste_limit = nil
	m.addwaste_limit = nil
	m.clearedFields[user.FieldWasteLimit] = struct{}{}
}

// WasteLimitCleared returns if the "waste_limit" field was cleared in this mutation.
func (m *UserMutation) WasteLimitCleared() bool {
	_, ok := m.clearedFields[user.FieldWasteLimit]
	return ok
}

// ResetWasteLimit resets all changes to the "waste_limit" field.
func (m *UserMutation) ResetWasteLimit() {
	m.waste_limit = nil
	m.addwaste_limit = nil
	delete(m.clearedFields, user.FieldWasteLimit)
}

// AddWasteIDs adds the "wastes" edge to the Waste entity by ids.
func (m *UserMutation) AddWasteIDs(ids ...uuid.UUID) {
	if m.wastes == nil {
		m.wastes = make(map[uuid.UUID]struct{})
	}
	for i := range ids {
		m.wastes[ids[i]] = struct{}{}
	}
}

// ClearWastes clears the "wastes" edge to the Waste entity.
func (m *UserMutation) ClearWastes() {
	m.clearedwastes = true
}

// WastesCleared reports if the "wastes" edge to the Waste entity was cleared.
func (m *UserMutation) WastesCleared() bool {
	return m.clearedwastes
}

// RemoveWasteIDs removes the "wastes" edge to the Waste entity by IDs.
func (m *UserMutation) RemoveWasteIDs(ids ...uuid.UUID) {
	if m.removedwastes == nil {
		m.removedwastes = make(map[uuid.UUID]struct{})
	}
	for i := range ids {
		delete(m.wastes, ids[i])
		m.removedwastes[ids[i]] = struct{}{}
	}
}

// RemovedWastes returns the removed IDs of the "wastes" edge to the Waste entity.
func (m *UserMutation) RemovedWastesIDs() (ids []uuid.UUID) {
	for id := range m.removedwastes {
		ids = append(ids, id)
	}
	return
}

// WastesIDs returns the "wastes" edge IDs in the mutation.
func (m *UserMutation) WastesIDs() (ids []uuid.UUID) {
	for id := range m.wastes {
		ids = append(ids, id)
	}
	return
}

// ResetWastes resets all changes to the "wastes" edge.
func (m *UserMutation) ResetWastes() {
	m.wastes = nil
	m.clearedwastes = false
	m.removedwastes = nil
}

// Where appends a list predicates to the UserMutation builder.
func (m *UserMutation) Where(ps ...predicate.User) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *UserMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (User).
func (m *UserMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *UserMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.first_name != nil {
		fields = append(fields, user.FieldFirstName)
	}
	if m.last_name != nil {
		fields = append(fields, user.FieldLastName)
	}
	if m.user_name != nil {
		fields = append(fields, user.FieldUserName)
	}
	if m.waste_limit != nil {
		fields = append(fields, user.FieldWasteLimit)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *UserMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case user.FieldFirstName:
		return m.FirstName()
	case user.FieldLastName:
		return m.LastName()
	case user.FieldUserName:
		return m.UserName()
	case user.FieldWasteLimit:
		return m.WasteLimit()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *UserMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case user.FieldFirstName:
		return m.OldFirstName(ctx)
	case user.FieldLastName:
		return m.OldLastName(ctx)
	case user.FieldUserName:
		return m.OldUserName(ctx)
	case user.FieldWasteLimit:
		return m.OldWasteLimit(ctx)
	}
	return nil, fmt.Errorf("unknown User field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *UserMutation) SetField(name string, value ent.Value) error {
	switch name {
	case user.FieldFirstName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFirstName(v)
		return nil
	case user.FieldLastName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLastName(v)
		return nil
	case user.FieldUserName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUserName(v)
		return nil
	case user.FieldWasteLimit:
		v, ok := value.(uint64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetWasteLimit(v)
		return nil
	}
	return fmt.Errorf("unknown User field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *UserMutation) AddedFields() []string {
	var fields []string
	if m.addwaste_limit != nil {
		fields = append(fields, user.FieldWasteLimit)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *UserMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case user.FieldWasteLimit:
		return m.AddedWasteLimit()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *UserMutation) AddField(name string, value ent.Value) error {
	switch name {
	case user.FieldWasteLimit:
		v, ok := value.(int64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddWasteLimit(v)
		return nil
	}
	return fmt.Errorf("unknown User numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *UserMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(user.FieldWasteLimit) {
		fields = append(fields, user.FieldWasteLimit)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *UserMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *UserMutation) ClearField(name string) error {
	switch name {
	case user.FieldWasteLimit:
		m.ClearWasteLimit()
		return nil
	}
	return fmt.Errorf("unknown User nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *UserMutation) ResetField(name string) error {
	switch name {
	case user.FieldFirstName:
		m.ResetFirstName()
		return nil
	case user.FieldLastName:
		m.ResetLastName()
		return nil
	case user.FieldUserName:
		m.ResetUserName()
		return nil
	case user.FieldWasteLimit:
		m.ResetWasteLimit()
		return nil
	}
	return fmt.Errorf("unknown User field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *UserMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.wastes != nil {
		edges = append(edges, user.EdgeWastes)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *UserMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case user.EdgeWastes:
		ids := make([]ent.Value, 0, len(m.wastes))
		for id := range m.wastes {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *UserMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedwastes != nil {
		edges = append(edges, user.EdgeWastes)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *UserMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case user.EdgeWastes:
		ids := make([]ent.Value, 0, len(m.removedwastes))
		for id := range m.removedwastes {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *UserMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedwastes {
		edges = append(edges, user.EdgeWastes)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *UserMutation) EdgeCleared(name string) bool {
	switch name {
	case user.EdgeWastes:
		return m.clearedwastes
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *UserMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown User unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *UserMutation) ResetEdge(name string) error {
	switch name {
	case user.EdgeWastes:
		m.ResetWastes()
		return nil
	}
	return fmt.Errorf("unknown User edge %s", name)
}

// WasteMutation represents an operation that mutates the Waste nodes in the graph.
type WasteMutation struct {
	config
	op            Op
	typ           string
	id            *uuid.UUID
	cost          *int64
	addcost       *int64
	category      *string
	date          *time.Time
	clearedFields map[string]struct{}
	user          *int64
	cleareduser   bool
	done          bool
	oldValue      func(context.Context) (*Waste, error)
	predicates    []predicate.Waste
}

var _ ent.Mutation = (*WasteMutation)(nil)

// wasteOption allows management of the mutation configuration using functional options.
type wasteOption func(*WasteMutation)

// newWasteMutation creates new mutation for the Waste entity.
func newWasteMutation(c config, op Op, opts ...wasteOption) *WasteMutation {
	m := &WasteMutation{
		config:        c,
		op:            op,
		typ:           TypeWaste,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withWasteID sets the ID field of the mutation.
func withWasteID(id uuid.UUID) wasteOption {
	return func(m *WasteMutation) {
		var (
			err   error
			once  sync.Once
			value *Waste
		)
		m.oldValue = func(ctx context.Context) (*Waste, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Waste.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withWaste sets the old Waste of the mutation.
func withWaste(node *Waste) wasteOption {
	return func(m *WasteMutation) {
		m.oldValue = func(context.Context) (*Waste, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m WasteMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m WasteMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Waste entities.
func (m *WasteMutation) SetID(id uuid.UUID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *WasteMutation) ID() (id uuid.UUID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *WasteMutation) IDs(ctx context.Context) ([]uuid.UUID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []uuid.UUID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Waste.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetCost sets the "cost" field.
func (m *WasteMutation) SetCost(i int64) {
	m.cost = &i
	m.addcost = nil
}

// Cost returns the value of the "cost" field in the mutation.
func (m *WasteMutation) Cost() (r int64, exists bool) {
	v := m.cost
	if v == nil {
		return
	}
	return *v, true
}

// OldCost returns the old "cost" field's value of the Waste entity.
// If the Waste object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *WasteMutation) OldCost(ctx context.Context) (v int64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCost is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCost requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCost: %w", err)
	}
	return oldValue.Cost, nil
}

// AddCost adds i to the "cost" field.
func (m *WasteMutation) AddCost(i int64) {
	if m.addcost != nil {
		*m.addcost += i
	} else {
		m.addcost = &i
	}
}

// AddedCost returns the value that was added to the "cost" field in this mutation.
func (m *WasteMutation) AddedCost() (r int64, exists bool) {
	v := m.addcost
	if v == nil {
		return
	}
	return *v, true
}

// ResetCost resets all changes to the "cost" field.
func (m *WasteMutation) ResetCost() {
	m.cost = nil
	m.addcost = nil
}

// SetCategory sets the "category" field.
func (m *WasteMutation) SetCategory(s string) {
	m.category = &s
}

// Category returns the value of the "category" field in the mutation.
func (m *WasteMutation) Category() (r string, exists bool) {
	v := m.category
	if v == nil {
		return
	}
	return *v, true
}

// OldCategory returns the old "category" field's value of the Waste entity.
// If the Waste object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *WasteMutation) OldCategory(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCategory is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCategory requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCategory: %w", err)
	}
	return oldValue.Category, nil
}

// ResetCategory resets all changes to the "category" field.
func (m *WasteMutation) ResetCategory() {
	m.category = nil
}

// SetDate sets the "date" field.
func (m *WasteMutation) SetDate(t time.Time) {
	m.date = &t
}

// Date returns the value of the "date" field in the mutation.
func (m *WasteMutation) Date() (r time.Time, exists bool) {
	v := m.date
	if v == nil {
		return
	}
	return *v, true
}

// OldDate returns the old "date" field's value of the Waste entity.
// If the Waste object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *WasteMutation) OldDate(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDate is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDate requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDate: %w", err)
	}
	return oldValue.Date, nil
}

// ResetDate resets all changes to the "date" field.
func (m *WasteMutation) ResetDate() {
	m.date = nil
}

// SetUserID sets the "user" edge to the User entity by id.
func (m *WasteMutation) SetUserID(id int64) {
	m.user = &id
}

// ClearUser clears the "user" edge to the User entity.
func (m *WasteMutation) ClearUser() {
	m.cleareduser = true
}

// UserCleared reports if the "user" edge to the User entity was cleared.
func (m *WasteMutation) UserCleared() bool {
	return m.cleareduser
}

// UserID returns the "user" edge ID in the mutation.
func (m *WasteMutation) UserID() (id int64, exists bool) {
	if m.user != nil {
		return *m.user, true
	}
	return
}

// UserIDs returns the "user" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// UserID instead. It exists only for internal usage by the builders.
func (m *WasteMutation) UserIDs() (ids []int64) {
	if id := m.user; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetUser resets all changes to the "user" edge.
func (m *WasteMutation) ResetUser() {
	m.user = nil
	m.cleareduser = false
}

// Where appends a list predicates to the WasteMutation builder.
func (m *WasteMutation) Where(ps ...predicate.Waste) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *WasteMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Waste).
func (m *WasteMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *WasteMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.cost != nil {
		fields = append(fields, waste.FieldCost)
	}
	if m.category != nil {
		fields = append(fields, waste.FieldCategory)
	}
	if m.date != nil {
		fields = append(fields, waste.FieldDate)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *WasteMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case waste.FieldCost:
		return m.Cost()
	case waste.FieldCategory:
		return m.Category()
	case waste.FieldDate:
		return m.Date()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *WasteMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case waste.FieldCost:
		return m.OldCost(ctx)
	case waste.FieldCategory:
		return m.OldCategory(ctx)
	case waste.FieldDate:
		return m.OldDate(ctx)
	}
	return nil, fmt.Errorf("unknown Waste field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *WasteMutation) SetField(name string, value ent.Value) error {
	switch name {
	case waste.FieldCost:
		v, ok := value.(int64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCost(v)
		return nil
	case waste.FieldCategory:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCategory(v)
		return nil
	case waste.FieldDate:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDate(v)
		return nil
	}
	return fmt.Errorf("unknown Waste field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *WasteMutation) AddedFields() []string {
	var fields []string
	if m.addcost != nil {
		fields = append(fields, waste.FieldCost)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *WasteMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case waste.FieldCost:
		return m.AddedCost()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *WasteMutation) AddField(name string, value ent.Value) error {
	switch name {
	case waste.FieldCost:
		v, ok := value.(int64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddCost(v)
		return nil
	}
	return fmt.Errorf("unknown Waste numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *WasteMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *WasteMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *WasteMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Waste nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *WasteMutation) ResetField(name string) error {
	switch name {
	case waste.FieldCost:
		m.ResetCost()
		return nil
	case waste.FieldCategory:
		m.ResetCategory()
		return nil
	case waste.FieldDate:
		m.ResetDate()
		return nil
	}
	return fmt.Errorf("unknown Waste field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *WasteMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.user != nil {
		edges = append(edges, waste.EdgeUser)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *WasteMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case waste.EdgeUser:
		if id := m.user; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *WasteMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *WasteMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *WasteMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.cleareduser {
		edges = append(edges, waste.EdgeUser)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *WasteMutation) EdgeCleared(name string) bool {
	switch name {
	case waste.EdgeUser:
		return m.cleareduser
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *WasteMutation) ClearEdge(name string) error {
	switch name {
	case waste.EdgeUser:
		m.ClearUser()
		return nil
	}
	return fmt.Errorf("unknown Waste unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *WasteMutation) ResetEdge(name string) error {
	switch name {
	case waste.EdgeUser:
		m.ResetUser()
		return nil
	}
	return fmt.Errorf("unknown Waste edge %s", name)
}