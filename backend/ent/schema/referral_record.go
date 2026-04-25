package schema

import (
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ReferralRecord holds the schema definition for the ReferralRecord entity.
type ReferralRecord struct {
	ent.Schema
}

func (ReferralRecord) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "referral_records"},
	}
}

func (ReferralRecord) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
		mixins.SoftDeleteMixin{},
	}
}

func (ReferralRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Int("referrer_id"),
		field.Int("invitee_id").
			Unique(),
		field.String("referral_code").
			MaxLen(32),
		field.Float("total_commission").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0),
		field.Float("total_recharged").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Default(0),
		field.String("status").
			MaxLen(20).
			Default("active"),
	}
}

func (ReferralRecord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("referrer", User.Type).
			Ref("referral_records_as_referrer").
			Field("referrer_id").
			Required().
			Unique(),
		edge.From("invitee", User.Type).
			Ref("referral_records_as_invitee").
			Field("invitee_id").
			Required().
			Unique(),
	}
}

func (ReferralRecord) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("referrer_id"),
		index.Fields("referral_code"),
	}
}
