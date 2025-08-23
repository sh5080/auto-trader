package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// StrategyStatus holds the schema definition for the StrategyStatus entity.
type StrategyStatus struct {
	ent.Schema
}

// Fields of the StrategyStatus.
func (StrategyStatus) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("strategy_id", uuid.UUID{}).
			Unique(),
		field.Enum("status").
			Values("active", "inactive", "paused", "error").
			Default("inactive"),
		field.Time("last_execution").
			Optional().
			Nillable(),
		field.Int64("execution_count").
			Default(0),
		field.Text("error_message").
			Optional().
			Nillable(),
		field.Int64("uptime_seconds").
			Default(0),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the StrategyStatus.
func (StrategyStatus) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("strategy", Strategy.Type).
			Ref("status").
			Field("strategy_id").
			Required().
			Unique(),
	}
}
