"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { RedirectIfAuthed } from "@/components/marketing/redirect-if-authed";

const features = [
  {
    title: "Instant short links",
    body: "Generate unique short codes with a Postgres-backed ticket server. Sub-10ms cache hits via Redis.",
  },
  {
    title: "Real analytics",
    body: "Track clicks, referrers, geographic data, and time-series usage on every link you create.",
  },
  {
    title: "API-first",
    body: "Every dashboard action is a documented REST call. Issue scoped API keys for your own integrations.",
  },
  {
    title: "Rate limiting",
    body: "Per-key and per-IP limits in Redis stop abuse without slowing real users down.",
  },
  {
    title: "Built to scale",
    body: "Stateless Go containers, horizontal scaling tested to 8 replicas at 1k+ RPS.",
  },
  {
    title: "Graceful degradation",
    body: "If Redis blinks, requests fall back to Postgres so your links keep resolving.",
  },
];

const stats = [
  { label: "Avg. cache hit", value: "<10ms" },
  { label: "Throughput / instance", value: "200 RPS" },
  { label: "p95 latency", value: "<100ms" },
  { label: "Uptime target", value: "99.9%" },
];

export default function LandingPage() {
  return (
    <main>
      <RedirectIfAuthed />
      {/* Hero */}
      <section className="relative overflow-hidden">
        <div className="gradient-mesh absolute inset-0 -z-10" />
        <div className="mx-auto max-w-6xl px-6 py-24 sm:py-32">
          <div className="mx-auto max-w-3xl text-center">
            <span className="inline-flex items-center gap-2 rounded-full border border-[var(--border)] bg-[var(--surface)]/60 px-3 py-1 text-xs text-muted backdrop-blur">
              <span className="h-1.5 w-1.5 rounded-full bg-[var(--success)]" />
              Production-ready · Open source
            </span>
            <h1 className="mt-6 text-4xl font-semibold tracking-tight sm:text-6xl">
              Short links with
              <span className="block bg-gradient-to-r from-[var(--accent)] to-pink-500 bg-clip-text text-transparent">
                real analytics built in.
              </span>
            </h1>
            <p className="mx-auto mt-6 max-w-2xl text-base text-muted sm:text-lg">
              Slugify is a self-hostable URL shortener written in Go. It ships
              with API keys, rate limiting, click analytics, and a dashboard you
              can extend.
            </p>
            <div className="mt-10 flex flex-wrap items-center justify-center gap-3">
              <Link href="/signup">
                <Button size="lg">Create free account</Button>
              </Link>
              <Link href="/login">
                <Button size="lg" variant="secondary">
                  Sign in
                </Button>
              </Link>
            </div>
          </div>

          {/* Mock dashboard preview */}
          <div className="mx-auto mt-16 max-w-4xl rounded-2xl border border-[var(--border)] bg-[var(--surface)] p-2 shadow-2xl shadow-indigo-500/10">
            <div className="rounded-xl border border-[var(--border)] bg-[var(--background)] p-6">
              <div className="flex items-center gap-2 border-b border-[var(--border)] pb-3">
                <span className="h-2.5 w-2.5 rounded-full bg-red-400/70" />
                <span className="h-2.5 w-2.5 rounded-full bg-yellow-400/70" />
                <span className="h-2.5 w-2.5 rounded-full bg-green-400/70" />
                <span className="ml-3 text-xs text-muted">slugify.app/dashboard</span>
              </div>
              <div className="grid gap-4 pt-6 sm:grid-cols-3">
                {[
                  { k: "Total clicks", v: "128,402" },
                  { k: "Active links", v: "1,204" },
                  { k: "Cache hit rate", v: "82%" },
                ].map((s) => (
                  <div
                    key={s.k}
                    className="rounded-lg border border-[var(--border)] bg-[var(--surface-2)] px-4 py-3"
                  >
                    <p className="text-xs text-muted">{s.k}</p>
                    <p className="mt-1 text-2xl font-semibold tracking-tight">
                      {s.v}
                    </p>
                  </div>
                ))}
              </div>
              <div className="mt-6 h-32 rounded-lg border border-[var(--border)] bg-gradient-to-t from-indigo-500/20 to-transparent" />
            </div>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="border-t border-[var(--border)] px-6 py-16 sm:py-20">
        <div className="mx-auto max-w-5xl">
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="text-3xl font-semibold tracking-tight">
              Everything you need.
            </h2>
            <p className="mt-4 text-base text-muted">
              Built for developers and powered by battle-tested infrastructure.
            </p>
          </div>
          <div className="mt-16 grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
            {features.map((f) => (
              <div key={f.title}>
                <h3 className="font-medium tracking-tight">{f.title}</h3>
                <p className="mt-2 text-sm text-muted">{f.body}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Stats */}
      <section className="border-t border-[var(--border)] px-6 py-16 sm:py-20">
        <div className="mx-auto max-w-5xl">
          <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
            {stats.map((s) => (
              <div key={s.label}>
                <p className="text-xs text-muted">{s.label}</p>
                <p className="mt-2 text-3xl font-semibold tracking-tight">
                  {s.value}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="border-t border-[var(--border)] px-6 py-16 sm:py-20">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="text-3xl font-semibold tracking-tight">
            Ready to get started?
          </h2>
          <p className="mt-4 text-base text-muted">
            Deploy Slugify in minutes. No credit card required.
          </p>
          <div className="mt-8 flex flex-wrap items-center justify-center gap-3">
            <Link href="/signup">
              <Button size="lg">Create account</Button>
            </Link>
            <Link href="/login">
              <Button size="lg" variant="secondary">
                Sign in
              </Button>
            </Link>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-[var(--border)] px-6 py-8">
        <div className="mx-auto max-w-5xl text-center text-sm text-muted">
          <p>Built with Go, Postgres, and Redis.</p>
        </div>
      </footer>
    </main>
  );
}