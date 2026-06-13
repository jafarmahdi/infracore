export type MonitoringType = "agent" | "snmp" | "wmi" | "icmp" | "ssh" | "api";
export type HostStatus = "up" | "down" | "warning" | "maintenance" | "unknown";

export interface MonitoredHost {
  id: string;
  tenant_id: string;
  asset_id?: string | null;
  name: string;
  display_name?: string | null;
  ip_address?: string | null;
  monitoring_type: MonitoringType;
  port?: number | null;
  check_interval_seconds: number;
  timeout_seconds: number;
  is_active: boolean;
  status: HostStatus;
  last_status?: string | null;
  last_check_at?: string | null;
  snmp_version?: string | null;
  snmp_community?: string | null;
  agent_id?: string | null;
  tags?: string[] | null;
  notes?: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateHostRequest {
  name: string;
  monitoring_type: MonitoringType;
  ip_address?: string;
  port?: number;
  display_name?: string;
  check_interval_seconds?: number;
  timeout_seconds?: number;
  snmp_version?: string;
  snmp_community?: string;
  notes?: string;
}

export interface ListHostsResponse {
  items: MonitoredHost[];
  total_items: number;
  total_pages: number;
  page: number;
  page_size: number;
}

export interface HostStatusCounts {
  up: number;
  down: number;
  warning: number;
  maintenance: number;
  unknown: number;
  total: number;
}
