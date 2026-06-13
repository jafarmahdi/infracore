export type TrendDirection = "up" | "down" | "neutral";
export type Severity = "critical" | "high" | "warning" | "info";

export interface DashboardSummary {
  totalAssets: number;
  onlineDevices: number;
  offlineDevices: number;
  activeAlerts: number;
  expiringLicenses: number;
  sla: number;
  trends: Record<string, { value: number; direction: TrendDirection }>;
}

export interface UtilizationPoint {
  time: string;
  cpu: number;
  memory: number;
  disk: number;
}

export interface TrafficPoint {
  time: string;
  inbound: number;
  outbound: number;
}

export interface RecentEvent {
  id: string;
  title: string;
  source: string;
  severity: Severity;
  timestamp: string;
}

export interface DashboardData {
  summary: DashboardSummary;
  utilization: UtilizationPoint[];
  traffic: TrafficPoint[];
  events: RecentEvent[];
  deviceHealth: { name: string; value: number; color: string }[];
}
