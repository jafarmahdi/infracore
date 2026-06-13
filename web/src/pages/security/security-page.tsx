import { useState } from "react";
import { AlertTriangle, CheckCircle2, Globe, Lock, Search, Shield, ShieldCheck, XCircle } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { MetricCard } from "@/components/dashboard/metric-card";

const SESSIONS = [
  { id: "SES-001", user: "Omar Hassan",      email: "admin@infracore.io", ip: "10.10.0.100", device: "Chrome / Windows 11", location: "Riyadh, SA",  loggedIn: "2026-06-13T09:00:00Z", current: true },
  { id: "SES-002", user: "Fatima Al-Rashid", email: "f.rashid@ops.io",   ip: "10.10.0.102", device: "Firefox / Ubuntu",    location: "Riyadh, SA",  loggedIn: "2026-06-13T08:30:00Z", current: false },
  { id: "SES-003", user: "Khaled Nasser",    email: "k.nasser@ops.io",   ip: "10.20.0.55",  device: "Safari / macOS",      location: "Jeddah, SA",  loggedIn: "2026-06-13T07:15:00Z", current: false },
  { id: "SES-004", user: "Nora Ibrahim",     email: "n.ibrahim@ops.io",  ip: "10.10.0.118", device: "Edge / Windows 10",   location: "Riyadh, SA",  loggedIn: "2026-06-13T11:45:00Z", current: false },
];

const AUDIT_LOG = [
  { id: "AUD-001", user: "Omar Hassan",      action: "LOGIN",          resource: "auth.session",      result: "success", ip: "10.10.0.100", time: "13 Jun 09:00" },
  { id: "AUD-002", user: "Omar Hassan",      action: "CREATE",         resource: "iam.users",         result: "success", ip: "10.10.0.100", time: "13 Jun 09:05" },
  { id: "AUD-003", user: "Fatima Al-Rashid", action: "UPDATE",         resource: "monitoring.hosts",  result: "success", ip: "10.10.0.102", time: "13 Jun 09:12" },
  { id: "AUD-004", user: "unknown",          action: "LOGIN",          resource: "auth.session",      result: "failed",  ip: "185.123.45.6",time: "13 Jun 10:01" },
  { id: "AUD-005", user: "unknown",          action: "LOGIN",          resource: "auth.session",      result: "failed",  ip: "185.123.45.6",time: "13 Jun 10:02" },
  { id: "AUD-006", user: "Khaled Nasser",    action: "DELETE",         resource: "dcim.racks",        result: "denied",  ip: "10.20.0.55",  time: "13 Jun 11:30" },
  { id: "AUD-007", user: "Omar Hassan",      action: "EXPORT",         resource: "asset.assets",      result: "success", ip: "10.10.0.100", time: "13 Jun 12:00" },
  { id: "AUD-008", user: "Nora Ibrahim",     action: "UPDATE",         resource: "iam.roles",         result: "success", ip: "10.10.0.118", time: "13 Jun 12:10" },
];

const RESULT_VARIANT: Record<string, "success" | "destructive" | "warning"> = {
  success: "success", failed: "destructive", denied: "warning",
};

const SECURITY_CHECKS = [
  { label: "MFA enforcement",          status: "pass" },
  { label: "Password policy (12+ chars)", status: "pass" },
  { label: "Account lockout (5 attempts)",status: "pass" },
  { label: "Session timeout (4h)",     status: "pass" },
  { label: "Audit logging",            status: "pass" },
  { label: "Failed login alerts",      status: "warn" },
  { label: "API rate limiting",        status: "pass" },
  { label: "HTTPS/TLS enforced",       status: "pass" },
];

