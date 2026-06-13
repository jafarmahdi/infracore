import {
  Activity, BellRing, Boxes, Building2, Cable, FileBarChart, FileCheck2,
  Gauge, KeyRound, Network, RadioTower, Settings, ShieldCheck,
} from "lucide-react";
import type { LucideIcon } from "lucide-react";
import type { UserRole } from "@/types/auth";

export interface NavigationItem {
  label: string;
  path: string;
  icon: LucideIcon;
  roles?: UserRole[];
  badge?: string;
}

export const navigation: NavigationItem[] = [
  { label: "Overview", path: "/dashboard", icon: Gauge },
  { label: "DCIM", path: "/dcim", icon: Building2 },
  { label: "Assets", path: "/assets", icon: Boxes },
  { label: "IPAM", path: "/ipam", icon: Network },
  { label: "Monitoring", path: "/monitoring", icon: Activity },
  { label: "Alerts", path: "/alerts", icon: BellRing, badge: "17" },
  { label: "Agents", path: "/agents", icon: RadioTower },
  { label: "Licenses", path: "/licenses", icon: KeyRound },
  { label: "Contracts", path: "/contracts", icon: FileCheck2 },
  { label: "Reports", path: "/reports", icon: FileBarChart },
  { label: "Cabling", path: "/cabling", icon: Cable },
  { label: "Security", path: "/security", icon: ShieldCheck, roles: ["admin"] },
  { label: "Settings", path: "/settings", icon: Settings, roles: ["admin"] },
];
