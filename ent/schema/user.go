package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.String("name").
			MaxLen(100).
			NotEmpty().
			Default("Unknown User"),
		field.String("nickname").
			MaxLen(50).
			NotEmpty().
			Default("Unknown"),
		field.String("email").
			MaxLen(255).
			NotEmpty().
			Unique(),
		field.String("password").
			MaxLen(255).
			NotEmpty().
			Default(""),
		field.Bool("is_valid").
			Default(true),
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

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("strategies", Strategy.Type),
		edge.To("portfolios", Portfolio.Type),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").
			Unique(),
		index.Fields("nickname"),
		index.Fields("is_valid"),
	}
}
