package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// StrategyTemplate holds the schema definition for the StrategyTemplate entity.
type StrategyTemplate struct {
	ent.Schema
}

// Fields of the StrategyTemplate.
func (StrategyTemplate) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.String("name").
			MaxLen(100),
		field.Text("description").
			Optional().
			Nillable(),
		field.String("category").
			MaxLen(50),
		field.JSON("template_config", map[string]interface{}{}),
		field.JSON("input_schema", map[string]interface{}{}),
		field.String("version").
			MaxLen(10).
			Default("1.0.0"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the StrategyTemplate.
func (StrategyTemplate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("strategies", Strategy.Type),
	}
}
