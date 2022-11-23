package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Waste holds the schema definition for the Waste entity.
type Waste struct {
	ent.Schema
}

// Fields of the Waste.
func (Waste) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.Int64("cost"),
		field.String("category"),
		field.Time("date"),
	}
}

// Edges of the Waste.
func (Waste) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("wastes").
			Unique(),
	}
}

// Indexes of the Waste.
func (Waste) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("category"),
	}
}
