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

// Portfolio holds the schema definition for the Portfolio entity.
type Portfolio struct {
	ent.Schema
}

// Fields of the Portfolio.
func (Portfolio) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("symbol").
			MaxLen(10),
		field.Other("quantity", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(15,4)",
			}).
			Default(decimal.Zero),
		field.Other("average_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(12,4)",
			}).
			Default(decimal.Zero),
		field.Other("current_price", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(12,4)",
			}).
			Optional().
			Nillable(),
		field.Other("market_value", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(15,4)",
			}).
			Default(decimal.Zero),
		field.Other("total_cost", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(15,4)",
			}).
			Default(decimal.Zero),
		field.Other("unrealized_pnl", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(15,4)",
			}).
			Default(decimal.Zero),
		field.Other("realized_pnl", decimal.Decimal{}).
			SchemaType(map[string]string{
				"postgres": "numeric(15,4)",
			}).
			Default(decimal.Zero),
		field.Time("last_updated").
			Default(time.Now).
			UpdateDefault(time.Now),
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

// Edges of the Portfolio.
func (Portfolio) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("portfolios").
			Field("user_id").
			Required().
			Unique(),
	}
}

// Indexes of the Portfolio.
func (Portfolio) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("symbol"),
		index.Fields("user_id", "symbol").
			Unique(),
	}
}
