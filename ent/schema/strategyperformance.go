package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// StrategyPerformance holds the schema definition for the StrategyPerformance entity.
type StrategyPerformance struct {
	ent.Schema
}

// Fields of the StrategyPerformance.
func (StrategyPerformance) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("strategy_id", uuid.UUID{}).
			Unique(),
		field.Other("total_return", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(15,4)",
			}).
			Default(decimal.Zero),
		field.Other("win_rate", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(5,2)",
			}).
			Default(decimal.Zero),
		field.Other("profit_loss", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(15,4)",
			}).
			Default(decimal.Zero),
		field.Int64("trade_count").
			Default(0),
		field.Time("last_trade_time").
			Optional().
			Nillable(),
		field.Other("max_drawdown", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(8,4)",
			}).
			Default(decimal.Zero),
		field.Other("sharpe_ratio", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(8,4)",
			}).
			Default(decimal.Zero),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the StrategyPerformance.
func (StrategyPerformance) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("strategy", Strategy.Type).
			Ref("performance").
			Field("strategy_id").
			Required().
			Unique(),
	}
}
