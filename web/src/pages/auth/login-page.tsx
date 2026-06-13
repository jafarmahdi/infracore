import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useMutation } from "@tanstack/react-query";
import { isAxiosError } from "axios";
import { Navigate, useLocation, useNavigate } from "react-router-dom";
import { Building2, Eye, EyeOff, LockKeyhole, Server, ShieldCheck, Wifi } from "lucide-react";
import { toast } from "sonner";
import { login } from "@/services/api/auth";
import { useAuthStore } from "@/stores/auth-store";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Logo } from "@/components/common/logo";

const loginSchema = z.object({
  tenant_slug: z.string().min(1, "Tenant is required"),
  email: z.string().email("Enter a valid work email"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});
type LoginForm = z.infer<typeof loginSchema>;

export function LoginPage() {
  const [showPassword, setShowPassword] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const { isAuthenticated, setSession } = useAuthStore();

  const form = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
    defaultValues: { tenant_slug: "demo", email: "admin@demo.com", password: "" },
  });

  const mutation = useMutation({
    mutationFn: login,
    onSuccess: ({ user, access_token }) => {
      // Normalise to what the rest of the UI expects
      const enriched = {
        ...user,
        name: `${user.first_name} ${user.last_name}`.trim() || user.username,
        role: (user.is_superuser ? "admin" : "operator") as "admin" | "operator",
        company: user.tenant_slug,
      };
      setSession(enriched, access_token);
      toast.success(`Welcome back, ${enriched.first_name || enriched.username}`);
      const destination = (location.state as { from?: string } | null)?.from ?? "/dashboard";
      navigate(destination, { replace: true });
    },
    onError: (err: unknown) => {
      const msg = isAxiosError<{ message?: string }>(err)
        ? err.response?.data?.message
        : undefined;
      toast.error(msg ?? "Unable to sign in. Check your credentials.");
    },
  });

  if (isAuthenticated) return <Navigate to="/dashboard" replace />;

  return (
    <div className="grid min-h-screen bg-slate-950 lg:grid-cols-[1.1fr_0.9fr]">
      {/* Left panel */}
      <section className="relative hidden overflow-hidden p-12 lg:flex lg:flex-col">
        <div className="absolute inset-0 subtle-grid opacity-20" />
        <div className="absolute -left-32 top-20 h-96 w-96 rounded-full bg-blue-600/20 blur-[120px]" />
        <div className="absolute bottom-0 right-0 h-96 w-96 rounded-full bg-cyan-500/10 blur-[120px]" />
        <Logo className="relative z-10 text-white [&_div:last-child_div:last-child]:text-slate-500" />
        <div className="relative z-10 my-auto max-w-xl">
          <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-blue-400/20 bg-blue-400/10 px-3 py-1 text-xs font-medium text-blue-300">
            <span className="h-1.5 w-1.5 animate-pulse rounded-full bg-emerald-400" /> All systems operational
          </div>
          <h1 className="text-5xl font-semibold leading-[1.08] tracking-tight text-white">
            Infrastructure clarity.<br /><span className="text-blue-400">One control plane.</span>
          </h1>
          <p className="mt-6 max-w-lg text-base leading-7 text-slate-400">
            Unify data centers, network inventory, monitoring, assets, and service health in one secure operations platform.
          </p>
          <div className="mt-10 grid grid-cols-3 gap-3">
            {[
              [Server, "DCIM", "Data center infra"],
              [Wifi, "IPAM", "IP management"],
              [ShieldCheck, "24/7", "Operations"],
            ].map(([Icon, value, label]) => {
              const Component = Icon as typeof Server;
              return (
                <div key={String(label)} className="rounded-xl border border-white/10 bg-white/[0.04] p-4 backdrop-blur">
                  <Component className="mb-3 h-5 w-5 text-blue-400" />
                  <div className="text-lg font-semibold text-white">{String(value)}</div>
                  <div className="mt-0.5 text-[11px] text-slate-500">{String(label)}</div>
                </div>
              );
            })}
          </div>
        </div>
        <p className="relative z-10 text-xs text-slate-600">InfraCore Enterprise Platform · Secure operations access</p>
      </section>

      {/* Right panel */}
      <section className="flex items-center justify-center bg-background px-6 py-12">
        <div className="w-full max-w-[420px]">
          <Logo className="mb-12 lg:hidden" />
          <div className="mb-8">
            <h2 className="text-2xl font-semibold tracking-tight">Sign in to InfraCore</h2>
            <p className="mt-2 text-sm text-muted-foreground">Use your organization credentials to continue.</p>
          </div>

          <form className="space-y-5" onSubmit={form.handleSubmit((v) => mutation.mutate(v))}>
            <div className="space-y-2">
              <Label htmlFor="tenant_slug">Organization</Label>
              <div className="relative">
                <Building2 className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input id="tenant_slug" autoComplete="organization" className="pl-9" placeholder="demo" {...form.register("tenant_slug")} />
              </div>
              {form.formState.errors.tenant_slug && <p className="text-xs text-red-500">{form.formState.errors.tenant_slug.message}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="email">Work email</Label>
              <Input id="email" autoComplete="email" placeholder="admin@demo.com" {...form.register("email")} />
              {form.formState.errors.email && <p className="text-xs text-red-500">{form.formState.errors.email.message}</p>}
            </div>

            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <Label htmlFor="password">Password</Label>
                <button type="button" className="text-xs font-medium text-primary hover:underline">Forgot password?</button>
              </div>
              <div className="relative">
                <LockKeyhole className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input id="password" type={showPassword ? "text" : "password"} autoComplete="current-password" className="pl-9 pr-10" {...form.register("password")} />
                <button type="button" onClick={() => setShowPassword(!showPassword)} className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground">
                  {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
              {form.formState.errors.password && <p className="text-xs text-red-500">{form.formState.errors.password.message}</p>}
            </div>

            <Button className="w-full" type="submit" disabled={mutation.isPending}>
              {mutation.isPending ? "Authenticating…" : "Sign in securely"}
            </Button>
          </form>

          <div className="mt-8 flex items-center justify-center gap-2 text-xs text-muted-foreground">
            <ShieldCheck className="h-4 w-4 text-emerald-500" /> Protected by enterprise security policies
          </div>
        </div>
      </section>
    </div>
  );
}
