-- Add referral/commission system tables and user fields.
-- This migration supports the tiered commission referral feature.

-- Add referral fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS referral_code VARCHAR(32) NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS referred_by BIGINT;

-- Index for referral code lookup
CREATE INDEX IF NOT EXISTS idx_users_referral_code ON users (referral_code) WHERE referral_code != '';

-- Referral records: tracks who invited whom
CREATE TABLE IF NOT EXISTS referral_records (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    referrer_id BIGINT NOT NULL REFERENCES users(id),
    invitee_id BIGINT NOT NULL UNIQUE REFERENCES users(id),
    referral_code VARCHAR(32) NOT NULL,
    total_commission DECIMAL(20,8) NOT NULL DEFAULT 0,
    total_recharged DECIMAL(20,8) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
);

CREATE INDEX IF NOT EXISTS idx_referral_records_referrer_id ON referral_records (referrer_id);
CREATE INDEX IF NOT EXISTS idx_referral_records_referral_code ON referral_records (referral_code);

-- Commission logs: records each commission event
CREATE TABLE IF NOT EXISTS commission_logs (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    referrer_id BIGINT NOT NULL REFERENCES users(id),
    invitee_id BIGINT NOT NULL REFERENCES users(id),
    order_id VARCHAR(64) NOT NULL,
    recharge_amount DECIMAL(20,8) NOT NULL,
    commission_rate DECIMAL(10,4) NOT NULL,
    commission_amount DECIMAL(20,8) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'completed',
    remark TEXT DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_commission_logs_referrer_id ON commission_logs (referrer_id);
CREATE INDEX IF NOT EXISTS idx_commission_logs_invitee_id ON commission_logs (invitee_id);
CREATE INDEX IF NOT EXISTS idx_commission_logs_order_id ON commission_logs (order_id);
