import { ExternalLink } from "lucide-react";
import type { RecentEvent } from "@/types/dashboard";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { formatTime } from "@/lib/utils";

const variants = { critical: "destructive", high: "destructive", warning: "warning", info: "default" } as const;

export function RecentEvents({ events }: { events: RecentEvent[] }) {
  return (
    <Card className="lg:col-span-2">
      <CardHeader className="flex-row items-center justify-between space-y-0">
        <div><CardTitle className="text-sm">Recent events</CardTitle><CardDescription>Latest changes across integrated systems</CardDescription></div>
        <Button variant="ghost" size="sm">View all <ExternalLink className="h-3.5 w-3.5" /></Button>
      </CardHeader>
      <CardContent className="px-0 pb-1">
        <div className="divide-y">
          {events.map((event) => (
            <div key={event.id} className="flex items-center gap-3 px-5 py-3.5 hover:bg-muted/50">
              <span className={`h-2 w-2 shrink-0 rounded-full ${event.severity === "critical" ? "bg-red-500" : event.severity === "warning" ? "bg-amber-500" : event.severity === "high" ? "bg-orange-500" : "bg-blue-500"}`} />
              <div className="min-w-0 flex-1">
                <p className="truncate text-sm font-medium">{event.title}</p>
                <p className="mt-0.5 text-[11px] text-muted-foreground">{event.source} · {formatTime(event.timestamp)}</p>
              </div>
              <Badge variant={variants[event.severity]}>{event.severity}</Badge>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
