"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { RedirectIfAuthed } from "@/components/marketing/redirect-if-authed";

const features = [
  { title: "Instant short links", body: "Generate unique short codes with a Postgres-backed ticket server. Sub-10ms cache hits via Redis." },
  { title: "Real analytics", body: "Track clicks, referrers, geographic data, and time-series usage on every link you create." },
  { title: "API-first", body: "Every dashboard action is a documented REST call. Issue scoped API keys for your own integrations." },
  { title: "Rate limiting", body: "Per-key and per-IP limits in Redis stop abuse without slowing real users down." },
  { title: "Built to scale", body: "Stateless Go containers, horizontal scaling tested to 8 replicas at 1k+ RPS." },
  { title: "Graceful degradation", body: "If Redis blinks, requests fall back to Postgres so your links keep resolving." },
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
              Slugify is a self-hostable URL shortener written in Go. It ships with API keys, rate limiting, click analytics, and a dashboard you can extend.
            </p>
            <div className="mt-10 flex flex-wrap items-center justify-center gap-3">
              <Link href="/signup"><Button size="lg">Create free account</Button></Link>
              <Link href="/login"><Button size="lg" variant="secondary">Sign in</Button></Link>
            </div>
          </div>

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
                  <div key={s.k} className="rounded-lg border border-[var(--border)] bg-[var(--surface-2)] px-4 py-3">
                    <p className="text-xs text-muted">{s.k}</p>
                    <p className="mt-1 text-2xl font-semibold tracking-tight">{s.v}</p>
                  </div>
                ))}
              </div>
              <div className="mt-6 h-32 rounded-lg border border-[var(--border)] bg-gradient-to-t from-indigo-500/20 to-transparent" />
            </div>
          </div>
        </div>
      </section>

      <section id="features" className="mx-auto max-w-6xl px-6 py-24">
        <div className="max-w-2xl">
          <p className="text-xs font-medium uppercase tracking-[0.2em] text-[var(--accent)]">Features</p>
          <h2 className="mt-3 text-3xl font-semibold tracking-tight sm:text-4xl">Everything you need to ship.</h2>
          <p className="mt-4 text-muted">A focused feature set covering the full URL-shortener lifecycle — creation, distribution, analytics, and operations.</p>
        </div>
        <div className="mt-12 grid gap-px overflow-hidden rounded-2xl border border-[var(--border)] bg-[var(--border)] sm:grid-cols-2 lg:grid-cols-3">
          {features.map((f) => (
            <div key={f.title} className="bg-[var(--background)] p-6 transition hover:bg-[var(--surface)]">
              <h3 className="font-medium tracking-tight">{f.title}</h3>
              <p className="mt-2 text-sm text-muted">{f.body}</p>
            </div>
          ))}
        </div>
      </section>

      <section id="how" className="border-y border-[var(--border)] bg-[var(--surface)]">
        <div className="mx-auto max-w-6xl px-6 py-24">
          <div className="max-w-2xl">
            <p className="text-xs font-medium uppercase tracking-[0.2em] text-[var(--accent)]">How it works</p>
            <h2 className="mt-3 text-3xl font-semibold tracking-tight sm:text-4xl">Three moving parts.</h2>
          </div>
          <ol className="mt-12 grid gap-6 sm:grid-cols-3">
            {[
              { n: "01", t: "Submit a long URL", b: "POST /api/shorten with your API key. We validate, sanitize, and reserve a unique short code." },
              { n: "02", t: "Redirect at the edge", b: "GET /:code reads from Redis first, falls back to Postgres, and 302s to the original URL." },
              { n: "03", t: "Analyze clicks", b: "Each redirect logs a hashed click event for daily, referrer, and country aggregations." },
            ].map((step) => (
              <li key={step.n} className="rounded-xl border border-[var(--border)] bg-[var(--background)] p-6">
                <span className="font-mono text-xs text-muted">{step.n}</span>
                <h3 className="mt-3 font-medium tracking-tight">{step.t}</h3>
                <p className="mt-2 text-sm text-muted">{step.b}</p>
              </li>
            ))}
          </ol>
        </div>
      </section>

      <section id="stats" className="mx-auto max-w-6xl px-6 py-24">
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {stats.map((s) => (
            <div key={s.label} className="rounded-xl border border-[var(--border)] p-6">
              <p className="text-3xl font-semibold tracking-tight">{s.value}</p>
              <p className="mt-2 text-sm text-muted">{s.label}</p>
            </div>
          ))}
        </div>
      </section>

      <section className="mx-auto max-w-6xl px-6 pb-24">
        <div className="overflow-hidden rounded-2xl border border-[var(--border)] bg-[var(--surface)] p-10 text-center">
          <h2 className="text-3xl font-semibold tracking-tight">Ready to shorten?</h2>
          <p className="mx-auto mt-3 max-w-xl text-muted">Spin up an account, drop in a URL, and watch the analytics roll in.</p>
          <div className="mt-8 flex justify-center gap-3">
            <Link href="/signup"><Button size="lg">Create free account</Button></Link>
            <Link href="/login"><Button size="lg" variant="ghost">I already have one →</Button></Link>
          </div>
        </div>
      </section>
    </main>
  );
}