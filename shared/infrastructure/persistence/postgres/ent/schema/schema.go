package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("name").Optional().MaxLen(256),
		field.String("password").MaxLen(30),
		field.Bool("is_manager").Default(false),
		field.String("photo_url").Optional().MaxLen(500),
		field.String("document").MaxLen(50).Unique(),
		field.String("email").Optional().MaxLen(256),
		field.UUID("manager_id", uuid.UUID{}).Optional(),
		field.Int("user_status_id").Optional(),
		field.Time("deactivated_at").Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("manager", User.Type).
			Ref("subordinates").
			Field("manager_id").
			Unique(),

		edge.To("subordinates", User.Type),

		edge.To("business", Business.Type).Unique(),

		edge.From("status", UserStatus.Type).
			Ref("users").
			Field("user_status_id").
			Unique(),

		edge.To("access_groups", UsersOnAccessGroups.Type),
	}
}

type Business struct {
	ent.Schema
}

func (Business) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Positive().Immutable(),
		field.String("name").Unique().MaxLen(60),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deactivated_at").Optional(),
	}
}

func (Business) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
	}
}

type AccessGroup struct {
	ent.Schema
}

func (AccessGroup) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Positive().Immutable(),
		field.String("name").MaxLen(100),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deactivated_at").Optional(),
	}
}

func (AccessGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", UsersOnAccessGroups.Type),
	}
}

type UsersOnAccessGroups struct {
	ent.Schema
}

func (UsersOnAccessGroups) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("user_id", uuid.UUID{}),
		field.Int("access_group_id"),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

func (UsersOnAccessGroups) Edges() []ent.Edge {
	return []ent.Edge{
		// Edge para User
		edge.From("user", User.Type).
			Ref("access_groups").
			Field("user_id").
			Unique().
			Required(),

		edge.From("access_group", AccessGroup.Type).
			Ref("users").
			Field("access_group_id").
			Unique().
			Required(),
	}
}

func (UsersOnAccessGroups) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "access_group_id").Unique(),
	}
}

type UserStatus struct {
	ent.Schema
}

func (UserStatus) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Positive().Immutable(),
		field.String("external_id").Unique().MaxLen(5),
		field.String("name").Unique().MaxLen(60),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (UserStatus) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
	}
}
