import { Cell, Pie, PieChart, ResponsiveContainer } from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { formatNumber } from "@/lib/utils";

export function HealthPanel({ data }: { data: { name: string; value: number; color: string }[] }) {
  const total = data.reduce((sum, item) => sum + item.value, 0);
  return (
    <Card>
      <CardHeader><CardTitle className="text-sm">Device health</CardTitle><CardDescription>Current managed estate status</CardDescription></CardHeader>
      <CardContent>
        <div className="relative mx-auto h-[145px]">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart><Pie data={data} dataKey="value" innerRadius={47} outerRadius={65} paddingAngle={3} stroke="none">{data.map((item) => <Cell key={item.name} fill={item.color} />)}</Pie></PieChart>
          </ResponsiveContainer>
          <div className="pointer-events-none absolute inset-0 grid place-items-center text-center">
            <div><div className="text-xl font-semibold">{formatNumber(total)}</div><div className="text-[10px] text-muted-foreground">Devices</div></div>
          </div>
        </div>
        <div className="mt-3 space-y-2.5">
          {data.map((item) => (
            <div key={item.name} className="flex items-center text-xs"><span className="mr-2 h-2 w-2 rounded-full" style={{ backgroundColor: item.color }} /><span className="text-muted-foreground">{item.name}</span><span className="ml-auto font-semibold">{formatNumber(item.value)}</span></div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
