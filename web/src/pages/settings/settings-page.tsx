import { useState } from "react";
import { Bell, Building2, Key, Mail, Moon, Save, Sun, Users } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useUIStore } from "@/stores/ui-store";

const TEAM_MEMBERS = [
  { name: "Omar Hassan",       email: "admin@infracore.io",  role: "admin",    status: "active", lastSeen: "Now" },
  { name: "Fatima Al-Rashid", email: "f.rashid@ops.io",     role: "operator", status: "active", lastSeen: "2h ago" },
  { name: "Khaled Nasser",    email: "k.nasser@ops.io",     role: "operator", status: "active", lastSeen: "5h ago" },
  { name: "Nora Ibrahim",     email: "n.ibrahim@ops.io",    role: "admin",    status: "active", lastSeen: "1h ago" },
  { name: "Ahmed Al-Zahrani", email: "a.zahrani@ops.io",    role: "viewer",   status: "invited",lastSeen: "–" },
];

const ROLE_VARIANT: Record<string, "default" | "secondary" | "warning"> = {
  admin: "default", operator: "secondary", viewer: "secondary",
};

const ST_VARIANT: Record<string, "success" | "secondary"> = {
  active: "success", invited: "secondary",
};

const INTEGRATIONS = [
  { name: "SMTP Email",         description: "Send alert notifications and reports", connected: true,  icon: Mail },
  { name: "Slack",              description: "Push alerts to Slack channels",        connected: false, icon: Bell },
  { name: "PagerDuty",          description: "Escalate critical incidents",          connected: true,  icon: Bell },
  { name: "LDAP / Active Dir.", description: "Sync users from directory service",    connected: false, icon: Users },
  { name: "Jira",               description: "Create tickets from alert events",     connected: false, icon: Key },
];

type SettingsTab = "organization" | "team" | "notifications" | "integrations" | "appearance";

