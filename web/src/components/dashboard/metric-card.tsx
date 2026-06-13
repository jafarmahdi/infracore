import type { LucideIcon } from "lucide-react";
import { ArrowDownRight, ArrowUpRight, Minus } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { cn, formatNumber } from "@/lib/utils";
import type { TrendDirection } from "@/types/dashboard";

interface MetricCardProps {
  title: string;
  value: number;
  suffix?: string;
  description: string;
  icon: LucideIcon;
  iconClassName: string;
  trend?: { value: number; direction: TrendDirection };
}

export function MetricCard({ title, value, suffix, description, icon: Icon, iconClassName, trend }: MetricCardProps) {
  const TrendIcon = trend?.direction === "up" ? ArrowUpRight : trend?.direction === "down" ? ArrowDownRight : Minus;
  return (
    <Card className="overflow-hidden">
      <CardContent className="p-5">
        <div className="flex items-start justify-between">
          <div>
            <p className="text-xs font-medium text-muted-foreground">{title}</p>
            <div className="mt-2 flex items-end gap-1">
              <span className="text-2xl font-semibold tracking-tight">{formatNumber(value)}</span>
              {suffix && <span className="mb-0.5 text-sm font-medium text-muted-foreground">{suffix}</span>}
            </div>
          </div>
          <div className={cn("rounded-lg p-2.5", iconClassName)}><Icon className="h-5 w-5" /></div>
        </div>
        <div className="mt-4 flex items-center gap-1.5 text-[11px] text-muted-foreground">
          {trend && (
            <span className={cn("flex items-center font-semibold", trend.direction === "down" ? "text-emerald-500" : "text-blue-500")}>
              <TrendIcon className="h-3.5 w-3.5" />{trend.value}%
            </span>
          )}
          {description}
        </div>
      </CardContent>
    </Card>
  );
}
