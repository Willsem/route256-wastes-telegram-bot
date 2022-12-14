// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent/predicate"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent/user"
	"gitlab.ozon.dev/stepanov.ao.dev/telegram-bot/internal/ent/waste"
)

// WasteQuery is the builder for querying Waste entities.
type WasteQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.Waste
	withUser   *UserQuery
	withFKs    bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the WasteQuery builder.
func (wq *WasteQuery) Where(ps ...predicate.Waste) *WasteQuery {
	wq.predicates = append(wq.predicates, ps...)
	return wq
}

// Limit adds a limit step to the query.
func (wq *WasteQuery) Limit(limit int) *WasteQuery {
	wq.limit = &limit
	return wq
}

// Offset adds an offset step to the query.
func (wq *WasteQuery) Offset(offset int) *WasteQuery {
	wq.offset = &offset
	return wq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (wq *WasteQuery) Unique(unique bool) *WasteQuery {
	wq.unique = &unique
	return wq
}

// Order adds an order step to the query.
func (wq *WasteQuery) Order(o ...OrderFunc) *WasteQuery {
	wq.order = append(wq.order, o...)
	return wq
}

// QueryUser chains the current query on the "user" edge.
func (wq *WasteQuery) QueryUser() *UserQuery {
	query := &UserQuery{config: wq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := wq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := wq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(waste.Table, waste.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, waste.UserTable, waste.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(wq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Waste entity from the query.
// Returns a *NotFoundError when no Waste was found.
func (wq *WasteQuery) First(ctx context.Context) (*Waste, error) {
	nodes, err := wq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{waste.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (wq *WasteQuery) FirstX(ctx context.Context) *Waste {
	node, err := wq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Waste ID from the query.
// Returns a *NotFoundError when no Waste ID was found.
func (wq *WasteQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = wq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{waste.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (wq *WasteQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := wq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Waste entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Waste entity is found.
// Returns a *NotFoundError when no Waste entities are found.
func (wq *WasteQuery) Only(ctx context.Context) (*Waste, error) {
	nodes, err := wq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{waste.Label}
	default:
		return nil, &NotSingularError{waste.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (wq *WasteQuery) OnlyX(ctx context.Context) *Waste {
	node, err := wq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Waste ID in the query.
// Returns a *NotSingularError when more than one Waste ID is found.
// Returns a *NotFoundError when no entities are found.
func (wq *WasteQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = wq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{waste.Label}
	default:
		err = &NotSingularError{waste.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (wq *WasteQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := wq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Wastes.
func (wq *WasteQuery) All(ctx context.Context) ([]*Waste, error) {
	if err := wq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return wq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (wq *WasteQuery) AllX(ctx context.Context) []*Waste {
	nodes, err := wq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Waste IDs.
func (wq *WasteQuery) IDs(ctx context.Context) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	if err := wq.Select(waste.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (wq *WasteQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := wq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (wq *WasteQuery) Count(ctx context.Context) (int, error) {
	if err := wq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return wq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (wq *WasteQuery) CountX(ctx context.Context) int {
	count, err := wq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (wq *WasteQuery) Exist(ctx context.Context) (bool, error) {
	if err := wq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return wq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (wq *WasteQuery) ExistX(ctx context.Context) bool {
	exist, err := wq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the WasteQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (wq *WasteQuery) Clone() *WasteQuery {
	if wq == nil {
		return nil
	}
	return &WasteQuery{
		config:     wq.config,
		limit:      wq.limit,
		offset:     wq.offset,
		order:      append([]OrderFunc{}, wq.order...),
		predicates: append([]predicate.Waste{}, wq.predicates...),
		withUser:   wq.withUser.Clone(),
		// clone intermediate query.
		sql:    wq.sql.Clone(),
		path:   wq.path,
		unique: wq.unique,
	}
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (wq *WasteQuery) WithUser(opts ...func(*UserQuery)) *WasteQuery {
	query := &UserQuery{config: wq.config}
	for _, opt := range opts {
		opt(query)
	}
	wq.withUser = query
	return wq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Cost int64 `json:"cost,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Waste.Query().
//		GroupBy(waste.FieldCost).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (wq *WasteQuery) GroupBy(field string, fields ...string) *WasteGroupBy {
	grbuild := &WasteGroupBy{config: wq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := wq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return wq.sqlQuery(ctx), nil
	}
	grbuild.label = waste.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Cost int64 `json:"cost,omitempty"`
//	}
//
//	client.Waste.Query().
//		Select(waste.FieldCost).
//		Scan(ctx, &v)
func (wq *WasteQuery) Select(fields ...string) *WasteSelect {
	wq.fields = append(wq.fields, fields...)
	selbuild := &WasteSelect{WasteQuery: wq}
	selbuild.label = waste.Label
	selbuild.flds, selbuild.scan = &wq.fields, selbuild.Scan
	return selbuild
}

func (wq *WasteQuery) prepareQuery(ctx context.Context) error {
	for _, f := range wq.fields {
		if !waste.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if wq.path != nil {
		prev, err := wq.path(ctx)
		if err != nil {
			return err
		}
		wq.sql = prev
	}
	return nil
}

func (wq *WasteQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Waste, error) {
	var (
		nodes       = []*Waste{}
		withFKs     = wq.withFKs
		_spec       = wq.querySpec()
		loadedTypes = [1]bool{
			wq.withUser != nil,
		}
	)
	if wq.withUser != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, waste.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Waste).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Waste{config: wq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, wq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := wq.withUser; query != nil {
		if err := wq.loadUser(ctx, query, nodes, nil,
			func(n *Waste, e *User) { n.Edges.User = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (wq *WasteQuery) loadUser(ctx context.Context, query *UserQuery, nodes []*Waste, init func(*Waste), assign func(*Waste, *User)) error {
	ids := make([]int64, 0, len(nodes))
	nodeids := make(map[int64][]*Waste)
	for i := range nodes {
		if nodes[i].user_wastes == nil {
			continue
		}
		fk := *nodes[i].user_wastes
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	query.Where(user.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "user_wastes" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (wq *WasteQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := wq.querySpec()
	_spec.Node.Columns = wq.fields
	if len(wq.fields) > 0 {
		_spec.Unique = wq.unique != nil && *wq.unique
	}
	return sqlgraph.CountNodes(ctx, wq.driver, _spec)
}

func (wq *WasteQuery) sqlExist(ctx context.Context) (bool, error) {
	switch _, err := wq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

func (wq *WasteQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   waste.Table,
			Columns: waste.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: waste.FieldID,
			},
		},
		From:   wq.sql,
		Unique: true,
	}
	if unique := wq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := wq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, waste.FieldID)
		for i := range fields {
			if fields[i] != waste.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := wq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := wq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := wq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := wq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (wq *WasteQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(wq.driver.Dialect())
	t1 := builder.Table(waste.Table)
	columns := wq.fields
	if len(columns) == 0 {
		columns = waste.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if wq.sql != nil {
		selector = wq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if wq.unique != nil && *wq.unique {
		selector.Distinct()
	}
	for _, p := range wq.predicates {
		p(selector)
	}
	for _, p := range wq.order {
		p(selector)
	}
	if offset := wq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := wq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// WasteGroupBy is the group-by builder for Waste entities.
type WasteGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (wgb *WasteGroupBy) Aggregate(fns ...AggregateFunc) *WasteGroupBy {
	wgb.fns = append(wgb.fns, fns...)
	return wgb
}

// Scan applies the group-by query and scans the result into the given value.
func (wgb *WasteGroupBy) Scan(ctx context.Context, v any) error {
	query, err := wgb.path(ctx)
	if err != nil {
		return err
	}
	wgb.sql = query
	return wgb.sqlScan(ctx, v)
}

func (wgb *WasteGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range wgb.fields {
		if !waste.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := wgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := wgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (wgb *WasteGroupBy) sqlQuery() *sql.Selector {
	selector := wgb.sql.Select()
	aggregation := make([]string, 0, len(wgb.fns))
	for _, fn := range wgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(wgb.fields)+len(wgb.fns))
		for _, f := range wgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(wgb.fields...)...)
}

// WasteSelect is the builder for selecting fields of Waste entities.
type WasteSelect struct {
	*WasteQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (ws *WasteSelect) Scan(ctx context.Context, v any) error {
	if err := ws.prepareQuery(ctx); err != nil {
		return err
	}
	ws.sql = ws.WasteQuery.sqlQuery(ctx)
	return ws.sqlScan(ctx, v)
}

func (ws *WasteSelect) sqlScan(ctx context.Context, v any) error {
	rows := &sql.Rows{}
	query, args := ws.sql.Query()
	if err := ws.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
