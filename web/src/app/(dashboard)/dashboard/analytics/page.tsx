'use client';

import { useCallback, useEffect, useMemo, useState } from 'react';
import { listLinks, type Link } from '@/lib/slugify-api';
import { useAuth } from '@/lib/auth-context';

function BarRow({
  label,
  value,
  pct,
}: {
  label: string;
  value: number;
  pct: number;
}) {
  return (
    <div className="space-y-1">
      <div className="flex items-center justify-between gap-3 text-sm">
        <code className="truncate font-mono text-xs">{label}</code>
        <span className="tabular-nums text-muted">
          {value.toLocaleString()}
        </span>
      </div>
      <div className="h-2 overflow-hidden rounded-full bg-[var(--surface-2)]">
        <div
          className="h-full rounded-full bg-[var(--accent)]"
          style={{ width: `${Math.max(2, pct)}%` }}
        />
      </div>
    </div>
  );
}

export default function AnalyticsPage() {
  const { apiKey } = useAuth();
  const [links, setLinks] = useState<Link[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    if (!apiKey) return;
    setError(null);
    setLoading(true);
    try {
      // Pull a wide window for client-side aggregation
      const data = await listLinks(apiKey, { limit: 100, offset: 0 });
      setLinks(data);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load analytics');
    } finally {
      setLoading(false);
    }
  }, [apiKey]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  const stats = useMemo(() => {
    const total = links.length;
    const active = links.filter((l) => l.is_active).length;
    const totalClicks = links.reduce((acc, l) => acc + l.clicks, 0);
    const avg = total > 0 ? Math.round(totalClicks / total) : 0;
    const top = [...links].sort((a, b) => b.clicks - a.clicks).slice(0, 8);
    const topClicks = top[0]?.clicks ?? 0;

    const created = links.reduce<Record<string, number>>((acc, l) => {
      const day = l.created_at.slice(0, 10);
      acc[day] = (acc[day] ?? 0) + 1;
      return acc;
    }, {});
    const last7 = Array.from({ length: 7 }).map((_, i) => {
      const d = new Date();
      d.setDate(d.getDate() - (6 - i));
      const key = d.toISOString().slice(0, 10);
      return { day: key, count: created[key] ?? 0 };
    });
    const max7 = Math.max(1, ...last7.map((d) => d.count));

    return { total, active, totalClicks, avg, top, topClicks, last7, max7 };
  }, [links]);

  return (
    <div className="mx-auto max-w-5xl px-6 py-10">
      <p className="text-xs uppercase tracking-[0.2em] text-muted">Insights</p>
      <h1 className="mt-1 text-2xl font-semibold tracking-tight">Analytics</h1>
      <p className="mt-2 text-sm text-muted">
        A quick look at how your links are performing.
      </p>

      {error && (
        <p className="mt-6 rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-600 dark:text-red-400">
          {error}
        </p>
      )}

      <section className="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {[
          { label: 'Total links', value: stats.total },
          { label: 'Active', value: stats.active },
          { label: 'Total clicks', value: stats.totalClicks },
          { label: 'Avg clicks / link', value: stats.avg },
        ].map((s) => (
          <div
            key={s.label}
            className="rounded-xl border border-[var(--border)] bg-[var(--surface)] p-5"
          >
            <p className="text-xs text-muted">{s.label}</p>
            <p className="mt-2 text-3xl font-semibold tracking-tight">
              {s.value.toLocaleString()}
            </p>
          </div>
        ))}
      </section>

      <section className="mt-6 grid gap-4 lg:grid-cols-3">
        <div className="rounded-xl border border-[var(--border)] bg-[var(--surface)] p-5 lg:col-span-2">
          <p className="text-sm font-medium">Top links by clicks</p>
          <p className="text-xs text-muted">
            Top 8 of your most-used short codes.
          </p>
          <div className="mt-5 space-y-4">
            {!loading && stats.top.length === 0 && (
              <p className="text-sm text-muted">
                No data yet — create a link first.
              </p>
            )}
            {stats.top.map((l) => (
              <BarRow
                key={l.id}
                label={l.short_code}
                value={l.clicks}
                pct={(l.clicks / Math.max(1, stats.topClicks)) * 100}
              />
            ))}
          </div>
        </div>

        <div className="rounded-xl border border-[var(--border)] bg-[var(--surface)] p-5">
          <p className="text-sm font-medium">Links created (last 7 days)</p>
          <div className="mt-5 flex h-40 items-end gap-2">
            {stats.last7.map((d) => (
              <div
                key={d.day}
                className="flex flex-1 flex-col items-center gap-1"
              >
                <div
                  className="w-full rounded-t bg-[var(--accent)]/80"
                  style={{
                    height: `${(d.count / stats.max7) * 100}%`,
                    minHeight: d.count > 0 ? '8%' : '2%',
                  }}
                  title={`${d.day}: ${d.count}`}
                />
                <span className="text-[10px] text-muted">{d.day.slice(5)}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="mt-10">
        <h2 className="text-lg font-medium">All links by activity</h2>
        <div className="mt-4 overflow-hidden rounded-xl border border-[var(--border)]">
          <table className="w-full text-left text-sm">
            <thead className="border-b border-[var(--border)] bg-[var(--surface-2)]">
              <tr>
                <th className="px-3 py-2 font-medium">Code</th>
                <th className="px-3 py-2 font-medium">Long URL</th>
                <th className="px-3 py-2 font-medium">Clicks</th>
                <th className="px-3 py-2 font-medium">Created</th>
              </tr>
            </thead>
            <tbody>
              {!loading && links.length === 0 && (
                <tr>
                  <td colSpan={4} className="px-3 py-8 text-center text-muted">
                    No data yet.
                  </td>
                </tr>
              )}
              {[...links]
                .sort((a, b) => b.clicks - a.clicks)
                .map((l) => (
                  <tr key={l.id} className="border-t border-[var(--border)]">
                    <td className="px-3 py-2 font-mono text-xs">
                      {l.short_code}
                    </td>
                    <td className="max-w-[24rem] truncate px-3 py-2 text-muted">
                      {l.long_url}
                    </td>
                    <td className="px-3 py-2 tabular-nums">{l.clicks}</td>
                    <td className="px-3 py-2 text-muted">
                      {new Date(l.created_at).toLocaleDateString()}
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
