-- 000001_init_extensions.up.sql
-- Enable required PostgreSQL extensions

CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";       -- trigram indexes for fast text search
CREATE EXTENSION IF NOT EXISTS "btree_gist";    -- exclusion constraints on ranges
CREATE EXTENSION IF NOT EXISTS "citext";        -- case-insensitive text type
CREATE EXTENSION IF NOT EXISTS "timescaledb" CASCADE; -- time-series engine

-- Utility function: automatically update updated_at on row changes
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
