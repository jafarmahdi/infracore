import { createBrowserRouter, Navigate } from "react-router-dom";
import { ProtectedRoute } from "@/components/auth/protected-route";
import { RoleGuard } from "@/components/auth/role-guard";
import { AppLayout } from "@/components/layout/app-layout";
import { LoginPage } from "@/pages/auth/login-page";
import { DashboardPage } from "@/pages/dashboard/dashboard-page";
import { DCIMPage } from "@/pages/dcim/dcim-page";
import { AssetsPage } from "@/pages/assets/assets-page";
import { IPAMPage } from "@/pages/ipam/ipam-page";
import { MonitoringPage } from "@/pages/monitoring/monitoring-page";
import { AlertsPage } from "@/pages/alerts/alerts-page";
import { AgentsPage } from "@/pages/agents/agents-page";
import { LicensesPage } from "@/pages/licenses/licenses-page";
import { ContractsPage } from "@/pages/contracts/contracts-page";
import { ReportsPage } from "@/pages/reports/reports-page";
import { CablingPage } from "@/pages/cabling/cabling-page";
import { SecurityPage } from "@/pages/security/security-page";
import { SettingsPage } from "@/pages/settings/settings-page";

export const router = createBrowserRouter([
  { path: "/login", element: <LoginPage /> },
  {
    element: <ProtectedRoute />,
    children: [{
      element: <AppLayout />,
      children: [
        { index: true, element: <Navigate to="/dashboard" replace /> },
        { path: "/dashboard",  element: <DashboardPage /> },
        { path: "/dcim",       element: <DCIMPage /> },
        { path: "/assets",     element: <AssetsPage /> },
        { path: "/ipam",       element: <IPAMPage /> },
        { path: "/monitoring", element: <MonitoringPage /> },
        { path: "/alerts",     element: <AlertsPage /> },
        { path: "/agents",     element: <AgentsPage /> },
        { path: "/licenses",   element: <LicensesPage /> },
        { path: "/contracts",  element: <ContractsPage /> },
        { path: "/reports",    element: <ReportsPage /> },
        { path: "/cabling",    element: <CablingPage /> },
        {
          element: <RoleGuard allowed={["admin"]} />,
          children: [
            { path: "/security", element: <SecurityPage /> },
            { path: "/settings", element: <SettingsPage /> },
          ],
        },
      ],
    }],
  },
  { path: "*", element: <Navigate to="/dashboard" replace /> },
]);
