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

// CommissionLog holds the schema definition for the CommissionLog entity.
//
// 佣金日志：记录每次被邀请人充值时产生的佣金事件
// 删除策略：永久保留（无软删除），佣金记录为财务凭证，不可删除
type CommissionLog struct {
	ent.Schema
}

func (CommissionLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "commission_logs"},
	}
}

func (CommissionLog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (CommissionLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("referrer_id").
			Comment("推荐人用户ID"),
		field.Int("invitee_id").
			Comment("被邀请人用户ID"),
		field.String("order_id").
			MaxLen(64).
			Comment("触发本次佣金的支付订单ID"),
		field.Float("recharge_amount").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Comment("被邀请人充值金额"),
		field.Float("commission_rate").
			SchemaType(map[string]string{dialect.Postgres: "decimal(10,4)"}).
			Comment("佣金比例（如 0.10 表示 10%）"),
		field.Float("commission_amount").
			SchemaType(map[string]string{dialect.Postgres: "decimal(20,8)"}).
			Comment("实际佣金金额"),
		field.String("status").
			MaxLen(20).
			Default("completed").
			Comment("状态：completed / cancelled / pending"),
		field.String("remark").
			SchemaType(map[string]string{dialect.Postgres: "text"}).
			Optional().
			Default("").
			Comment("备注"),
	}
}

func (CommissionLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("referrer", User.Type).
			Ref("referral_commission_logs").
			Field("referrer_id").
			Unique().
			Required(),
		edge.From("invitee", User.Type).
			Ref("invitee_commission_logs").
			Field("invitee_id").
			Unique().
			Required(),
	}
}

func (CommissionLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("referrer_id"),
		index.Fields("invitee_id"),
		index.Fields("order_id"),
	}
}
