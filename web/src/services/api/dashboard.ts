import type { DashboardData } from "@/types/dashboard";
import { apiClient } from "./client";

const useMocks = import.meta.env.VITE_USE_MOCKS !== "false";

const mockDashboard: DashboardData = {
  summary: {
    totalAssets: 2847,
    onlineDevices: 2614,
    offlineDevices: 38,
    activeAlerts: 17,
    expiringLicenses: 12,
    sla: 99.96,
    trends: {
      totalAssets: { value: 3.8, direction: "up" },
      onlineDevices: { value: 1.2, direction: "up" },
      activeAlerts: { value: 8.4, direction: "down" },
      expiringLicenses: { value: 4, direction: "up" },
    },
  },
  utilization: [
    { time: "00:00", cpu: 42, memory: 61, disk: 47 },
    { time: "04:00", cpu: 38, memory: 58, disk: 47 },
    { time: "08:00", cpu: 64, memory: 69, disk: 48 },
    { time: "12:00", cpu: 72, memory: 74, disk: 49 },
    { time: "16:00", cpu: 58, memory: 71, disk: 49 },
    { time: "20:00", cpu: 48, memory: 65, disk: 50 },
    { time: "Now", cpu: 53, memory: 68, disk: 50 },
  ],
  traffic: [
    { time: "00:00", inbound: 2.1, outbound: 1.2 },
    { time: "04:00", inbound: 1.8, outbound: 0.9 },
    { time: "08:00", inbound: 4.5, outbound: 2.7 },
    { time: "12:00", inbound: 6.2, outbound: 3.8 },
    { time: "16:00", inbound: 5.4, outbound: 3.1 },
    { time: "20:00", inbound: 3.8, outbound: 2.4 },
    { time: "Now", inbound: 4.6, outbound: 2.8 },
  ],
  deviceHealth: [
    { name: "Healthy", value: 2614, color: "#22c55e" },
    { name: "Warning", value: 195, color: "#f59e0b" },
    { name: "Offline", value: 38, color: "#ef4444" },
  ],
  events: [
    { id: "evt-1", title: "Core switch packet loss above threshold", source: "SW-CORE-01", severity: "critical", timestamp: new Date(Date.now() - 120000).toISOString() },
    { id: "evt-2", title: "Power supply redundancy restored", source: "DC1-R14-SRV08", severity: "info", timestamp: new Date(Date.now() - 480000).toISOString() },
    { id: "evt-3", title: "Storage volume reached 85% capacity", source: "SAN-PROD-02", severity: "warning", timestamp: new Date(Date.now() - 960000).toISOString() },
    { id: "evt-4", title: "Agent version out of compliance", source: "BR-BGD-FW01", severity: "high", timestamp: new Date(Date.now() - 1500000).toISOString() },
  ],
};

export async function getDashboard(): Promise<DashboardData> {
  if (useMocks) {
    await new Promise((resolve) => setTimeout(resolve, 550));
    return mockDashboard;
  }
  const { data } = await apiClient.get<DashboardData>("/dashboard");
  return data;
}