export function SecurityPage() {
  const [search, setSearch] = useState("");
  const passChecks = SECURITY_CHECKS.filter(c => c.status === "pass").length;
  const failedLogins = AUDIT_LOG.filter(a => a.result === "failed").length;

  const filteredAudit = AUDIT_LOG.filter(a =>
    a.user.toLowerCase().includes(search.toLowerCase()) ||
    a.resource.toLowerCase().includes(search.toLowerCase()) ||
    a.action.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="space-y-5">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <MetricCard title="Active sessions" value={SESSIONS.length} description="logged in now" icon={Globe} iconClassName="bg-blue-500/10 text-blue-500" />
        <MetricCard title="Security score" value={passChecks} suffix={`/${SECURITY_CHECKS.length}`} description="checks passing" icon={ShieldCheck} iconClassName="bg-emerald-500/10 text-emerald-500" />
        <MetricCard title="Failed logins" value={failedLogins} description="last 24 hours" icon={XCircle} iconClassName="bg-red-500/10 text-red-500" trend={{ value: 2, direction: "up" }} />
        <MetricCard title="Audit events" value={AUDIT_LOG.length} description="today" icon={Shield} iconClassName="bg-violet-500/10 text-violet-500" />
      </div>

      <div className="grid gap-4 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle className="text-sm font-semibold">Active sessions</CardTitle>
          </CardHeader>
          <CardContent className="pt-0 overflow-x-auto">
            <table className="w-full text-sm">
              <thead><tr className="border-b text-left">
                {["User", "IP Address", "Device", "Location", "Logged in", ""].map(h => (
                  <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground whitespace-nowrap">{h}</th>
                ))}
              </tr></thead>
              <tbody className="divide-y">
                {SESSIONS.map(s => (
                  <tr key={s.id} className="hover:bg-muted/30">
                    <td className="py-3 pr-4">
                      <div className="font-medium">{s.user}</div>
                      <div className="text-xs text-muted-foreground">{s.email}</div>
                    </td>
                    <td className="py-3 pr-4 font-mono text-xs">{s.ip}</td>
                    <td className="py-3 pr-4 text-muted-foreground text-xs">{s.device}</td>
                    <td className="py-3 pr-4 text-muted-foreground">{s.location}</td>
                    <td className="py-3 pr-4 font-mono text-xs">{new Date(s.loggedIn).toLocaleTimeString()}</td>
                    <td className="py-3">
                      {s.current
                        ? <Badge variant="success">current</Badge>
                        : <Button variant="ghost" size="sm" className="text-xs text-destructive hover:text-destructive"><Lock className="h-3 w-3" /> Revoke</Button>
                      }
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-semibold">Security posture</CardTitle>
          </CardHeader>
          <CardContent className="pt-0 space-y-2">
            {SECURITY_CHECKS.map(c => (
              <div key={c.label} className="flex items-center justify-between py-1">
                <span className="text-sm">{c.label}</span>
                {c.status === "pass"
                  ? <CheckCircle2 className="h-4 w-4 text-emerald-500 shrink-0" />
                  : <AlertTriangle className="h-4 w-4 text-amber-500 shrink-0" />
                }
              </div>
            ))}
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="text-sm font-semibold">Audit log</CardTitle>
            <div className="relative">
              <Search className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
              <Input placeholder="User, action, resource…" className="pl-8 h-8 w-52 text-sm" value={search} onChange={e => setSearch(e.target.value)} />
            </div>
          </div>
        </CardHeader>
        <CardContent className="pt-0 overflow-x-auto">
          <table className="w-full text-sm">
            <thead><tr className="border-b text-left">
              {["ID", "User", "Action", "Resource", "Result", "IP", "Time"].map(h => (
                <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground whitespace-nowrap">{h}</th>
              ))}
            </tr></thead>
            <tbody className="divide-y">
              {filteredAudit.map(a => (
                <tr key={a.id} className="hover:bg-muted/30">
                  <td className="py-3 pr-4 font-mono text-xs">{a.id}</td>
                  <td className="py-3 pr-4">{a.user}</td>
                  <td className="py-3 pr-4 font-mono text-xs font-medium">{a.action}</td>
                  <td className="py-3 pr-4 text-muted-foreground">{a.resource}</td>
                  <td className="py-3 pr-4"><Badge variant={RESULT_VARIANT[a.result]}>{a.result}</Badge></td>
                  <td className="py-3 pr-4 font-mono text-xs">{a.ip}</td>
                  <td className="py-3 text-muted-foreground text-xs">{a.time}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
