-- 000004_assets.up.sql
-- Asset Management schema

-- ============================================================
-- MANUFACTURERS / VENDORS
-- ============================================================
CREATE TABLE manufacturers (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL REFERENCES tenants(id),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    website_url TEXT,
    support_url TEXT,
    support_phone   VARCHAR(50),
    support_email   CITEXT,
    is_global   BOOLEAN      NOT NULL DEFAULT FALSE,  -- visible across all tenants
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, name)
);

CREATE TRIGGER set_manufacturers_updated_at
    BEFORE UPDATE ON manufacturers
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- DEVICE TYPES  (hardware model catalog / templates)
-- ============================================================
CREATE TABLE device_types (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    manufacturer_id UUID        NOT NULL REFERENCES manufacturers(id),
    name            VARCHAR(255) NOT NULL,
    model           VARCHAR(255) NOT NULL,
    part_number     VARCHAR(255),
    category        VARCHAR(100) NOT NULL
                        CHECK (category IN ('server','switch','router','firewall','ups',
                                            'printer','access_point','storage','vm','pdu',
                                            'patch_panel','kvm','other')),
    form_factor     VARCHAR(50)
                        CHECK (form_factor IN ('1u','2u','4u','blade','tower','desktop',
                                               'modular','wall-mount','handheld','other')),
    rack_units      INTEGER      CHECK (rack_units > 0 AND rack_units <= 50),
    weight_kg       DECIMAL(8,2),
    front_image_url TEXT,
    rear_image_url  TEXT,
    specs           JSONB        NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, manufacturer_id, model)
);

CREATE TRIGGER set_device_types_updated_at
    BEFORE UPDATE ON device_types
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- ASSETS  (unified single-table, discriminated by category)
-- ============================================================
CREATE TABLE assets (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID        NOT NULL REFERENCES tenants(id),
    site_id             UUID        REFERENCES sites(id),
    department_id       UUID        REFERENCES departments(id),
    device_type_id      UUID        REFERENCES device_types(id),
    rack_id             UUID        REFERENCES racks(id),
    rack_position_start INTEGER     CHECK (rack_position_start > 0),
    rack_position_end   INTEGER     CHECK (rack_position_end > 0),
    rack_face           VARCHAR(10)
                            CHECK (rack_face IN ('front','rear')),
    asset_tag           VARCHAR(100),   -- physical barcode / asset label
    serial_number       VARCHAR(255),
    name                VARCHAR(255) NOT NULL,
    category            VARCHAR(100) NOT NULL
                            CHECK (category IN ('server','switch','router','firewall','ups',
                                                'printer','access_point','storage','vm','other')),
    status              VARCHAR(50)  NOT NULL DEFAULT 'active'
                            CHECK (status IN ('active','spare','offline','failed',
                                              'decommissioned','planned','in_repair','in_transit')),
    role                VARCHAR(100),   -- e.g. 'web-server', 'core-switch', 'edge-router'

    -- Denormalized quick-display fields (avoid JOIN overhead on list views)
    manufacturer_name   VARCHAR(255),
    model_name          VARCHAR(255),

    -- Connectivity
    primary_ip          INET,
    oob_ip              INET,           -- Out-of-band management IP
    mac_address         MACADDR,

    -- OS / Platform
    os_type             VARCHAR(100),   -- 'Windows Server 2022', 'Ubuntu 22.04', 'ESXi 8.0'
    os_version          VARCHAR(100),

    -- Hardware specs (for quick listing; full specs in device_type.specs)
    cpu_count           INTEGER,
    cpu_model           VARCHAR(255),
    ram_gb              INTEGER,
    storage_gb          INTEGER,

    -- Lifecycle
    purchase_date       DATE,
    purchase_cost       DECIMAL(10,2),
    purchase_currency   CHAR(3)         DEFAULT 'USD',
    purchase_order      VARCHAR(100),
    warranty_expires_at DATE,
    eol_date            DATE,

    -- Ownership
    assigned_to_user_id UUID            REFERENCES users(id),
    managed_by_user_id  UUID            REFERENCES users(id),

    notes               TEXT,
    custom_fields       JSONB           NOT NULL DEFAULT '{}',
    tags                JSONB           NOT NULL DEFAULT '[]',

    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ
);

CREATE TRIGGER set_assets_updated_at
    BEFORE UPDATE ON assets
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- NETWORK INTERFACES  (for physical & virtual assets)
-- ============================================================
CREATE TABLE network_interfaces (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    asset_id        UUID        NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,  -- eth0, Gi0/0, vmnic0
    description     TEXT,
    interface_type  VARCHAR(50),            -- 1000base-t, 10gbase-x-sfpp, 25gbase-x-sfp28
    speed_mbps      INTEGER,
    duplex          VARCHAR(10)
                        CHECK (duplex IN ('auto','full','half')),
    mac_address     MACADDR,
    mtu             INTEGER      DEFAULT 1500,
    is_management   BOOLEAN      NOT NULL DEFAULT FALSE,
    is_enabled      BOOLEAN      NOT NULL DEFAULT TRUE,
    vlan_mode       VARCHAR(10)
                        CHECK (vlan_mode IN ('access','trunk','hybrid')),
    untagged_vlan_id    UUID,   -- FK to vlans (added in IPAM migration)
    tagged_vlan_ids     JSONB   NOT NULL DEFAULT '[]',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (asset_id, name)
);

CREATE TRIGGER set_network_interfaces_updated_at
    BEFORE UPDATE ON network_interfaces
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- VIRTUAL MACHINES  (extends assets with hypervisor metadata)
-- ============================================================
CREATE TABLE virtual_machines (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id            UUID        NOT NULL REFERENCES assets(id) ON DELETE CASCADE UNIQUE,
    tenant_id           UUID        NOT NULL REFERENCES tenants(id),
    host_asset_id       UUID        REFERENCES assets(id),       -- physical hypervisor host
    hypervisor_type     VARCHAR(50)
                            CHECK (hypervisor_type IN ('vmware','hyperv','proxmox','kvm','xen','other')),
    vm_id_external      VARCHAR(255),  -- VM's ID in the hypervisor platform
    vcpu_count          INTEGER,
    vram_gb             INTEGER,
    storage_gb          INTEGER,
    cluster_name        VARCHAR(255),
    datastore_name      VARCHAR(255),
    template_name       VARCHAR(255),
    power_state         VARCHAR(20)
                            CHECK (power_state IN ('running','stopped','suspended','paused','unknown')),
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_virtual_machines_updated_at
    BEFORE UPDATE ON virtual_machines
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- ASSET RELATIONSHIPS  (e.g. blade → chassis, vm → host)
-- ============================================================
CREATE TABLE asset_relationships (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    parent_asset_id UUID        NOT NULL REFERENCES assets(id),
    child_asset_id  UUID        NOT NULL REFERENCES assets(id),
    relationship    VARCHAR(50)  NOT NULL
                        CHECK (relationship IN ('parent_child','cluster','stack','ha_pair')),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (parent_asset_id, child_asset_id, relationship)
);
