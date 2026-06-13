-- 000008_agents.up.sql
-- Agent Management schema

-- ============================================================
-- AGENT GROUPS
-- ============================================================
CREATE TABLE agent_groups (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL REFERENCES tenants(id),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    is_default  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, name)
);

CREATE TRIGGER set_agent_groups_updated_at
    BEFORE UPDATE ON agent_groups
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- AGENTS
-- ============================================================
CREATE TABLE agents (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    site_id         UUID        REFERENCES sites(id),
    group_id        UUID        REFERENCES agent_groups(id),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    hostname        VARCHAR(255),
    ip_address      INET,
    os_type         VARCHAR(50)
                        CHECK (os_type IN ('linux','windows','darwin')),
    os_version      VARCHAR(100),
    arch            VARCHAR(20)
                        CHECK (arch IN ('amd64','arm64','arm','386')),
    version         VARCHAR(50),
    -- API key is hashed (SHA-256); original shown only once at registration
    api_key_hash    VARCHAR(64)  NOT NULL UNIQUE,
    status          VARCHAR(50)  NOT NULL DEFAULT 'offline'
                        CHECK (status IN ('online','offline','error','updating','maintenance')),
    last_heartbeat_at   TIMESTAMPTZ,
    last_seen_ip        INET,
    -- Array of supported probe types
    capabilities    JSONB        NOT NULL DEFAULT '[]',
    -- Push configuration to agent via this field
    configuration   JSONB        NOT NULL DEFAULT '{}',
    registered_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_agents_updated_at
    BEFORE UPDATE ON agents
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- Back-fill FK from monitored_hosts to agents
ALTER TABLE monitored_hosts
    ADD CONSTRAINT fk_monitored_hosts_agent
    FOREIGN KEY (agent_id) REFERENCES agents(id);

-- ============================================================
-- AGENT HEALTH METRICS  (TimescaleDB hypertable)
-- ============================================================
CREATE TABLE agent_health_metrics (
    time                TIMESTAMPTZ     NOT NULL,
    agent_id            UUID            NOT NULL,
    tenant_id           UUID            NOT NULL,
    cpu_percent         DOUBLE PRECISION,
    ram_percent         DOUBLE PRECISION,
    goroutines          INTEGER,
    checks_per_minute   INTEGER,
    errors_per_minute   INTEGER,
    queue_depth         INTEGER,
    latency_ms          INTEGER         -- agent → server round-trip
);

SELECT create_hypertable(
    'agent_health_metrics', 'time',
    chunk_time_interval => INTERVAL '1 day',
    if_not_exists => TRUE
);

SELECT add_retention_policy('agent_health_metrics', INTERVAL '90 days', if_not_exists => TRUE);

-- ============================================================
-- AGENT UPDATE CHANNELS  (version distribution)
-- ============================================================
CREATE TABLE agent_update_channels (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    channel_name    VARCHAR(50)  NOT NULL DEFAULT 'stable'
                        CHECK (channel_name IN ('stable','beta','edge')),
    os_type         VARCHAR(50)  NOT NULL,
    arch            VARCHAR(20)  NOT NULL,
    version         VARCHAR(50)  NOT NULL,
    download_url    TEXT         NOT NULL,
    checksum_sha256 CHAR(64),
    release_notes   TEXT,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ============================================================
-- AGENT TASKS  (commands pushed to agents; pull-based queue)
-- ============================================================
CREATE TABLE agent_tasks (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    agent_id        UUID        NOT NULL REFERENCES agents(id),
    task_type       VARCHAR(100) NOT NULL,
    -- update_agent, run_discovery, reload_config, run_check, collect_logs
    payload         JSONB        NOT NULL DEFAULT '{}',
    status          VARCHAR(20)  NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending','delivered','running','completed','failed','expired')),
    result          JSONB,
    error_message   TEXT,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    delivered_at    TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    expires_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW() + INTERVAL '1 hour'
);
