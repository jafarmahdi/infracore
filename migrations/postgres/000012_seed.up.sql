-- 000012_seed.up.sql
-- Initial seed: audit_log partitions, default tenant, superuser admin.
-- Uses pgcrypto crypt() which produces a $2a$ bcrypt hash compatible with
-- Go's golang.org/x/crypto/bcrypt.

-- ── Audit log partitions (required before any INSERT can land) ────────────
CREATE TABLE audit_logs_2026
    PARTITION OF audit_logs
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');

CREATE TABLE audit_logs_2027
    PARTITION OF audit_logs
    FOR VALUES FROM ('2027-01-01') TO ('2028-01-01');

CREATE TABLE audit_logs_2028
    PARTITION OF audit_logs
    FOR VALUES FROM ('2028-01-01') TO ('2029-01-01');

-- ── Default tenant ─────────────────────────────────────────────────────────
INSERT INTO tenants (id, name, slug, plan, max_users, max_assets)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Demo Organization',
    'demo',
    'enterprise',
    1000,
    100000
);

-- ── Superuser admin ────────────────────────────────────────────────────────
-- Default password: Admin@123456
-- Hash generated via pgcrypto bcrypt (compatible with Go's bcrypt library).
INSERT INTO users (
    id, tenant_id, email, username, password_hash,
    first_name, last_name, is_active, is_superuser, email_verified
)
VALUES (
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000001',
    'admin@demo.com',
    'admin',
    crypt('Admin@123456', gen_salt('bf', 10)),
    'Admin',
    'User',
    TRUE,
    TRUE,
    TRUE
);
