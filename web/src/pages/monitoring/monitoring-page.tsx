import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { isAxiosError } from "axios";
import {
  Activity,
  AlertTriangle,
  CheckCircle2,
  ChevronLeft,
  ChevronRight,
  CircleDashed,
  Loader2,
  Plus,
  ServerCrash,
  Trash2,
  WrenchIcon,
  XCircle,
} from "lucide-react";
import { toast } from "sonner";
import { createHost, deleteHost, getStatusCounts, listHosts } from "@/services/api/monitoring";
import type { CreateHostRequest, HostStatus, MonitoringType } from "@/types/monitoring";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Textarea } from "@/components/ui/textarea";

// ─── Add Host form schema ───────────────────────────────────────────────────

const addHostSchema = z.object({
  name: z.string().min(1, "Name is required").max(255),
  display_name: z.string().optional(),
  ip_address: z.string().optional(),
  monitoring_type: z.enum(["agent", "snmp", "wmi", "icmp", "ssh", "api"]),
  port: z.coerce.number().int().min(1).max(65535).optional().or(z.literal("")),
  check_interval_seconds: z.coerce.number().int().min(10).default(60),
  timeout_seconds: z.coerce.number().int().min(1).default(30),
  snmp_version: z.string().optional(),
  snmp_community: z.string().optional(),
  notes: z.string().optional(),
});
type AddHostForm = z.infer<typeof addHostSchema>;

// ─── Helpers ────────────────────────────────────────────────────────────────

const STATUS_LABELS: Record<HostStatus, string> = {
  up: "Up",
  down: "Down",
  warning: "Warning",
  maintenance: "Maintenance",
  unknown: "Unknown",
};

function StatusBadge({ status }: { status: HostStatus }) {
  const map: Record<HostStatus, { cls: string; Icon: React.ElementType }> = {
    up: { cls: "bg-emerald-500/10 text-emerald-600 border-emerald-500/20", Icon: CheckCircle2 },
    down: { cls: "bg-red-500/10 text-red-600 border-red-500/20", Icon: XCircle },
    warning: { cls: "bg-amber-500/10 text-amber-600 border-amber-500/20", Icon: AlertTriangle },
    maintenance: { cls: "bg-blue-500/10 text-blue-600 border-blue-500/20", Icon: WrenchIcon },
    unknown: { cls: "bg-slate-500/10 text-slate-600 border-slate-500/20", Icon: CircleDashed },
  };
  const { cls, Icon } = map[status] ?? map.unknown;
  return (
    <span className={`inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium ${cls}`}>
      <Icon className="h-3 w-3" />
      {STATUS_LABELS[status] ?? status}
    </span>
  );
}

const FILTER_OPTIONS: { label: string; value: string | undefined }[] = [
  { label: "All", value: undefined },
  { label: "Up", value: "up" },
  { label: "Warning", value: "warning" },
  { label: "Down", value: "down" },
  { label: "Maintenance", value: "maintenance" },
];

// ─── Page ────────────────────────────────────────────────────────────────────

