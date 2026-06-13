import { CartesianGrid, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import type { UtilizationPoint } from "@/types/dashboard";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export function OperationsChart({ data }: { data: UtilizationPoint[] }) {
  return (
    <Card className="lg:col-span-2">
      <CardHeader className="flex-row items-center justify-between space-y-0">
        <div><CardTitle className="text-sm">Resource utilization</CardTitle><CardDescription>Aggregate infrastructure performance</CardDescription></div>
        <Button variant="outline" size="sm">Last 24 hours</Button>
      </CardHeader>
      <CardContent>
        <div className="mb-4 flex gap-5 text-xs text-muted-foreground">
          {[["CPU", "#3b82f6"], ["Memory", "#8b5cf6"], ["Disk", "#14b8a6"]].map(([label, color]) => (
            <span key={label} className="flex items-center gap-1.5"><span className="h-2 w-2 rounded-full" style={{ backgroundColor: color }} />{label}</span>
          ))}
        </div>
        <ResponsiveContainer width="100%" height={250}>
          <LineChart data={data}>
            <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="hsl(var(--border))" />
            <XAxis dataKey="time" axisLine={false} tickLine={false} tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }} />
            <YAxis unit="%" axisLine={false} tickLine={false} tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }} />
            <Tooltip contentStyle={{ borderRadius: 10, border: "1px solid hsl(var(--border))", background: "hsl(var(--card))", fontSize: 12 }} />
            <Line type="monotone" dataKey="cpu" stroke="#3b82f6" strokeWidth={2} dot={false} />
            <Line type="monotone" dataKey="memory" stroke="#8b5cf6" strokeWidth={2} dot={false} />
            <Line type="monotone" dataKey="disk" stroke="#14b8a6" strokeWidth={2} dot={false} />
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
