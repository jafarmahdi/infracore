-- 000010_discovery.up.sql
-- Network & Infrastructure Discovery schema

-- ============================================================
-- DISCOVERY JOBS
-- ============================================================
CREATE TABLE discovery_jobs (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    site_id         UUID        REFERENCES sites(id),
    agent_id        UUID        REFERENCES agents(id),
    name            VARCHAR(255) NOT NULL,
    job_type        VARCHAR(50)  NOT NULL
                        CHECK (job_type IN ('network','snmp','vmware','hyperv',
                                            'proxmox','active_directory','nmap')),
    status          VARCHAR(50)  NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending','queued','running',
                                          'completed','failed','cancelled')),
    -- Cron expression (NULL = run once)
    schedule_cron   VARCHAR(100),
    next_run_at     TIMESTAMPTZ,
    last_run_at     TIMESTAMPTZ,
    -- job_type-specific configuration:
    -- network: {"network": "10.0.0.0/24", "ports": [22,80,443,161,3389], "timeout": 5}
    -- snmp:    {"communities": ["public","private"], "versions": ["v2c","v3"]}
    -- vmware:  {"vcenter_host": "vc01.local", "username": "...", "password_enc": "..."}
    -- ad:      {"domain": "corp.local", "ldap_server": "dc01.corp.local"}
    configuration   JSONB        NOT NULL DEFAULT '{}',
    progress_percent INTEGER     NOT NULL DEFAULT 0 CHECK (progress_percent BETWEEN 0 AND 100),
    discovered_count INTEGER     NOT NULL DEFAULT 0,
    imported_count   INTEGER     NOT NULL DEFAULT 0,
    error_message   TEXT,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_discovery_jobs_updated_at
    BEFORE UPDATE ON discovery_jobs
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- DISCOVERY RESULTS  (one row per discovered device)
-- ============================================================
CREATE TABLE discovery_results (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id          UUID        NOT NULL REFERENCES discovery_jobs(id) ON DELETE CASCADE,
    tenant_id       UUID        NOT NULL,
    ip_address      INET,
    hostname        VARCHAR(255),
    fqdn            VARCHAR(255),
    mac_address     MACADDR,
    device_type     VARCHAR(100),   -- inferred: 'switch', 'server', 'printer'
    vendor          VARCHAR(255),   -- from OUI lookup or SNMP sysDescr
    model           VARCHAR(255),
    os_info         VARCHAR(255),
    open_ports      JSONB        NOT NULL DEFAULT '[]',    -- [{port: 22, service: "ssh"}]
    snmp_data       JSONB        NOT NULL DEFAULT '{}',   -- {sysDescr, sysName, interfaces: [...]}
    vmware_data     JSONB        NOT NULL DEFAULT '{}',
    ad_data         JSONB        NOT NULL DEFAULT '{}',
    raw_data        JSONB        NOT NULL DEFAULT '{}',
    status          VARCHAR(50)  NOT NULL DEFAULT 'new'
                        CHECK (status IN ('new','imported','ignored','duplicate','pending_review')),
    matched_asset_id UUID        REFERENCES assets(id),  -- if matched to existing asset
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ============================================================
-- DISCOVERY RULES  (auto-import / auto-assign rules)
-- ============================================================
CREATE TABLE discovery_rules (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    is_enabled      BOOLEAN      NOT NULL DEFAULT TRUE,
    priority        INTEGER      NOT NULL DEFAULT 100,   -- lower = higher priority
    -- Matching conditions (all must match = AND logic)
    -- {"ip_range": "10.0.1.0/24", "vendor": "Cisco", "open_ports": [161], "hostname_pattern": "sw-*"}
    match_conditions JSONB       NOT NULL DEFAULT '{}',
    -- Actions to take on match
    -- {"auto_import": true, "category": "switch", "site_id": "uuid", "role": "core-switch", "tags": ["auto"]}
    actions         JSONB        NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_discovery_rules_updated_at
    BEFORE UPDATE ON discovery_rules
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
