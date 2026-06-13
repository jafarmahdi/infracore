-- 000002_identity.up.sql
-- Identity & Access Management schema

-- ============================================================
-- TENANTS
-- ============================================================
CREATE TABLE tenants (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    slug            VARCHAR(100) NOT NULL UNIQUE,
    plan            VARCHAR(50)  NOT NULL DEFAULT 'free'
                        CHECK (plan IN ('free', 'pro', 'enterprise')),
    max_users       INTEGER      NOT NULL DEFAULT 10,
    max_assets      INTEGER      NOT NULL DEFAULT 100,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    settings        JSONB        NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TRIGGER set_tenants_updated_at
    BEFORE UPDATE ON tenants
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- SITES  (physical branches / locations within a tenant)
-- ============================================================
CREATE TABLE sites (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    name            VARCHAR(255) NOT NULL,
    slug            VARCHAR(100) NOT NULL,
    code            VARCHAR(50),
    description     TEXT,
    street          VARCHAR(255),
    city            VARCHAR(100),
    state           VARCHAR(100),
    country         CHAR(2),    -- ISO 3166-1 alpha-2
    postal          VARCHAR(20),
    latitude        DECIMAL(10,7),
    longitude       DECIMAL(10,7),
    time_zone       VARCHAR(64)  NOT NULL DEFAULT 'UTC',
    contact_name    VARCHAR(255),
    contact_email   CITEXT,
    contact_phone   VARCHAR(50),
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    UNIQUE (tenant_id, slug)
);

CREATE TRIGGER set_sites_updated_at
    BEFORE UPDATE ON sites
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- DEPARTMENTS
-- ============================================================
CREATE TABLE departments (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL REFERENCES tenants(id),
    site_id         UUID        REFERENCES sites(id),
    parent_id       UUID        REFERENCES departments(id),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    manager_user_id UUID,       -- FK added after users table
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE TRIGGER set_departments_updated_at
    BEFORE UPDATE ON departments
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- USERS
-- ============================================================
CREATE TABLE users (
    id                      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id               UUID        NOT NULL REFERENCES tenants(id),
    email                   CITEXT      NOT NULL,
    username                VARCHAR(100) NOT NULL,
    password_hash           VARCHAR(255) NOT NULL,
    first_name              VARCHAR(100),
    last_name               VARCHAR(100),
    phone                   VARCHAR(50),
    avatar_url              TEXT,
    is_active               BOOLEAN      NOT NULL DEFAULT TRUE,
    is_superuser            BOOLEAN      NOT NULL DEFAULT FALSE,
    email_verified          BOOLEAN      NOT NULL DEFAULT FALSE,
    last_login_at           TIMESTAMPTZ,
    last_login_ip           INET,
    failed_login_attempts   INTEGER      NOT NULL DEFAULT 0,
    locked_until            TIMESTAMPTZ,
    mfa_enabled             BOOLEAN      NOT NULL DEFAULT FALSE,
    mfa_secret              VARCHAR(255),
    preferences             JSONB        NOT NULL DEFAULT '{}',
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMPTZ,
    UNIQUE (tenant_id, email),
    UNIQUE (tenant_id, username)
);

CREATE TRIGGER set_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- Back-fill FK from departments to users
ALTER TABLE departments
    ADD CONSTRAINT fk_departments_manager
    FOREIGN KEY (manager_user_id) REFERENCES users(id);

-- ============================================================
-- ROLES
-- ============================================================
CREATE TABLE roles (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   UUID        NOT NULL REFERENCES tenants(id),
    name        VARCHAR(100) NOT NULL,
    slug        VARCHAR(100) NOT NULL,
    description TEXT,
    is_system   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, slug)
);

CREATE TRIGGER set_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION trigger_set_updated_at();

-- ============================================================
-- PERMISSIONS  (system-wide, not per-tenant)
-- ============================================================
CREATE TABLE permissions (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    resource    VARCHAR(100) NOT NULL,  -- e.g. 'dcim.racks', 'asset.servers'
    action      VARCHAR(50)  NOT NULL,  -- create | read | update | delete | list | export
    description TEXT,
    UNIQUE (resource, action)
);

-- ============================================================
-- ROLE ↔ PERMISSION  (M:N)
-- ============================================================
CREATE TABLE role_permissions (
    role_id       UUID NOT NULL REFERENCES roles(id)       ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- ============================================================
-- USER ROLES  (user gets a role, optionally scoped to a site)
-- ============================================================
CREATE TABLE user_roles (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID        NOT NULL REFERENCES users(id)       ON DELETE CASCADE,
    role_id       UUID        NOT NULL REFERENCES roles(id)       ON DELETE CASCADE,
    site_id       UUID        REFERENCES sites(id),
    department_id UUID        REFERENCES departments(id),
    granted_by    UUID        REFERENCES users(id),
    granted_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at    TIMESTAMPTZ,
    UNIQUE NULLS NOT DISTINCT (user_id, role_id, site_id)
);

-- ============================================================
-- REFRESH TOKENS
-- ============================================================
CREATE TABLE refresh_tokens (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    ip_address  INET,
    user_agent  TEXT,
    revoked_at  TIMESTAMPTZ
);

-- ============================================================
-- AUDIT LOGS
-- ============================================================
CREATE TABLE audit_logs (
    id            UUID        NOT NULL DEFAULT gen_random_uuid(),
    tenant_id     UUID        NOT NULL REFERENCES tenants(id),
    user_id       UUID        REFERENCES users(id),
    resource_type VARCHAR(100) NOT NULL,
    resource_id   UUID,
    action        VARCHAR(50)  NOT NULL,
    old_values    JSONB,
    new_values    JSONB,
    ip_address    INET,
    user_agent    TEXT,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Seed system-wide permissions
INSERT INTO permissions (resource, action, description) VALUES
  -- Tenant management
  ('iam.tenants',     'read',   'View tenant details'),
  ('iam.tenants',     'update', 'Modify tenant settings'),
  -- User management
  ('iam.users',       'create', 'Invite or create users'),
  ('iam.users',       'read',   'View user details'),
  ('iam.users',       'update', 'Modify user accounts'),
  ('iam.users',       'delete', 'Remove user accounts'),
  ('iam.users',       'list',   'List users'),
  -- Role management
  ('iam.roles',       'create', 'Create roles'),
  ('iam.roles',       'read',   'View role details'),
  ('iam.roles',       'update', 'Modify roles and permissions'),
  ('iam.roles',       'delete', 'Delete roles'),
  ('iam.roles',       'list',   'List roles'),
  -- Sites
  ('iam.sites',       'create', 'Create sites'),
  ('iam.sites',       'read',   'View site details'),
  ('iam.sites',       'update', 'Modify sites'),
  ('iam.sites',       'delete', 'Delete sites'),
  ('iam.sites',       'list',   'List sites'),
  -- DCIM
  ('dcim.datacenters','create', 'Create data centers'),
  ('dcim.datacenters','read',   'View data centers'),
  ('dcim.datacenters','update', 'Modify data centers'),
  ('dcim.datacenters','delete', 'Delete data centers'),
  ('dcim.datacenters','list',   'List data centers'),
  ('dcim.racks',      'create', 'Create racks'),
  ('dcim.racks',      'read',   'View racks'),
  ('dcim.racks',      'update', 'Modify racks'),
  ('dcim.racks',      'delete', 'Delete racks'),
  ('dcim.racks',      'list',   'List racks'),
  -- Assets
  ('asset.assets',    'create', 'Create assets'),
  ('asset.assets',    'read',   'View asset details'),
  ('asset.assets',    'update', 'Modify assets'),
  ('asset.assets',    'delete', 'Delete/decommission assets'),
  ('asset.assets',    'list',   'List assets'),
  ('asset.assets',    'export', 'Export asset data'),
  -- IPAM
  ('ipam.prefixes',   'create', 'Create IP prefixes'),
  ('ipam.prefixes',   'read',   'View IP prefixes'),
  ('ipam.prefixes',   'update', 'Modify IP prefixes'),
  ('ipam.prefixes',   'delete', 'Delete IP prefixes'),
  ('ipam.prefixes',   'list',   'List IP prefixes'),
  ('ipam.addresses',  'create', 'Allocate IP addresses'),
  ('ipam.addresses',  'read',   'View IP addresses'),
  ('ipam.addresses',  'update', 'Modify IP addresses'),
  ('ipam.addresses',  'delete', 'Release IP addresses'),
  ('ipam.addresses',  'list',   'List IP addresses'),
  -- Monitoring
  ('monitoring.hosts','create', 'Add monitored hosts'),
  ('monitoring.hosts','read',   'View monitoring data'),
  ('monitoring.hosts','update', 'Modify monitoring config'),
  ('monitoring.hosts','delete', 'Remove monitored hosts'),
  ('monitoring.hosts','list',   'List monitored hosts'),
  -- Alerting
  ('alerting.rules',  'create', 'Create alert rules'),
  ('alerting.rules',  'read',   'View alert rules'),
  ('alerting.rules',  'update', 'Modify alert rules'),
  ('alerting.rules',  'delete', 'Delete alert rules'),
  ('alerting.rules',  'list',   'List alert rules'),
  -- Agents
  ('agent.agents',    'create', 'Register agents'),
  ('agent.agents',    'read',   'View agent details'),
  ('agent.agents',    'update', 'Modify agent config'),
  ('agent.agents',    'delete', 'Deregister agents'),
  ('agent.agents',    'list',   'List agents'),
  -- Licenses
  ('license.licenses','create', 'Add software licenses'),
  ('license.licenses','read',   'View license details'),
  ('license.licenses','update', 'Modify licenses'),
  ('license.licenses','delete', 'Delete licenses'),
  ('license.licenses','list',   'List licenses'),
  -- Contracts
  ('contract.contracts','create','Add contracts'),
  ('contract.contracts','read',  'View contracts'),
  ('contract.contracts','update','Modify contracts'),
  ('contract.contracts','delete','Delete contracts'),
  ('contract.contracts','list',  'List contracts'),
  -- Reports
  ('reporting.reports','read',  'View and generate reports'),
  ('reporting.reports','export','Export reports');
