CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS refresh_tokens (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  ip_hash TEXT,
  ua_hash TEXT,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_revoked 
ON refresh_tokens(user_id, revoked) 
WHERE NOT revoked;

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires 
ON refresh_tokens(expires_at) 
WHERE NOT revoked;

