import { Navigate, Outlet } from "react-router-dom";
import { useAuthStore } from "@/stores/auth-store";
import type { UserRole } from "@/types/auth";

export function RoleGuard({ allowed }: { allowed: UserRole[] }) {
  const role = useAuthStore((state) => state.user?.role);
  return role && allowed.includes(role) ? <Outlet /> : <Navigate to="/dashboard" replace />;
}
