'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  listLinks,
  normalizeUrl,
  shorten,
  type Link as LinkRow,
} from '@/lib/slugify-api';
import { useAuth } from '@/lib/auth-context';

function Sparkline({ values }: { values: number[] }) {
  const max = Math.max(1, ...values);
  const points = values
    .map(
      (v, i) =>
        `${(i / (values.length - 1 || 1)) * 100},${100 - (v / max) * 100}`,
    )
    .join(' ');
  return (
    <svg
      viewBox="0 0 100 100"
      preserveAspectRatio="none"
      className="h-16 w-full"
    >
      <polyline
        fill="none"
        stroke="var(--accent)"
        strokeWidth="2"
        vectorEffect="non-scaling-stroke"
        points={points}
      />
    </svg>
  );
}

export default function DashboardOverview() {
  const { apiKey } = useAuth();
  const [links, setLinks] = useState<LinkRow[]>([]);
  const [longUrl, setLongUrl] = useState('');
  const [shortUrl, setShortUrl] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const refresh = useCallback(async () => {
    if (!apiKey) {
      setLinks([]);
      return;
    }
    setError(null);
    setLoading(true);
    try {
      const data = await listLinks(apiKey, { limit: 5, offset: 0 });
      setLinks(data);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load links');
    } finally {
      setLoading(false);
    }
  }, [apiKey]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  const stats = useMemo(() => {
    const totalClicks = links.reduce((acc, l) => acc + l.clicks, 0);
    const active = links.filter((l) => l.is_active).length;
    return { totalClicks, active, total: links.length };
  }, [links]);

  const onShorten = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setShortUrl(null);

    if (!apiKey) {
      setError(
        'Your API key is still being provisioned. Try again in a moment.',
      );
      return;
    }

    const normalized = normalizeUrl(longUrl);
    if (!normalized) {
      setError("That doesn't look like a valid URL.");
      return;
    }

    setLoading(true);
    try {
      const r = await shorten(normalized, apiKey);
      setShortUrl(r.short_url);
      setLongUrl('');
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Shorten failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mx-auto max-w-5xl px-6 py-10">
      <div className="flex items-end justify-between">
        <div>
          <p className="text-xs uppercase tracking-[0.2em] text-muted">
            Overview
          </p>
          <h1 className="mt-1 text-2xl font-semibold tracking-tight">
            Dashboard
          </h1>
        </div>
      </div>

      <section className="mt-8 grid gap-4 sm:grid-cols-3">
        {[
          { label: 'Total clicks', value: stats.totalClicks.toLocaleString() },
          { label: 'Active links', value: stats.active.toLocaleString() },
          { label: 'Total links', value: stats.total.toLocaleString() },
        ].map((s) => (
          <div
            key={s.label}
            className="rounded-xl border border-[var(--border)] bg-[var(--surface)] p-5"
          >
            <p className="text-xs text-muted">{s.label}</p>
            <p className="mt-2 text-3xl font-semibold tracking-tight">
              {s.value}
            </p>
          </div>
        ))}
      </section>

      <section className="mt-6 grid gap-4 lg:grid-cols-3">
        <div className="rounded-xl border border-[var(--border)] bg-[var(--surface)] p-5 lg:col-span-2">
          <p className="text-xs text-muted">Clicks across recent links</p>
          {links.length > 0 ? (
            <Sparkline values={links.map((l) => l.clicks)} />
          ) : (
            <p className="mt-4 text-sm text-muted">No data yet.</p>
          )}
        </div>

        <form
          onSubmit={onShorten}
          className="rounded-xl border border-[var(--border)] bg-[var(--surface)] p-5"
        >
          <p className="text-sm font-medium">Quick shorten</p>
          <Input
            type="text"
            inputMode="url"
            autoComplete="off"
            spellCheck={false}
            required
            placeholder="example.com"
            value={longUrl}
            onChange={(e) => setLongUrl(e.target.value)}
            className="mt-3"
          />
          <Button type="submit" disabled={loading} className="mt-3 w-full">
            {loading ? 'Working…' : 'Shorten'}
          </Button>
          {shortUrl && (
            <div className="mt-3 truncate text-xs">
              <code className="font-mono">{shortUrl}</code>
            </div>
          )}
        </form>
      </section>

      {error && (
        <p className="mt-6 rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-600 dark:text-red-400">
          {error}
        </p>
      )}

      <section className="mt-10">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-medium">Recent links</h2>
          <Link
            href="/dashboard/links"
            className="text-sm text-muted underline underline-offset-4"
          >
            View all
          </Link>
        </div>

        <div className="mt-4 overflow-hidden rounded-xl border border-[var(--border)]">
          <table className="w-full text-left text-sm">
            <thead className="border-b border-[var(--border)] bg-[var(--surface-2)]">
              <tr>
                <th className="px-3 py-2 font-medium">Code</th>
                <th className="px-3 py-2 font-medium">Long URL</th>
                <th className="px-3 py-2 font-medium">Clicks</th>
                <th className="px-3 py-2 font-medium">Status</th>
              </tr>
            </thead>
            <tbody>
              {!apiKey && (
                <tr>
                  <td colSpan={4} className="px-3 py-8 text-center text-muted">
                    Loading API key…
                  </td>
                </tr>
              )}
              {apiKey && links.length === 0 && (
                <tr>
                  <td colSpan={4} className="px-3 py-8 text-center text-muted">
                    No links yet. Use Quick shorten above.
                  </td>
                </tr>
              )}
              {links.map((l) => (
                <tr key={l.id} className="border-t border-[var(--border)]">
                  <td className="px-3 py-2 font-mono text-xs">
                    {l.short_code}
                  </td>
                  <td className="max-w-[24rem] truncate px-3 py-2 text-muted">
                    {l.long_url}
                  </td>
                  <td className="px-3 py-2 tabular-nums">{l.clicks}</td>
                  <td className="px-3 py-2">
                    <span
                      className={
                        l.is_active ? 'text-[var(--success)]' : 'text-muted'
                      }
                    >
                      {l.is_active ? 'Active' : 'Off'}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  );
}
