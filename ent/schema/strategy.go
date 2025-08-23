package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Strategy holds the schema definition for the Strategy entity.
type Strategy struct {
	ent.Schema
}

// Fields of the Strategy.
func (Strategy) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("strategy_id").
			MaxLen(50),
		field.UUID("template_id", uuid.UUID{}).
			Optional().
			Nillable(),
		field.String("name").
			MaxLen(100),
		field.Text("description").
			Optional().
			Nillable(),
		field.String("symbol").
			MaxLen(10),
		field.JSON("user_inputs", map[string]interface{}{}).
			Default(map[string]interface{}{}),
		field.JSON("settings", map[string]interface{}{}).
			Default(map[string]interface{}{}),
		field.Bool("active").
			Default(false),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Optional().
			Nillable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Optional().
			Nillable(),
	}
}

// Edges of the Strategy.
func (Strategy) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("strategies").
			Field("user_id").
			Required().
			Unique(),
		edge.From("template", StrategyTemplate.Type).
			Ref("strategies").
			Field("template_id").
			Unique(),
		edge.To("executions", StrategyExecution.Type),
		edge.To("performance", StrategyPerformance.Type).
			Unique(),
		edge.To("status", StrategyStatus.Type).
			Unique(),
	}
}

// Indexes of the Strategy.
func (Strategy) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("symbol"),
		index.Fields("template_id"),
		index.Fields("active"),
		index.Fields("user_id", "strategy_id").
			Unique(),
	}
}
