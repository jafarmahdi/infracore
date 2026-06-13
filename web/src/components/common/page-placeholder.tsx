import { Construction } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

export function PagePlaceholder({ title }: { title: string }) {
  return (
    <Card>
      <CardContent className="flex min-h-[420px] flex-col items-center justify-center text-center">
        <div className="mb-4 rounded-xl bg-primary/10 p-3 text-primary">
          <Construction className="h-6 w-6" />
        </div>
        <h2 className="text-lg font-semibold">{title}</h2>
        <p className="mt-2 max-w-md text-sm text-muted-foreground">
          The module route and access boundary are ready. Data views and workflows will connect through the shared API service layer.
        </p>
      </CardContent>
    </Card>
  );
}
