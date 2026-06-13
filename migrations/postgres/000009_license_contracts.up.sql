-- 000009_license_contracts.up.sql
-- License Management and Contract Management schemas

-- ============================================================
-- SOFTWARE LICENSES
-- ============================================================
CREATE TABLE software_licenses (
    id                      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID        NOT NULL REFERENCES tenants(id),
    site_id                 UUID        REFERENCES sites(id),
    department_id           UUID        REFERENCES departments(id),
    vendor_id               UUID        REFERENCES manufacturers(id),
    name                    VARCHAR(255) NOT NULL,
    product_name            VARCHAR(255) NOT NULL,
    version                 VARCHAR(100),
    license_type            VARCHAR(100)
                                CHECK (license_type IN ('perpetual','subscription','concurrent',
                                                        'per_seat','per_device','per_core',
                                                        'per_server','open_source','freeware','trial')),
    license_key             TEXT,       -- stored encrypted
    license_count           INTEGER      NOT NULL DEFAULT 1 CHECK (license_count > 0),
    seats_used              INTEGER      NOT NULL DEFAULT 0 CHECK (seats_used >= 0),
    purchase_date           DATE,
    purchase_cost           DECIMAL(12,2),
    purchase_currency       CHAR(3)      DEFAULT 'USD',
    purchase_order          VARCHAR(100),
    invoice_number          VARCHAR(100),
    expiration_date         DATE,
    support_expiration_date DATE,
    renewal_cost            DECIMAL(12,2),
    notes                   TEXT,
    status                  VARCHAR(50)  NOT NULL DEFAULT 'active'
                                CHECK (status IN ('active','expired','cancelled','pending')),
    notify_days_before      INTEGER[]    NOT NULL DEFAULT '{30,14,7}',  -- alert N days before expiry
    document_urls           JSONB        NOT NULL DEFAULT '[]',
    tags                    JSONB        NOT NULL DEFAULT '[]',
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ
);

CREATE TRIGGER set_software_licenses_updated_at
    BEFORE UPDATE ON software_licenses
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- LICENSE ASSIGNMENTS  (license → asset or user)
-- ============================================================
CREATE TABLE license_assignments (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           UUID        NOT NULL REFERENCES tenants(id),
    license_id          UUID        NOT NULL REFERENCES software_licenses(id) ON DELETE CASCADE,
    assigned_to_type    VARCHAR(20)  NOT NULL CHECK (assigned_to_type IN ('asset','user')),
    assigned_to_id      UUID        NOT NULL,
    license_key_segment TEXT,           -- specific key/seat assigned
    assigned_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    assigned_by         UUID        REFERENCES users(id),
    notes               TEXT,
    UNIQUE (license_id, assigned_to_type, assigned_to_id)
);

-- ============================================================
-- CONTRACTS
-- ============================================================
CREATE TABLE contracts (
    id                      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID        NOT NULL REFERENCES tenants(id),
    site_id                 UUID        REFERENCES sites(id),
    department_id           UUID        REFERENCES departments(id),
    vendor_id               UUID        REFERENCES manufacturers(id),
    contract_number         VARCHAR(100),
    title                   VARCHAR(255) NOT NULL,
    contract_type           VARCHAR(100) NOT NULL
                                CHECK (contract_type IN ('maintenance','support','warranty',
                                                         'service','lease','nda','sla','other')),
    description             TEXT,
    start_date              DATE         NOT NULL,
    end_date                DATE,
    auto_renew              BOOLEAN      NOT NULL DEFAULT FALSE,
    renewal_notice_days     INTEGER      DEFAULT 30,
    value                   DECIMAL(12,2),
    currency                CHAR(3)      DEFAULT 'USD',
    billing_cycle           VARCHAR(20)
                                CHECK (billing_cycle IN ('monthly','quarterly',
                                                         'annually','one_time')),
    status                  VARCHAR(50)  NOT NULL DEFAULT 'active'
                                CHECK (status IN ('draft','active','expired',
                                                  'cancelled','renewed','under_review')),
    sla_response_time_hours INTEGER,
    sla_resolution_time_hours INTEGER,
    contact_name            VARCHAR(255),
    contact_email           CITEXT,
    contact_phone           VARCHAR(50),
    notes                   TEXT,
    document_urls           JSONB        NOT NULL DEFAULT '[]',
    notify_days_before      INTEGER[]    NOT NULL DEFAULT '{30,14,7}',
    tags                    JSONB        NOT NULL DEFAULT '[]',
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ
);

CREATE TRIGGER set_contracts_updated_at
    BEFORE UPDATE ON contracts
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- CONTRACT ↔ ASSETS  (M:N)
-- ============================================================
CREATE TABLE contract_assets (
    contract_id UUID NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
    asset_id    UUID NOT NULL REFERENCES assets(id)    ON DELETE CASCADE,
    PRIMARY KEY (contract_id, asset_id)
);

-- ============================================================
-- WARRANTIES  (can also be tracked via contracts, but dedicated table
--              allows quick asset-centric warranty lookup)
-- ============================================================
CREATE TABLE warranties (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    asset_id        UUID        NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    vendor_id       UUID        REFERENCES manufacturers(id),
    contract_id     UUID        REFERENCES contracts(id),
    warranty_type   VARCHAR(50)
                        CHECK (warranty_type IN ('manufacturer','extended',
                                                 'accidental','on_site','next_day')),
    start_date      DATE        NOT NULL,
    end_date        DATE        NOT NULL,
    description     TEXT,
    support_contact TEXT,
    rma_process     TEXT,
    notify_days_before INTEGER[] NOT NULL DEFAULT '{30,14}',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_warranties_updated_at
    BEFORE UPDATE ON warranties
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();
