"use client";

import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  deleteLink,
  listLinks,
  normalizeUrl,
  patchLinkActive,
  shorten,
  type Link,
} from "@/lib/slugify-api";
import { useAuth } from "@/lib/auth-context";

export default function LinksPage() {
  const { apiKey } = useAuth();
  const [links, setLinks] = useState<Link[]>([]);
  const [longUrl, setLongUrl] = useState("");
  const [shortUrl, setShortUrl] = useState<string | null>(null);
  const [offset, setOffset] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const limit = 20;

  const refresh = useCallback(async () => {
    if (!apiKey) {
      setLinks([]);
      return;
    }
    setError(null);
    setLoading(true);
    try {
      const data = await listLinks(apiKey, { limit, offset });
      setLinks(data);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load links');
    } finally {
      setLoading(false);
    }
  }, [apiKey, offset, limit]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

    const onShorten = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setShortUrl(null);

    if (!apiKey) {
      setError("Your API key is still being provisioned. Try again in a moment.");
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
      setLongUrl("");
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Shorten failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mx-auto max-w-5xl px-6 py-10">
      <h1 className="text-2xl font-semibold tracking-tight">Links</h1>

      <form
        onSubmit={onShorten}
        className="mt-6 flex flex-col gap-3 sm:flex-row sm:items-center"
      >
        <Input
          type="text"
          inputMode="url"
          autoComplete="off"
          spellCheck={false}
          required
          placeholder="example.com or <https://example.com/path>"
          value={longUrl}
          onChange={(e) => setLongUrl(e.target.value)}
        />
        <Button type="submit" disabled={loading}>
          {loading ? "Working…" : "Shorten"}
        </Button>
      </form>

      {shortUrl && (
        <div className="mt-4 rounded-lg border border-[var(--border)] bg-[var(--surface)] p-3 text-sm">
          <code className="font-mono">{shortUrl}</code>
          <button
            type="button"
            onClick={() => navigator.clipboard.writeText(shortUrl)}
            className="ml-3 underline underline-offset-4"
          >
            Copy
          </button>
        </div>
      )}

      {error && (
        <p className="mt-4 rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-600 dark:text-red-400">
          {error}
        </p>
      )}

      <div className="mt-8 overflow-hidden rounded-xl border border-[var(--border)]">
        <table className="w-full text-left text-sm">
          <thead className="border-b border-[var(--border)] bg-[var(--surface-2)]">
            <tr>
              <th className="px-3 py-2 font-medium">Code</th>
              <th className="px-3 py-2 font-medium">Long URL</th>
              <th className="px-3 py-2 font-medium">Clicks</th>
              <th className="px-3 py-2 font-medium">Active</th>
              <th className="px-3 py-2 font-medium" />
            </tr>
          </thead>
          <tbody>
            {!apiKey && (
              <tr>
                <td colSpan={5} className="px-3 py-8 text-center text-muted">
                  Loading API key…
                </td>
              </tr>
            )}
            {apiKey && links.length === 0 && !loading && (
              <tr>
                <td colSpan={5} className="px-3 py-8 text-center text-muted">
                  No links yet.
                </td>
              </tr>
            )}
            {links.map((link) => (
              <tr key={link.id} className="border-t border-[var(--border)]">
                <td className="px-3 py-2 font-mono text-xs">{link.short_code}</td>
                <td className="max-w-[24rem] truncate px-3 py-2 text-muted">{link.long_url}</td>
                <td className="px-3 py-2 tabular-nums">{link.clicks}</td>
                <td className="px-3 py-2">
                  <button
                    type="button"
                    className="text-xs underline underline-offset-4"
                    onClick={async () => {
                      if (!apiKey) return;
                      try {
                        await patchLinkActive(link.id, !link.is_active, apiKey);
                        await refresh();
                      } catch (e) {
                        setError(e instanceof Error ? e.message : "Update failed");
                      }
                    }}
                  >
                    {link.is_active ? "On" : "Off"}
                  </button>
                </td>
                <td className="px-3 py-2 text-right">
                  <button
                    type="button"
                    className="text-xs text-[var(--danger)] underline underline-offset-4"
                    onClick={async () => {
                      if (!apiKey) return;
                      if (!confirm("Delete this link?")) return;
                      try {
                        await deleteLink(link.id, apiKey);
                        await refresh();
                      } catch (e) {
                        setError(e instanceof Error ? e.message : "Delete failed");
                      }
                    }}
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="mt-4 flex justify-between text-sm text-muted">
        <button
          type="button"
          disabled={offset === 0}
          onClick={() => setOffset((o) => Math.max(0, o - limit))}
          className="underline underline-offset-4 disabled:opacity-30"
        >
          Previous
        </button>
        <button
          type="button"
          disabled={links.length < limit}
          onClick={() => setOffset((o) => o + limit)}
          className="underline underline-offset-4 disabled:opacity-30"
        >
          Next
        </button>
      </div>
    </div>
  );
}