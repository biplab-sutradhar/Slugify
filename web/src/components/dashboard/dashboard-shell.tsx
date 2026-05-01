"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";
import { useAuth } from "@/lib/auth-context";

const items = [
  { href: "/dashboard", label: "Overview" },
  { href: "/dashboard/links", label: "Links" },
  { href: "/dashboard/keys", label: "API keys" },
  { href: "/dashboard/analytics", label: "Analytics" },
];

export function DashboardShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const {
    user,
    loading,
    logout,
    ensureApiKey,
    apiKey,
    provisioning,
    provisioningError,
  } = useAuth();

  useEffect(() => {
    if (!loading && !user) router.replace("/login");
  }, [loading, user, router]);

  useEffect(() => {
    if (user && !apiKey && !provisioning && !provisioningError) {
      void ensureApiKey();
    }
  }, [user, apiKey, provisioning, provisioningError, ensureApiKey]);

  if (loading || !user) {
    return (
      <div className="flex min-h-screen items-center justify-center text-sm text-muted">
        Loading…
      </div>
    );
  }

  return (
    <div className="grid min-h-screen grid-cols-1 lg:grid-cols-[16rem_1fr]">
      <aside className="hidden border-r border-[var(--border)] bg-[var(--surface)] p-5 lg:flex lg:flex-col">
        <Link href="/dashboard" className="flex items-center gap-2">
          <span className="grid h-7 w-7 place-items-center rounded-md bg-[var(--accent)] text-[var(--accent-foreground)] text-xs font-bold">
            S
          </span>
          <span className="text-sm font-semibold tracking-tight">Slugify</span>
        </Link>

        <nav className="mt-8 flex flex-col gap-1 text-sm">
          {items.map((item) => {
            const active =
              pathname === item.href ||
              (item.href !== "/dashboard" && pathname.startsWith(item.href));
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`rounded-lg px-3 py-2 transition ${
                  active
                    ? "bg-[var(--surface-2)] text-foreground"
                    : "text-muted hover:bg-[var(--surface-2)] hover:text-foreground"
                }`}
              >
                {item.label}
              </Link>
            );
          })}
        </nav>

        <div className="mt-auto rounded-lg border border-[var(--border)] p-3 text-sm">
          <p className="truncate font-medium">{user.name || user.email}</p>
          <p className="truncate text-xs text-muted">{user.email}</p>
          <button
            type="button"
            onClick={() => {
              logout();
              router.push("/");
            }}
            className="mt-2 text-xs text-muted underline underline-offset-4 hover:text-foreground"
          >
            Sign out
          </button>
        </div>
      </aside>

      <div className="flex min-w-0 flex-col">
        {provisioningError && (
          <div className="border-b border-red-500/30 bg-red-500/10 px-6 py-3 text-sm text-red-700 dark:text-red-300">
            <span className="font-medium">Couldn't load your API key:</span>{" "}
            {provisioningError}{" "}
            <button
              type="button"
              onClick={() => void ensureApiKey()}
              className="ml-2 underline underline-offset-4"
            >
              Retry
            </button>
          </div>
        )}
        {!apiKey && provisioning && !provisioningError && (
          <div className="border-b border-[var(--border)] bg-[var(--surface)] px-6 py-3 text-sm text-muted">
            Provisioning your API key…
          </div>
        )}
        <main className="flex-1 bg-[var(--background)]">{children}</main>
      </div>
    </div>
  );
}