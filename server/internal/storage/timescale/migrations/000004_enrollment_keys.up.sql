CREATE TABLE enrollment_keys (
    agent_id UUID PRIMARY KEY REFERENCES agents(id),
    key_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ
);