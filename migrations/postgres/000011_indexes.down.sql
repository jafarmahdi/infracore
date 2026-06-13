-- 000011_indexes.down.sql
-- Drop all custom indexes (PostgreSQL drops indexes with tables automatically,
-- but this handles any future standalone index removals)

-- Discovery
DROP INDEX IF EXISTS idx_discovery_results_ip;
DROP INDEX IF EXISTS idx_discovery_results_job;
DROP INDEX IF EXISTS idx_discovery_jobs_tenant;

-- Licenses & Contracts
DROP INDEX IF EXISTS idx_warranties_end_date;
DROP INDEX IF EXISTS idx_warranties_asset;
DROP INDEX IF EXISTS idx_contracts_end_date;
DROP INDEX IF EXISTS idx_contracts_tenant;
DROP INDEX IF EXISTS idx_license_assignments;
DROP INDEX IF EXISTS idx_licenses_expiry;
DROP INDEX IF EXISTS idx_licenses_tenant;

-- Agents
DROP INDEX IF EXISTS idx_agent_tasks_agent;
DROP INDEX IF EXISTS idx_agent_health_agent;
DROP INDEX IF EXISTS idx_agents_api_key;
DROP INDEX IF EXISTS idx_agents_status;
DROP INDEX IF EXISTS idx_agents_tenant;

-- Alerting
DROP INDEX IF EXISTS idx_silences_tenant;
DROP INDEX IF EXISTS idx_alert_notifs_pending;
DROP INDEX IF EXISTS idx_alert_notifs_event;
DROP INDEX IF EXISTS idx_alert_events_rule;
DROP INDEX IF EXISTS idx_alert_events_host;
DROP INDEX IF EXISTS idx_alert_events_status;
DROP INDEX IF EXISTS idx_alert_events_tenant;

-- Monitoring
DROP INDEX IF EXISTS idx_service_checks_host;
DROP INDEX IF EXISTS idx_avail_records_host;
DROP INDEX IF EXISTS idx_metric_data_tenant;
DROP INDEX IF EXISTS idx_metric_data_host;
DROP INDEX IF EXISTS idx_check_defs_host;
DROP INDEX IF EXISTS idx_monitored_hosts_status;
DROP INDEX IF EXISTS idx_monitored_hosts_agent;
DROP INDEX IF EXISTS idx_monitored_hosts_asset;
DROP INDEX IF EXISTS idx_monitored_hosts_tenant;

-- IPAM
DROP INDEX IF EXISTS idx_dns_records_name;
DROP INDEX IF EXISTS idx_dns_records_zone;
DROP INDEX IF EXISTS idx_dhcp_leases_mac;
DROP INDEX IF EXISTS idx_dhcp_leases_ip;
DROP INDEX IF EXISTS idx_ip_addresses_obj;
DROP INDEX IF EXISTS idx_ip_addresses_dns;
DROP INDEX IF EXISTS idx_ip_addresses_addr;
DROP INDEX IF EXISTS idx_ip_addresses_tenant;
DROP INDEX IF EXISTS idx_prefixes_prefix;
DROP INDEX IF EXISTS idx_prefixes_tenant;
DROP INDEX IF EXISTS idx_vlans_site;
DROP INDEX IF EXISTS idx_vlans_tenant;
DROP INDEX IF EXISTS idx_vrfs_tenant;

-- Assets
DROP INDEX IF EXISTS idx_vm_host;
DROP INDEX IF EXISTS idx_ni_mac;
DROP INDEX IF EXISTS idx_ni_asset;
DROP INDEX IF EXISTS idx_assets_name_trgm;
DROP INDEX IF EXISTS idx_assets_ip;
DROP INDEX IF EXISTS idx_assets_tag;
DROP INDEX IF EXISTS idx_assets_serial;
DROP INDEX IF EXISTS idx_assets_rack;
DROP INDEX IF EXISTS idx_assets_status;
DROP INDEX IF EXISTS idx_assets_category;
DROP INDEX IF EXISTS idx_assets_site;
DROP INDEX IF EXISTS idx_assets_tenant;

-- DCIM
DROP INDEX IF EXISTS idx_cable_terminations;
DROP INDEX IF EXISTS idx_cables_type;
DROP INDEX IF EXISTS idx_patch_panels_rack;
DROP INDEX IF EXISTS idx_pdus_rack;
DROP INDEX IF EXISTS idx_power_feeds_rack;
DROP INDEX IF EXISTS idx_racks_room;
DROP INDEX IF EXISTS idx_racks_site;
DROP INDEX IF EXISTS idx_racks_tenant;
DROP INDEX IF EXISTS idx_rooms_dc;
DROP INDEX IF EXISTS idx_data_centers_tenant;

-- Identity
DROP INDEX IF EXISTS idx_audit_logs_user;
DROP INDEX IF EXISTS idx_audit_logs_resource;
DROP INDEX IF EXISTS idx_audit_logs_tenant;
DROP INDEX IF EXISTS idx_refresh_tokens_user;
DROP INDEX IF EXISTS idx_user_roles_role;
DROP INDEX IF EXISTS idx_user_roles_user;
DROP INDEX IF EXISTS idx_roles_tenant;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_tenant;
DROP INDEX IF EXISTS idx_departments_parent;
DROP INDEX IF EXISTS idx_departments_tenant;
DROP INDEX IF EXISTS idx_sites_slug;
DROP INDEX IF EXISTS idx_sites_tenant;
DROP INDEX IF EXISTS idx_tenants_slug;
