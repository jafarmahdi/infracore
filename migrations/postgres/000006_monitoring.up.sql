-- 000006_monitoring.up.sql
-- Monitoring system schema

-- ============================================================
-- MONITORING PROFILES  (check interval templates)
-- ============================================================
CREATE TABLE monitoring_profiles (
    id                          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id                   UUID        NOT NULL REFERENCES tenants(id),
    name                        VARCHAR(255) NOT NULL,
    description                 TEXT,
    check_interval_seconds      INTEGER      NOT NULL DEFAULT 60   CHECK (check_interval_seconds >= 10),
    retry_interval_seconds      INTEGER      NOT NULL DEFAULT 30   CHECK (retry_interval_seconds >= 5),
    max_retries                 INTEGER      NOT NULL DEFAULT 3    CHECK (max_retries BETWEEN 0 AND 10),
    notification_period         VARCHAR(100) DEFAULT '24x7',
    is_default                  BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at                  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, name)
);

CREATE TRIGGER set_monitoring_profiles_updated_at
    BEFORE UPDATE ON monitoring_profiles
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- SSH KEYS  (for SSH-based monitoring)
-- ============================================================
CREATE TABLE ssh_keys (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    name            VARCHAR(255) NOT NULL,
    public_key      TEXT        NOT NULL,
    private_key_enc TEXT        NOT NULL,   -- AES-256-GCM encrypted
    fingerprint     VARCHAR(255),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_ssh_keys_updated_at
    BEFORE UPDATE ON ssh_keys
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- MONITORED HOSTS
-- ============================================================
CREATE TABLE monitored_hosts (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID        NOT NULL REFERENCES tenants(id),
    asset_id            UUID        REFERENCES assets(id),
    agent_id            UUID,       -- FK added after agents migration
    profile_id          UUID        REFERENCES monitoring_profiles(id),
    name                VARCHAR(255) NOT NULL,
    display_name        VARCHAR(255),
    ip_address          INET,
    hostname            VARCHAR(255),
    monitoring_type     VARCHAR(20)  NOT NULL DEFAULT 'agent'
                            CHECK (monitoring_type IN ('agent','snmp','wmi','icmp','ssh','api')),

    -- SNMP configuration
    snmp_version        VARCHAR(5)
                            CHECK (snmp_version IN ('v1','v2c','v3')),
    snmp_community      VARCHAR(255),
    snmp_port           INTEGER      DEFAULT 161,
    snmp_timeout_secs   INTEGER      DEFAULT 5,
    snmp_v3_username    VARCHAR(255),
    snmp_v3_auth_proto  VARCHAR(10)
                            CHECK (snmp_v3_auth_proto IN ('MD5','SHA','SHA256','SHA512')),
    snmp_v3_priv_proto  VARCHAR(10)
                            CHECK (snmp_v3_priv_proto IN ('DES','AES','AES192','AES256')),
    snmp_v3_sec_level   VARCHAR(20)
                            CHECK (snmp_v3_sec_level IN ('noAuthNoPriv','authNoPriv','authPriv')),

    -- SSH configuration
    ssh_port            INTEGER      DEFAULT 22,
    ssh_username        VARCHAR(100),
    ssh_key_id          UUID         REFERENCES ssh_keys(id),

    -- WMI configuration
    wmi_username        VARCHAR(255),
    wmi_domain          VARCHAR(255),

    -- State
    status              VARCHAR(50)  NOT NULL DEFAULT 'active'
                            CHECK (status IN ('active','maintenance','disabled','pending')),
    maintenance_start   TIMESTAMPTZ,
    maintenance_end     TIMESTAMPTZ,
    last_check_at       TIMESTAMPTZ,
    last_status         VARCHAR(20)
                            CHECK (last_status IN ('up','down','warning','unknown','maintenance')),
    uptime_percent      DECIMAL(6,3),

    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE TRIGGER set_monitored_hosts_updated_at
    BEFORE UPDATE ON monitored_hosts
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- CHECK DEFINITIONS  (what to check on each host)
-- ============================================================
CREATE TABLE check_definitions (
    id                      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID        NOT NULL REFERENCES tenants(id),
    host_id                 UUID        NOT NULL REFERENCES monitored_hosts(id) ON DELETE CASCADE,
    name                    VARCHAR(255) NOT NULL,
    check_type              VARCHAR(100) NOT NULL
                                CHECK (check_type IN ('cpu','ram','disk','network','service',
                                                      'process','icmp','http','https','tcp','udp',
                                                      'dns','smtp','ftp','snmp_oid','custom')),
    parameters              JSONB        NOT NULL DEFAULT '{}',
    -- Examples:
    -- cpu:      {"warning": 80, "critical": 95}
    -- disk:     {"mount": "/", "warning": 80, "critical": 90}
    -- http:     {"url": "http://app/health", "expected_code": 200, "timeout": 10}
    -- snmp_oid: {"oid": "1.3.6.1.2.1.1.1.0", "warning": 80, "critical": 90}
    -- service:  {"name": "nginx", "check_running": true}
    thresholds              JSONB        NOT NULL DEFAULT '{}',
    check_interval_seconds  INTEGER,    -- NULL = inherit from profile
    is_enabled              BOOLEAN      NOT NULL DEFAULT TRUE,
    last_value              DOUBLE PRECISION,
    last_status             VARCHAR(20),
    last_checked_at         TIMESTAMPTZ,
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_check_definitions_updated_at
    BEFORE UPDATE ON check_definitions
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- METRIC DATA  (TimescaleDB hypertable — raw time-series)
-- ============================================================
CREATE TABLE metric_data (
    time        TIMESTAMPTZ     NOT NULL,
    tenant_id   UUID            NOT NULL,
    host_id     UUID            NOT NULL,
    check_id    UUID            NOT NULL,
    metric_name VARCHAR(100)    NOT NULL,
    value       DOUBLE PRECISION NOT NULL,
    labels      JSONB           NOT NULL DEFAULT '{}'
);

-- Convert to TimescaleDB hypertable, partition by week
SELECT create_hypertable(
    'metric_data', 'time',
    chunk_time_interval => INTERVAL '1 week',
    if_not_exists => TRUE
);

-- Compress chunks older than 30 days
SELECT add_compression_policy('metric_data', INTERVAL '30 days', if_not_exists => TRUE);

-- Auto-drop chunks older than 1 year (configurable per tenant)
SELECT add_retention_policy('metric_data', INTERVAL '365 days', if_not_exists => TRUE);

-- ============================================================
-- AVAILABILITY RECORDS  (host up/down event log)
-- ============================================================
CREATE TABLE availability_records (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL,
    host_id         UUID        NOT NULL REFERENCES monitored_hosts(id),
    status          VARCHAR(20)  NOT NULL
                        CHECK (status IN ('up','down','warning','unknown','maintenance')),
    started_at      TIMESTAMPTZ  NOT NULL,
    ended_at        TIMESTAMPTZ,
    duration_seconds INTEGER     GENERATED ALWAYS AS
                        (EXTRACT(EPOCH FROM (ended_at - started_at))::INTEGER) STORED,
    reason          TEXT
);

-- ============================================================
-- SERVICE CHECKS  (application-level service monitoring)
-- ============================================================
CREATE TABLE service_checks (
    id                      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID        NOT NULL REFERENCES tenants(id),
    host_id                 UUID        NOT NULL REFERENCES monitored_hosts(id) ON DELETE CASCADE,
    service_name            VARCHAR(255) NOT NULL,
    service_display_name    VARCHAR(255),
    check_type              VARCHAR(50)  NOT NULL
                                CHECK (check_type IN ('tcp','http','https','dns','ftp',
                                                      'smtp','imap','pop3','ldap','ssh','custom')),
    target_host             VARCHAR(255),
    target_port             INTEGER      CHECK (target_port BETWEEN 1 AND 65535),
    expected_response       TEXT,
    expected_http_code      INTEGER,
    timeout_seconds         INTEGER      NOT NULL DEFAULT 10,
    follow_redirects        BOOLEAN      NOT NULL DEFAULT TRUE,
    is_enabled              BOOLEAN      NOT NULL DEFAULT TRUE,
    last_status             VARCHAR(20)
                                CHECK (last_status IN ('ok','warning','critical','unknown')),
    last_checked_at         TIMESTAMPTZ,
    last_response_time_ms   INTEGER,
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_service_checks_updated_at
    BEFORE UPDATE ON service_checks
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
