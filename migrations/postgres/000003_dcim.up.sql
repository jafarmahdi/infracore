-- 000003_dcim.up.sql
-- Data Center Infrastructure Management

-- ============================================================
-- DATA CENTERS
-- ============================================================
CREATE TABLE data_centers (
    id                      UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID         NOT NULL REFERENCES tenants(id),
    site_id                 UUID         NOT NULL REFERENCES sites(id),
    name                    VARCHAR(255) NOT NULL,
    description             TEXT,
    facility_code           VARCHAR(50),
    physical_address        TEXT,
    noc_contact             TEXT,
    power_capacity_kva      DECIMAL(10,2),
    cooling_capacity_tons   DECIMAL(10,2),
    status                  VARCHAR(50)  NOT NULL DEFAULT 'active'
                                CHECK (status IN ('active', 'planned', 'decommissioned')),
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ
);

CREATE TRIGGER set_data_centers_updated_at
    BEFORE UPDATE ON data_centers
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- ROOMS
-- ============================================================
CREATE TABLE rooms (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    data_center_id  UUID        NOT NULL REFERENCES data_centers(id),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    floor           VARCHAR(50),
    room_number     VARCHAR(50),
    width_meters    DECIMAL(6,2),
    length_meters   DECIMAL(6,2),
    height_meters   DECIMAL(6,2),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TRIGGER set_rooms_updated_at
    BEFORE UPDATE ON rooms
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- RACKS
-- ============================================================
CREATE TABLE racks (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    site_id         UUID        NOT NULL REFERENCES sites(id),
    room_id         UUID        REFERENCES rooms(id),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    facility_id     VARCHAR(50),
    rack_type       VARCHAR(50)
                        CHECK (rack_type IN ('open-frame','enclosed','wall-mount','2-post','4-post')),
    status          VARCHAR(50)  NOT NULL DEFAULT 'active'
                        CHECK (status IN ('active','planned','reserved','decommissioned')),
    width_inches    INTEGER      NOT NULL DEFAULT 19
                        CHECK (width_inches IN (19, 21, 23)),
    total_units     INTEGER      NOT NULL DEFAULT 42
                        CHECK (total_units > 0 AND total_units <= 100),
    unit_numbering  VARCHAR(20)  NOT NULL DEFAULT 'bottom-to-top'
                        CHECK (unit_numbering IN ('bottom-to-top','top-to-bottom')),
    position_x      INTEGER,
    position_y      INTEGER,
    max_weight_kg   DECIMAL(8,2),
    max_power_watts INTEGER,
    comment         TEXT,
    tags            JSONB        NOT NULL DEFAULT '[]',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TRIGGER set_racks_updated_at
    BEFORE UPDATE ON racks
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- POWER FEEDS
-- ============================================================
CREATE TABLE power_feeds (
    id                      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID        NOT NULL REFERENCES tenants(id),
    rack_id                 UUID        NOT NULL REFERENCES racks(id),
    name                    VARCHAR(255) NOT NULL,
    supply                  VARCHAR(10)  NOT NULL DEFAULT 'AC'
                                CHECK (supply IN ('AC','DC')),
    phase                   VARCHAR(20)  NOT NULL DEFAULT 'single-phase'
                                CHECK (phase IN ('single-phase','three-phase')),
    voltage                 INTEGER      NOT NULL
                                CHECK (voltage IN (100,110,120,208,220,230,240,380,400,415,480)),
    amperage                INTEGER      NOT NULL CHECK (amperage > 0),
    max_utilization_percent INTEGER      NOT NULL DEFAULT 80
                                CHECK (max_utilization_percent BETWEEN 1 AND 100),
    status                  VARCHAR(50)  NOT NULL DEFAULT 'active'
                                CHECK (status IN ('active','planned','offline','failed')),
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_power_feeds_updated_at
    BEFORE UPDATE ON power_feeds
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- PDUs  (Power Distribution Units)
-- ============================================================
CREATE TABLE pdus (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    rack_id         UUID        NOT NULL REFERENCES racks(id),
    power_feed_id   UUID        REFERENCES power_feeds(id),
    name            VARCHAR(255) NOT NULL,
    model           VARCHAR(255),
    manufacturer    VARCHAR(255),
    serial_number   VARCHAR(255),
    pdu_type        VARCHAR(50)
                        CHECK (pdu_type IN ('basic','metered','monitored','switched','smart')),
    total_outlets   INTEGER,
    amperage        INTEGER,
    voltage         INTEGER,
    status          VARCHAR(50)  NOT NULL DEFAULT 'active'
                        CHECK (status IN ('active','offline','failed')),
    ip_address      INET,
    management_url  TEXT,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TRIGGER set_pdus_updated_at
    BEFORE UPDATE ON pdus
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- PATCH PANELS
-- ============================================================
CREATE TABLE patch_panels (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    rack_id         UUID        NOT NULL REFERENCES racks(id),
    name            VARCHAR(255) NOT NULL,
    port_count      INTEGER      NOT NULL CHECK (port_count > 0),
    port_type       VARCHAR(50), -- RJ45, LC, SC, FC, SFP, QSFP
    rack_unit_start INTEGER      NOT NULL,
    rack_unit_height INTEGER     NOT NULL DEFAULT 1,
    label           VARCHAR(255),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TRIGGER set_patch_panels_updated_at
    BEFORE UPDATE ON patch_panels
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- CABLES
-- ============================================================
CREATE TABLE cables (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    cable_type      VARCHAR(50)  NOT NULL
                        CHECK (cable_type IN ('cat5e','cat6','cat6a','cat7','cat8',
                                              'mmf','smf','coax','power','dac','aoc','other')),
    status          VARCHAR(50)  NOT NULL DEFAULT 'connected'
                        CHECK (status IN ('connected','planned','decommissioned')),
    label           VARCHAR(255),
    color           VARCHAR(50),
    length_meters   DECIMAL(8,2),
    description     TEXT,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_cables_updated_at
    BEFORE UPDATE ON cables
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- CABLE TERMINATIONS  (each cable has exactly 2: A and B)
-- ============================================================
CREATE TABLE cable_terminations (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    cable_id         UUID        NOT NULL REFERENCES cables(id) ON DELETE CASCADE,
    -- polymorphic: type indicates which table the termination_id references
    termination_type VARCHAR(100) NOT NULL,  -- 'network_interface', 'patch_panel_port', 'pdu_outlet'
    termination_id   UUID        NOT NULL,
    cable_end        CHAR(1)     NOT NULL CHECK (cable_end IN ('A','B')),
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (cable_id, cable_end)
);