export function MonitoringPage() {
  const [statusFilter, setStatusFilter] = useState<string | undefined>(undefined);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [addOpen, setAddOpen] = useState(false);
  const [deleteId, setDeleteId] = useState<string | null>(null);

  const queryClient = useQueryClient();

  const hostsQuery = useQuery({
    queryKey: ["hosts", page, statusFilter, search],
    queryFn: () => listHosts({ page, page_size: 20, status: statusFilter, search: search || undefined }),
  });

  const countsQuery = useQuery({
    queryKey: ["host-counts"],
    queryFn: getStatusCounts,
  });

  const form = useForm<AddHostForm>({
    resolver: zodResolver(addHostSchema),
    defaultValues: {
      monitoring_type: "icmp",
      check_interval_seconds: 60,
      timeout_seconds: 30,
    },
  });

  const monType = form.watch("monitoring_type") as MonitoringType;

  const createMutation = useMutation({
    mutationFn: (body: CreateHostRequest) => createHost(body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["hosts"] });
      queryClient.invalidateQueries({ queryKey: ["host-counts"] });
      toast.success("Host added successfully");
      setAddOpen(false);
      form.reset();
    },
    onError: (err: unknown) => {
      const message = isAxiosError<{ message?: string }>(err)
        ? err.response?.data?.message
        : undefined;
      toast.error(message ?? "Failed to add host");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteHost,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["hosts"] });
      queryClient.invalidateQueries({ queryKey: ["host-counts"] });
      toast.success("Host removed");
      setDeleteId(null);
    },
    onError: () => toast.error("Failed to remove host"),
  });

  function onSubmit(values: AddHostForm) {
    const body: CreateHostRequest = {
      name: values.name,
      monitoring_type: values.monitoring_type,
      ...(values.display_name && { display_name: values.display_name }),
      ...(values.ip_address && { ip_address: values.ip_address }),
      ...(values.port && { port: Number(values.port) }),
      check_interval_seconds: values.check_interval_seconds,
      timeout_seconds: values.timeout_seconds,
      ...(values.snmp_version && { snmp_version: values.snmp_version }),
      ...(values.snmp_community && { snmp_community: values.snmp_community }),
      ...(values.notes && { notes: values.notes }),
    };
    createMutation.mutate(body);
  }

  const hosts = hostsQuery.data?.items ?? [];
  const counts = countsQuery.data;
  const totalPages = hostsQuery.data?.total_pages ?? 1;

  return (
    <div className="space-y-6 p-6">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Monitoring</h1>
          <p className="mt-1 text-sm text-muted-foreground">Track host health and availability in real time.</p>
        </div>
        <Button onClick={() => { form.reset(); setAddOpen(true); }}>
          <Plus className="mr-2 h-4 w-4" /> Add Host
        </Button>
      </div>

      {/* Status summary cards */}
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-5">
        {([
          ["up", "Up", CheckCircle2, "text-emerald-500"],
          ["down", "Down", XCircle, "text-red-500"],
          ["warning", "Warning", AlertTriangle, "text-amber-500"],
          ["maintenance", "Maintenance", WrenchIcon, "text-blue-500"],
          ["unknown", "Unknown", CircleDashed, "text-slate-400"],
        ] as const).map(([key, label, Icon, iconCls]) => (
          <div key={key} className="rounded-lg border bg-card p-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">{label}</span>
              <Icon className={`h-4 w-4 ${iconCls}`} />
            </div>
            <div className="mt-1 text-2xl font-semibold">
              {countsQuery.isLoading ? "—" : (counts?.[key] ?? 0)}
            </div>
          </div>
        ))}
      </div>

      {/* Filters */}
      <div className="flex flex-wrap items-center gap-3">
        <div className="flex gap-1 rounded-lg border bg-muted p-1">
          {FILTER_OPTIONS.map(({ label, value }) => (
            <button
              key={label}
              onClick={() => { setStatusFilter(value); setPage(1); }}
              className={`rounded-md px-3 py-1.5 text-xs font-medium transition-colors ${
                statusFilter === value
                  ? "bg-background shadow text-foreground"
                  : "text-muted-foreground hover:text-foreground"
              }`}
            >
              {label}
            </button>
          ))}
        </div>
        <Input
          placeholder="Search hosts…"
          value={search}
          onChange={(e) => { setSearch(e.target.value); setPage(1); }}
          className="h-8 w-52 text-sm"
        />
        {hostsQuery.isFetching && <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />}
      </div>

      {/* Table */}
      <div className="rounded-lg border bg-card">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Host</TableHead>
              <TableHead>IP Address</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Last Status</TableHead>
              <TableHead>Last Check</TableHead>
              <TableHead className="w-12" />
            </TableRow>
          </TableHeader>
          <TableBody>
            {hostsQuery.isLoading ? (
              <TableRow>
                <TableCell colSpan={7} className="py-16 text-center">
                  <Loader2 className="mx-auto h-6 w-6 animate-spin text-muted-foreground" />
                </TableCell>
              </TableRow>
            ) : hostsQuery.isError ? (
              <TableRow>
                <TableCell colSpan={7} className="py-16 text-center">
                  <div className="flex flex-col items-center gap-2 text-muted-foreground">
                    <AlertTriangle className="h-8 w-8 text-red-400" />
                    <p className="text-sm">Failed to load hosts. Check the API connection.</p>
                    <Button variant="outline" size="sm" onClick={() => hostsQuery.refetch()}>Retry</Button>
                  </div>
                </TableCell>
              </TableRow>
            ) : hosts.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className="py-20 text-center">
                  <div className="flex flex-col items-center gap-3 text-muted-foreground">
                    <ServerCrash className="h-10 w-10 text-slate-300 dark:text-slate-600" />
                    <p className="font-medium">No hosts yet</p>
                    <p className="text-sm">Add your first monitored host to start tracking availability.</p>
                    <Button size="sm" onClick={() => { form.reset(); setAddOpen(true); }}>
                      <Plus className="mr-2 h-3.5 w-3.5" /> Add Host
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ) : (
              hosts.map((host) => (
                <TableRow key={host.id}>
                  <TableCell className="font-medium">
                    <div>{host.display_name ?? host.name}</div>
                    {host.display_name && (
                      <div className="text-xs text-muted-foreground">{host.name}</div>
                    )}
                  </TableCell>
                  <TableCell className="text-sm text-muted-foreground">{host.ip_address ?? "—"}</TableCell>
                  <TableCell>
                    <Badge variant="secondary" className="uppercase text-[11px]">
                      {host.monitoring_type}
                    </Badge>
                  </TableCell>
                  <TableCell><StatusBadge status={host.status} /></TableCell>
                  <TableCell className="text-sm text-muted-foreground">{host.last_status ?? "—"}</TableCell>
                  <TableCell className="text-sm text-muted-foreground">
                    {host.last_check_at
                      ? new Date(host.last_check_at).toLocaleString()
                      : "Not checked yet"}
                  </TableCell>
                  <TableCell>
                    <button
                      onClick={() => setDeleteId(host.id)}
                      className="rounded p-1 text-muted-foreground hover:bg-destructive/10 hover:text-destructive"
                    >
                      <Trash2 className="h-4 w-4" />
                    </button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between border-t px-4 py-3">
            <p className="text-sm text-muted-foreground">
              Page {page} of {totalPages} · {hostsQuery.data?.total_items ?? 0} hosts
            </p>
            <div className="flex gap-1">
              <Button variant="outline" size="icon" className="h-8 w-8" disabled={page <= 1} onClick={() => setPage(p => p - 1)}>
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <Button variant="outline" size="icon" className="h-8 w-8" disabled={page >= totalPages} onClick={() => setPage(p => p + 1)}>
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </div>
        )}
      </div>

      {/* Add Host modal */}
      <Dialog open={addOpen} onOpenChange={setAddOpen}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Activity className="h-5 w-5 text-primary" /> Add Monitored Host
            </DialogTitle>
          </DialogHeader>

          <form id="add-host-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="col-span-2 space-y-1.5">
                <Label>Host name <span className="text-red-500">*</span></Label>
                <Input placeholder="web-server-01" {...form.register("name")} />
                {form.formState.errors.name && (
                  <p className="text-xs text-red-500">{form.formState.errors.name.message}</p>
                )}
              </div>

              <div className="space-y-1.5">
                <Label>Display name</Label>
                <Input placeholder="Production Web" {...form.register("display_name")} />
              </div>

              <div className="space-y-1.5">
                <Label>IP / Hostname</Label>
                <Input placeholder="192.168.1.10" {...form.register("ip_address")} />
              </div>

              <div className="space-y-1.5">
                <Label>Monitoring type <span className="text-red-500">*</span></Label>
                <Select
                  value={form.watch("monitoring_type")}
                  onValueChange={(v) => form.setValue("monitoring_type", v as MonitoringType)}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {(["icmp", "snmp", "ssh", "wmi", "agent", "api"] as MonitoringType[]).map((t) => (
                      <SelectItem key={t} value={t}>{t.toUpperCase()}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-1.5">
                <Label>Port</Label>
                <Input type="number" placeholder="optional" {...form.register("port")} />
              </div>

              <div className="space-y-1.5">
                <Label>Check interval (s)</Label>
                <Input type="number" {...form.register("check_interval_seconds")} />
              </div>

              <div className="space-y-1.5">
                <Label>Timeout (s)</Label>
                <Input type="number" {...form.register("timeout_seconds")} />
              </div>

              {monType === "snmp" && (
                <>
                  <div className="space-y-1.5">
                    <Label>SNMP version</Label>
                    <Select onValueChange={(v) => form.setValue("snmp_version", v)}>
                      <SelectTrigger><SelectValue placeholder="v2c" /></SelectTrigger>
                      <SelectContent>
                        <SelectItem value="v1">v1</SelectItem>
                        <SelectItem value="v2c">v2c</SelectItem>
                        <SelectItem value="v3">v3</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-1.5">
                    <Label>Community string</Label>
                    <Input placeholder="public" {...form.register("snmp_community")} />
                  </div>
                </>
              )}

              <div className="col-span-2 space-y-1.5">
                <Label>Notes</Label>
                <Textarea rows={2} placeholder="Optional notes…" {...form.register("notes")} />
              </div>
            </div>
          </form>

          <DialogFooter>
            <Button variant="outline" onClick={() => setAddOpen(false)} disabled={createMutation.isPending}>
              Cancel
            </Button>
            <Button type="submit" form="add-host-form" disabled={createMutation.isPending}>
              {createMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Add Host
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete confirm modal */}
      <Dialog open={!!deleteId} onOpenChange={(o) => !o && setDeleteId(null)}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Remove host?</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            This will permanently delete the host and all its monitoring history.
          </p>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteId(null)} disabled={deleteMutation.isPending}>
              Cancel
            </Button>
            <Button
              variant="destructive"
              disabled={deleteMutation.isPending}
              onClick={() => deleteId && deleteMutation.mutate(deleteId)}
            >
              {deleteMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