export function SettingsPage() {
  const [tab, setTab] = useState<SettingsTab>("organization");
  const { theme, toggleTheme, locale, setLocale } = useUIStore();

  const tabs: { value: SettingsTab; label: string }[] = [
    { value: "organization", label: "Organization" },
    { value: "team", label: "Team" },
    { value: "notifications", label: "Notifications" },
    { value: "integrations", label: "Integrations" },
    { value: "appearance", label: "Appearance" },
  ];

  return (
    <div className="space-y-5 max-w-4xl">
      <div className="flex gap-1 rounded-lg border bg-muted/40 p-1 w-fit">
        {tabs.map(t => (
          <button key={t.value} onClick={() => setTab(t.value)}
            className={`rounded-md px-4 py-1.5 text-sm font-medium transition-colors ${tab === t.value ? "bg-background shadow-sm" : "text-muted-foreground hover:text-foreground"}`}>
            {t.label}
          </button>
        ))}
      </div>

      {tab === "organization" && (
        <Card>
          <CardHeader><CardTitle className="text-sm">Organization details</CardTitle></CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-1.5">
                <Label>Organization name</Label>
                <Input defaultValue="InfraCore Operations" />
              </div>
              <div className="space-y-1.5">
                <Label>Slug</Label>
                <Input defaultValue="infracore" className="font-mono" />
              </div>
              <div className="space-y-1.5">
                <Label>Contact email</Label>
                <Input type="email" defaultValue="ops@infracore.io" />
              </div>
              <div className="space-y-1.5">
                <Label>Timezone</Label>
                <select className="h-10 w-full rounded-md border bg-background px-3 text-sm text-foreground">
                  <option>Asia/Riyadh (UTC+3)</option>
                  <option>UTC</option>
                </select>
              </div>
            </div>
            <div className="space-y-1.5">
              <Label>Plan</Label>
              <div className="flex items-center gap-3 rounded-lg border bg-muted/30 p-3">
                <Building2 className="h-5 w-5 text-primary" />
                <div>
                  <p className="text-sm font-medium">Enterprise</p>
                  <p className="text-xs text-muted-foreground">1,000 users · 100,000 assets · Unlimited sites</p>
                </div>
                <Badge className="ml-auto">Active</Badge>
              </div>
            </div>
            <div className="flex justify-end">
              <Button><Save className="h-4 w-4" /> Save changes</Button>
            </div>
          </CardContent>
        </Card>
      )}

      {tab === "team" && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="text-sm">Team members</CardTitle>
              <Button size="sm"><Users className="h-3.5 w-3.5" /> Invite member</Button>
            </div>
          </CardHeader>
          <CardContent className="pt-0 overflow-x-auto">
            <table className="w-full text-sm">
              <thead><tr className="border-b text-left">
                {["Member", "Role", "Status", "Last seen", ""].map(h => (
                  <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground">{h}</th>
                ))}
              </tr></thead>
              <tbody className="divide-y">
                {TEAM_MEMBERS.map(m => (
                  <tr key={m.email} className="hover:bg-muted/30">
                    <td className="py-3 pr-4">
                      <div className="font-medium">{m.name}</div>
                      <div className="text-xs text-muted-foreground">{m.email}</div>
                    </td>
                    <td className="py-3 pr-4 capitalize"><Badge variant={ROLE_VARIANT[m.role]}>{m.role}</Badge></td>
                    <td className="py-3 pr-4"><Badge variant={ST_VARIANT[m.status]}>{m.status}</Badge></td>
                    <td className="py-3 pr-4 text-muted-foreground">{m.lastSeen}</td>
                    <td className="py-3 text-right">
                      <Button variant="ghost" size="sm" className="text-xs">Edit</Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      )}

      {tab === "notifications" && (
        <Card>
          <CardHeader><CardTitle className="text-sm">Notification preferences</CardTitle></CardHeader>
          <CardContent className="space-y-4">
            {[
              { label: "Critical alerts", description: "Notify immediately for critical severity events", checked: true },
              { label: "High alerts", description: "Notify for high severity events", checked: true },
              { label: "Warning alerts", description: "Notify for warning events", checked: false },
              { label: "License expiry", description: "30 and 7 days before expiry", checked: true },
              { label: "Contract expiry", description: "90 and 30 days before expiry", checked: true },
              { label: "Agent offline", description: "When a monitoring agent stops responding", checked: true },
              { label: "Weekly digest", description: "Summary report every Monday morning", checked: false },
            ].map(n => (
              <div key={n.label} className="flex items-center justify-between py-1 border-b last:border-0">
                <div>
                  <p className="text-sm font-medium">{n.label}</p>
                  <p className="text-xs text-muted-foreground">{n.description}</p>
                </div>
                <input type="checkbox" defaultChecked={n.checked} className="h-4 w-4 accent-blue-600" />
              </div>
            ))}
            <div className="flex justify-end">
              <Button><Save className="h-4 w-4" /> Save preferences</Button>
            </div>
          </CardContent>
        </Card>
      )}

      {tab === "integrations" && (
        <div className="space-y-3">
          {INTEGRATIONS.map(i => (
            <Card key={i.name}>
              <CardContent className="flex items-center gap-4 p-4">
                <div className="rounded-lg bg-muted/60 p-2.5"><i.icon className="h-5 w-5 text-muted-foreground" /></div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium">{i.name}</p>
                  <p className="text-xs text-muted-foreground">{i.description}</p>
                </div>
                <div className="flex items-center gap-3">
                  <Badge variant={i.connected ? "success" : "secondary"}>{i.connected ? "Connected" : "Not connected"}</Badge>
                  <Button variant="outline" size="sm">{i.connected ? "Configure" : "Connect"}</Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {tab === "appearance" && (
        <Card>
          <CardHeader><CardTitle className="text-sm">Appearance & language</CardTitle></CardHeader>
          <CardContent className="space-y-5">
            <div>
              <Label className="text-sm font-medium">Theme</Label>
              <div className="mt-2 flex gap-3">
                <button onClick={() => theme !== "light" && toggleTheme()}
                  className={`flex items-center gap-2 rounded-lg border-2 px-4 py-2.5 text-sm font-medium transition-colors ${theme === "light" ? "border-primary bg-primary/5 text-primary" : "border-border text-muted-foreground hover:text-foreground"}`}>
                  <Sun className="h-4 w-4" /> Light
                </button>
                <button onClick={() => theme !== "dark" && toggleTheme()}
                  className={`flex items-center gap-2 rounded-lg border-2 px-4 py-2.5 text-sm font-medium transition-colors ${theme === "dark" ? "border-primary bg-primary/5 text-primary" : "border-border text-muted-foreground hover:text-foreground"}`}>
                  <Moon className="h-4 w-4" /> Dark
                </button>
              </div>
            </div>
            <div>
              <Label className="text-sm font-medium">Language</Label>
              <div className="mt-2 flex gap-3">
                {(["en", "ar"] as const).map(l => (
                  <button key={l} onClick={() => setLocale(l)}
                    className={`rounded-lg border-2 px-5 py-2.5 text-sm font-medium transition-colors ${locale === l ? "border-primary bg-primary/5 text-primary" : "border-border text-muted-foreground hover:text-foreground"}`}>
                    {l === "en" ? "English" : "العربية"}
                  </button>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
