-- 000011_indexes.up.sql
-- Performance indexes for all tables

-- ============================================================
-- IDENTITY
-- ============================================================
CREATE INDEX idx_tenants_slug         ON tenants (slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_sites_tenant         ON sites (tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_sites_slug           ON sites (tenant_id, slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_departments_tenant   ON departments (tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_departments_parent   ON departments (parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX idx_users_tenant         ON users (tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_email          ON users (tenant_id, email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_username       ON users (tenant_id, username) WHERE deleted_at IS NULL;
CREATE INDEX idx_roles_tenant         ON roles (tenant_id);
CREATE INDEX idx_user_roles_user      ON user_roles (user_id);
CREATE INDEX idx_user_roles_role      ON user_roles (role_id);
CREATE INDEX idx_refresh_tokens_user  ON refresh_tokens (user_id) WHERE revoked_at IS NULL;
CREATE INDEX idx_audit_logs_tenant    ON audit_logs (tenant_id, created_at DESC);
CREATE INDEX idx_audit_logs_resource  ON audit_logs (resource_type, resource_id);
CREATE INDEX idx_audit_logs_user      ON audit_logs (user_id, created_at DESC);

-- ============================================================
-- DCIM
-- ============================================================
CREATE INDEX idx_data_centers_tenant  ON data_centers (tenant_id, site_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_rooms_dc             ON rooms (data_center_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_racks_tenant         ON racks (tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_racks_site           ON racks (site_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_racks_room           ON racks (room_id) WHERE room_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_power_feeds_rack     ON power_feeds (rack_id);
CREATE INDEX idx_pdus_rack            ON pdus (rack_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_patch_panels_rack    ON patch_panels (rack_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_cables_type          ON cables (tenant_id, cable_type);
CREATE INDEX idx_cable_terminations   ON cable_terminations (termination_type, termination_id);

-- ============================================================
-- ASSETS
-- ============================================================
CREATE INDEX idx_assets_tenant        ON assets (tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_assets_site          ON assets (site_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_assets_category      ON assets (tenant_id, category) WHERE deleted_at IS NULL;
CREATE INDEX idx_assets_status        ON assets (tenant_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_assets_rack          ON assets (rack_id) WHERE rack_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_assets_serial        ON assets (tenant_id, serial_number) WHERE serial_number IS NOT NULL;
CREATE INDEX idx_assets_tag           ON assets (tenant_id, asset_tag) WHERE asset_tag IS NOT NULL;
CREATE INDEX idx_assets_ip            ON assets USING gist (primary_ip inet_ops) WHERE primary_ip IS NOT NULL;
CREATE INDEX idx_assets_name_trgm     ON assets USING gin (name gin_trgm_ops);
CREATE INDEX idx_ni_asset             ON network_interfaces (asset_id);
CREATE INDEX idx_ni_mac               ON network_interfaces (mac_address) WHERE mac_address IS NOT NULL;
CREATE INDEX idx_vm_host              ON virtual_machines (host_asset_id) WHERE host_asset_id IS NOT NULL;

-- ============================================================
-- IPAM
-- ============================================================
CREATE INDEX idx_vrfs_tenant          ON vrfs (tenant_id);
CREATE INDEX idx_vlans_tenant         ON vlans (tenant_id, vid);
CREATE INDEX idx_vlans_site           ON vlans (site_id) WHERE site_id IS NOT NULL;
CREATE INDEX idx_prefixes_tenant      ON prefixes (tenant_id, vrf_id);
CREATE INDEX idx_prefixes_prefix      ON prefixes USING gist (prefix inet_ops);
CREATE INDEX idx_ip_addresses_tenant  ON ip_addresses (tenant_id, vrf_id);
CREATE INDEX idx_ip_addresses_addr    ON ip_addresses USING gist (address inet_ops);
CREATE INDEX idx_ip_addresses_dns     ON ip_addresses (dns_name) WHERE dns_name IS NOT NULL;
CREATE INDEX idx_ip_addresses_obj     ON ip_addresses (assigned_object_type, assigned_object_id)
    WHERE assigned_object_id IS NOT NULL;
CREATE INDEX idx_dhcp_leases_ip       ON dhcp_leases USING gist (ip_address inet_ops);
CREATE INDEX idx_dhcp_leases_mac      ON dhcp_leases (mac_address) WHERE mac_address IS NOT NULL;
CREATE INDEX idx_dns_records_zone     ON dns_records (zone_id, record_type);
CREATE INDEX idx_dns_records_name     ON dns_records (zone_id, name);

-- ============================================================
-- MONITORING
-- ============================================================
CREATE INDEX idx_monitored_hosts_tenant  ON monitored_hosts (tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_monitored_hosts_asset   ON monitored_hosts (asset_id) WHERE asset_id IS NOT NULL;
CREATE INDEX idx_monitored_hosts_agent   ON monitored_hosts (agent_id) WHERE agent_id IS NOT NULL;
CREATE INDEX idx_monitored_hosts_status  ON monitored_hosts (tenant_id, last_status);
CREATE INDEX idx_check_defs_host         ON check_definitions (host_id) WHERE is_enabled;
CREATE INDEX idx_metric_data_host        ON metric_data (host_id, metric_name, time DESC);
CREATE INDEX idx_metric_data_tenant      ON metric_data (tenant_id, time DESC);
CREATE INDEX idx_avail_records_host      ON availability_records (host_id, started_at DESC);
CREATE INDEX idx_service_checks_host     ON service_checks (host_id) WHERE is_enabled;

-- ============================================================
-- ALERTING
-- ============================================================
CREATE INDEX idx_alert_events_tenant     ON alert_events (tenant_id, fired_at DESC);
CREATE INDEX idx_alert_events_status     ON alert_events (tenant_id, status) WHERE status = 'firing';
CREATE INDEX idx_alert_events_host       ON alert_events (host_id, fired_at DESC)
    WHERE host_id IS NOT NULL;
CREATE INDEX idx_alert_events_rule       ON alert_events (rule_id, fired_at DESC);
CREATE INDEX idx_alert_notifs_event      ON alert_notifications (event_id);
CREATE INDEX idx_alert_notifs_pending    ON alert_notifications (status, created_at)
    WHERE status = 'pending';
CREATE INDEX idx_silences_tenant         ON silences (tenant_id, ends_at)
    WHERE ends_at IS NOT NULL;

-- ============================================================
-- AGENTS
-- ============================================================
CREATE INDEX idx_agents_tenant           ON agents (tenant_id);
CREATE INDEX idx_agents_status           ON agents (tenant_id, status);
CREATE INDEX idx_agents_api_key          ON agents (api_key_hash);
CREATE INDEX idx_agent_health_agent      ON agent_health_metrics (agent_id, time DESC);
CREATE INDEX idx_agent_tasks_agent       ON agent_tasks (agent_id, status)
    WHERE status IN ('pending','delivered');

-- ============================================================
-- LICENSES & CONTRACTS
-- ============================================================
CREATE INDEX idx_licenses_tenant         ON software_licenses (tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_licenses_expiry         ON software_licenses (expiration_date)
    WHERE expiration_date IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_license_assignments     ON license_assignments (license_id);
CREATE INDEX idx_contracts_tenant        ON contracts (tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_contracts_end_date      ON contracts (end_date)
    WHERE end_date IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_warranties_asset        ON warranties (asset_id);
CREATE INDEX idx_warranties_end_date     ON warranties (end_date);

-- ============================================================
-- DISCOVERY
-- ============================================================
CREATE INDEX idx_discovery_jobs_tenant   ON discovery_jobs (tenant_id, status);
CREATE INDEX idx_discovery_results_job   ON discovery_results (job_id, status);
CREATE INDEX idx_discovery_results_ip    ON discovery_results USING gist (ip_address inet_ops)
    WHERE ip_address IS NOT NULL;
