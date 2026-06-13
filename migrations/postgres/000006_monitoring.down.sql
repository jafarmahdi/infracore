-- 000006_monitoring.down.sql
DROP TABLE IF EXISTS service_checks;
DROP TABLE IF EXISTS availability_records;
DROP TABLE IF EXISTS metric_data;
DROP TABLE IF EXISTS check_definitions;
DROP TABLE IF EXISTS monitored_hosts;
DROP TABLE IF EXISTS ssh_keys;
DROP TABLE IF EXISTS monitoring_profiles;
