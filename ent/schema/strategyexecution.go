package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// StrategyExecution holds the schema definition for the StrategyExecution entity.
type StrategyExecution struct {
	ent.Schema
}

// Fields of the StrategyExecution.
func (StrategyExecution) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive(),
		field.UUID("strategy_id", uuid.UUID{}).
			Optional().
			Nillable(),
		field.String("symbol").
			MaxLen(10),
		field.Enum("action").
			Values("BUY", "SELL", "HOLD"),
		field.Other("price", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(12,4)",
			}).
			Optional().
			Nillable(),
		field.Int("quantity").
			Optional().
			Nillable(),
		field.Text("reasoning").
			Optional().
			Nillable(),
		field.Time("executed_at").
			Optional().
			Nillable().
			Default(time.Now),
	}
}

// Edges of the StrategyExecution.
func (StrategyExecution) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("strategy", Strategy.Type).
			Ref("executions").
			Field("strategy_id").
			Unique(),
	}
}

// Indexes of the StrategyExecution.
func (StrategyExecution) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("strategy_id"),
		index.Fields("symbol"),
		index.Fields("executed_at"),
	}
}
