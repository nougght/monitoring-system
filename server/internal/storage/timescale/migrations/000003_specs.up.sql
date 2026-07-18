CREATE TABLE agent_specs (
    agent_id UUID PRIMARY KEY REFERENCES agents(id),
    hostname VARCHAR(256),
    os_type VARCHAR(64),
    os VARCHAR(64),
    os_arch VARCHAR(16),
    cpu_cores_count INT,
    memory_total BIGINT,
    full_specs JSONB NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);