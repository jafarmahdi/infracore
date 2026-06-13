-- 000001_init_extensions.down.sql
DROP FUNCTION IF EXISTS trigger_set_updated_at();
DROP EXTENSION IF EXISTS "timescaledb";
DROP EXTENSION IF EXISTS "btree_gist";
DROP EXTENSION IF EXISTS "pg_trgm";
DROP EXTENSION IF EXISTS "citext";
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP EXTENSION IF EXISTS "pgcrypto";
