-- 000008_agents.down.sql
DROP TABLE IF EXISTS agent_tasks;
DROP TABLE IF EXISTS agent_update_channels;
DROP TABLE IF EXISTS agent_health_metrics;
ALTER TABLE monitored_hosts DROP CONSTRAINT IF EXISTS fk_monitored_hosts_agent;
DROP TABLE IF EXISTS agents;
DROP TABLE IF EXISTS agent_groups;
