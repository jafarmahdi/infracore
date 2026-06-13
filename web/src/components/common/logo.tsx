import { cn } from "@/lib/utils";

export function Logo({ compact = false, className }: { compact?: boolean; className?: string }) {
  return (
    <div className={cn("flex items-center gap-3", className)}>
      <span className="grid h-11 w-11 shrink-0 place-items-center rounded-lg bg-white p-1 shadow-sm ring-1 ring-slate-200/80">
        <img
          src="/brand/infracore-mark.png"
          alt="InfraCore"
          className="h-full w-full object-contain"
        />
      </span>
      {!compact && (
        <div>
          <div className="text-[15px] font-bold tracking-tight">InfraCore</div>
          <div className="text-[8px] font-semibold uppercase tracking-[0.14em] text-muted-foreground">
            Infrastructure management
          </div>
        </div>
      )}
    </div>
  );
}
