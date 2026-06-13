-- 000005_ipam.up.sql
-- IP Address Management schema

-- ============================================================
-- VRFs  (Virtual Routing and Forwarding instances)
-- ============================================================
CREATE TABLE vrfs (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    name            VARCHAR(255) NOT NULL,
    rd              VARCHAR(100),       -- Route Distinguisher, e.g. '65000:1'
    description     TEXT,
    enforce_unique  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, name)
);

CREATE TRIGGER set_vrfs_updated_at
    BEFORE UPDATE ON vrfs
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- VLAN GROUPS
-- ============================================================
CREATE TABLE vlan_groups (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    site_id         UUID        REFERENCES sites(id),
    name            VARCHAR(255) NOT NULL,
    slug            VARCHAR(100) NOT NULL,
    vlan_id_ranges  JSONB        NOT NULL DEFAULT '[{"min":1,"max":4094}]',
    description     TEXT,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, slug)
);

CREATE TRIGGER set_vlan_groups_updated_at
    BEFORE UPDATE ON vlan_groups
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- VLANs
-- ============================================================
CREATE TABLE vlans (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    site_id         UUID        REFERENCES sites(id),
    vlan_group_id   UUID        REFERENCES vlan_groups(id),
    vrf_id          UUID        REFERENCES vrfs(id),
    vid             INTEGER      NOT NULL CHECK (vid BETWEEN 1 AND 4094),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    status          VARCHAR(50)  NOT NULL DEFAULT 'active'
                        CHECK (status IN ('active','reserved','deprecated')),
    role            VARCHAR(100),   -- data, voice, management, storage, native, quarantine
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE NULLS NOT DISTINCT (tenant_id, site_id, vid)
);

CREATE TRIGGER set_vlans_updated_at
    BEFORE UPDATE ON vlans
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- Back-fill VLAN FK on network_interfaces
ALTER TABLE network_interfaces
    ADD CONSTRAINT fk_ni_untagged_vlan
    FOREIGN KEY (untagged_vlan_id) REFERENCES vlans(id);

-- ============================================================
-- PREFIXES  (IP networks / subnets)
-- ============================================================
CREATE TABLE prefixes (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    site_id         UUID        REFERENCES sites(id),
    vrf_id          UUID        REFERENCES vrfs(id),
    vlan_id         UUID        REFERENCES vlans(id),
    prefix          CIDR        NOT NULL,
    family          SMALLINT    NOT NULL GENERATED ALWAYS AS
                        (CASE WHEN family(prefix) = 4 THEN 4 ELSE 6 END) STORED,
    prefix_length   SMALLINT    NOT NULL GENERATED ALWAYS AS
                        (masklen(prefix)) STORED,
    status          VARCHAR(50)  NOT NULL DEFAULT 'active'
                        CHECK (status IN ('active','reserved','deprecated','container')),
    role            VARCHAR(100),   -- loopback, link, management, container, pool
    description     TEXT,
    is_pool         BOOLEAN      NOT NULL DEFAULT FALSE,
    mark_utilized   BOOLEAN      NOT NULL DEFAULT FALSE,
    tags            JSONB        NOT NULL DEFAULT '[]',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE NULLS NOT DISTINCT (tenant_id, vrf_id, prefix)
);

CREATE TRIGGER set_prefixes_updated_at
    BEFORE UPDATE ON prefixes
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- IP ADDRESSES
-- ============================================================
CREATE TABLE ip_addresses (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID        NOT NULL REFERENCES tenants(id),
    vrf_id              UUID        REFERENCES vrfs(id),
    prefix_id           UUID        REFERENCES prefixes(id),
    address             INET        NOT NULL,       -- includes prefix length: 192.168.1.10/24
    family              SMALLINT    NOT NULL GENERATED ALWAYS AS
                            (CASE WHEN family(address) = 4 THEN 4 ELSE 6 END) STORED,
    status              VARCHAR(50)  NOT NULL DEFAULT 'active'
                            CHECK (status IN ('active','reserved','deprecated','dhcp','slaac')),
    role                VARCHAR(50)
                            CHECK (role IN ('loopback','secondary','anycast','virtual',
                                            'vip','vrrp','hsrp','glbp')),
    -- Polymorphic assignment: 'network_interface', 'virtual_machine', 'asset'
    assigned_object_type VARCHAR(100),
    assigned_object_id   UUID,
    dns_name            VARCHAR(255),
    description         TEXT,
    nat_inside_id       UUID        REFERENCES ip_addresses(id),
    tags                JSONB        NOT NULL DEFAULT '[]',
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE NULLS NOT DISTINCT (tenant_id, vrf_id, address)
);

CREATE TRIGGER set_ip_addresses_updated_at
    BEFORE UPDATE ON ip_addresses
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- DHCP LEASES
-- ============================================================
CREATE TABLE dhcp_leases (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    prefix_id       UUID        REFERENCES prefixes(id),
    ip_address      INET        NOT NULL,
    mac_address     MACADDR,
    hostname        VARCHAR(255),
    client_id       VARCHAR(255),
    asset_id        UUID        REFERENCES assets(id),
    lease_start     TIMESTAMPTZ,
    lease_end       TIMESTAMPTZ,
    is_static       BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE NULLS NOT DISTINCT (tenant_id, prefix_id, ip_address)
);

CREATE TRIGGER set_dhcp_leases_updated_at
    BEFORE UPDATE ON dhcp_leases
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- DNS ZONES
-- ============================================================
CREATE TABLE dns_zones (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    zone_type       VARCHAR(20)  NOT NULL DEFAULT 'forward'
                        CHECK (zone_type IN ('forward','reverse')),
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    soa_primary_ns  VARCHAR(255),
    soa_email       CITEXT,
    soa_serial      INTEGER,
    soa_refresh     INTEGER      DEFAULT 3600,
    soa_retry       INTEGER      DEFAULT 900,
    soa_expire      INTEGER      DEFAULT 604800,
    soa_minimum     INTEGER      DEFAULT 300,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, name)
);

CREATE TRIGGER set_dns_zones_updated_at
    BEFORE UPDATE ON dns_zones
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- DNS RECORDS
-- ============================================================
CREATE TABLE dns_records (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    zone_id         UUID        NOT NULL REFERENCES dns_zones(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,  -- relative label, '@' for zone apex
    record_type     VARCHAR(20)  NOT NULL
                        CHECK (record_type IN ('A','AAAA','CNAME','MX','TXT','PTR',
                                               'SRV','NS','CAA','NAPTR','SOA')),
    value           TEXT         NOT NULL,
    ttl             INTEGER      NOT NULL DEFAULT 3600 CHECK (ttl > 0),
    priority        INTEGER,    -- MX, SRV
    ip_address_id   UUID        REFERENCES ip_addresses(id),
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_dns_records_updated_at
    BEFORE UPDATE ON dns_records
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
