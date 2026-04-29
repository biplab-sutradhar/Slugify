"use client";

import { useCallback, useEffect, useState } from "react";
import {
  deleteLink,
  listLinks,
  patchLinkActive,
  shorten,
  type Link,
} from "@/lib/slugify-api";

const STORAGE_KEY = "slugify_api_key";

export default function Home() {
  const [apiKey, setApiKey] = useState("");
  const [longUrl, setLongUrl] = useState("");
  const [shortUrl, setShortUrl] = useState<string | null>(null);
  const [links, setLinks] = useState<Link[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const limit = 20;
  const [offset, setOffset] = useState(0);

  useEffect(() => {
    const k = sessionStorage.getItem(STORAGE_KEY);
    if (k) setApiKey(k);
  }, []);

  const saveKey = () => {
    sessionStorage.setItem(STORAGE_KEY, apiKey.trim());
  };

  const refresh = useCallback(async () => {
    if (!apiKey.trim()) {
      setError("Add your API key first.");
      return;
    }
    setError(null);
    setLoading(true);
    try {
      const data = await listLinks(apiKey.trim(), { limit, offset });
      setLinks(data);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Request failed");
    } finally {
      setLoading(false);
    }
  }, [apiKey, offset]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  const onShorten = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setShortUrl(null);
    if (!apiKey.trim()) {
      setError("API key required.");
      return;
    }
    setLoading(true);
    try {
      saveKey();
      const r = await shorten(longUrl.trim(), apiKey.trim());
      setShortUrl(r.short_url);
      setLongUrl("");
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Shorten failed");
    } finally {
      setLoading(false);
    }
  };

  const copy = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
    } catch {
      setError("Could not copy to clipboard");
    }
  };

  return (
    <div className="min-h-screen bg-[var(--background)] text-[var(--foreground)]">
      <main className="mx-auto max-w-2xl px-6 py-16">
        <p className="text-xs uppercase tracking-[0.2em] text-neutral-500">
          Slugify
        </p>
        <h1 className="mt-3 text-3xl font-medium tracking-tight">
          Short links
        </h1>
        <p className="mt-2 text-sm text-neutral-500">
          Paste your API key once. Keys are stored in session storage only in
          this browser.
        </p>

        <section className="mt-10 space-y-4">
          <label className="block text-sm font-medium">API key</label>
          <input
            type="password"
            autoComplete="off"
            placeholder="xxxxxxxx"
            value={apiKey}
            onChange={(e) => setApiKey(e.target.value)}
            className="w-full rounded-lg border border-black/10 bg-transparent px-3 py-2 text-sm outline-none ring-foreground/15 focus:ring-2 dark:border-white/15"
          />
          <button
            type="button"
            onClick={saveKey}
            className="text-sm text-neutral-600 underline underline-offset-4 hover:text-foreground dark:text-neutral-400"
          >
            Save key in this tab
          </button>
        </section>

        <form onSubmit={onShorten} className="mt-10 space-y-4">
          <label className="block text-sm font-medium">Long URL</label>
          <input
            required
            type="url"
            placeholder="<https://example.com/very/long/path>"
            value={longUrl}
            onChange={(e) => setLongUrl(e.target.value)}
            className="w-full rounded-lg border border-black/10 bg-transparent px-3 py-2 text-sm outline-none ring-foreground/15 focus:ring-2 dark:border-white/15"
          />
          <button
            type="submit"
            disabled={loading}
            className="rounded-full bg-foreground px-5 py-2 text-sm font-medium text-background transition hover:opacity-90 disabled:opacity-40"
          >
            {loading ? "Working…" : "Shorten"}
          </button>
        </form>

        {error && (
          <p className="mt-6 rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-600 dark:text-red-400">
            {error}
          </p>
        )}

        {shortUrl && (
          <div className="mt-8 rounded-xl border border-black/10 p-4 dark:border-white/10">
            <p className="text-xs text-neutral-500">Short URL</p>
            <div className="mt-2 flex flex-wrap items-center gap-3">
              <code className="font-mono text-sm">{shortUrl}</code>
              <button
                type="button"
                onClick={() => void copy(shortUrl)}
                className="text-sm underline underline-offset-4"
              >
                Copy
              </button>
            </div>
          </div>
        )}

        <section className="mt-14">
          <div className="flex items-center justify-between gap-4">
            <h2 className="text-lg font-medium">Links</h2>
            <button
              type="button"
              onClick={() => void refresh()}
              className="text-sm text-neutral-600 underline underline-offset-4 dark:text-neutral-400"
            >
              Refresh
            </button>
          </div>

          <div className="mt-6 overflow-hidden rounded-xl border border-black/10 dark:border-white/10">
            <table className="w-full text-left text-sm">
              <thead className="border-b border-black/10 bg-black/[0.03] dark:border-white/10 dark:bg-white/[0.04]">
                <tr>
                  <th className="px-3 py-2 font-medium">Code</th>
                  <th className="px-3 py-2 font-medium">Clicks</th>
                  <th className="px-3 py-2 font-medium">Active</th>
                  <th className="px-3 py-2 font-medium" />
                </tr>
              </thead>
              <tbody>
                {links.length === 0 && !loading && (
                  <tr>
                    <td
                      colSpan={4}
                      className="px-3 py-8 text-center text-neutral-500"
                    >
                      No links yet.
                    </td>
                  </tr>
                )}
                {links.map((link) => (
                  <tr
                    key={link.id}
                    className="border-t border-black/5 dark:border-white/10"
                  >
                    <td className="max-w-[10rem] truncate px-3 py-2 font-mono text-xs">
                      {link.short_code}
                    </td>
                    <td className="px-3 py-2 tabular-nums">{link.clicks}</td>
                    <td className="px-3 py-2">
                      <button
                        type="button"
                        className="text-xs underline underline-offset-4"
                        onClick={async () => {
                          try {
                            await patchLinkActive(
                              link.id,
                              !link.is_active,
                              apiKey.trim()
                            );
                            await refresh();
                          } catch (e) {
                            setError(
                              e instanceof Error ? e.message : "Update failed"
                            );
                          }
                        }}
                      >
                        {link.is_active ? "On" : "Off"}
                      </button>
                    </td>
                    <td className="px-3 py-2 text-right">
                      <button
                        type="button"
                        className="text-xs text-red-600 underline underline-offset-4 dark:text-red-400"
                        onClick={async () => {
                          if (!confirm("Delete this link?")) return;
                          try {
                            await deleteLink(link.id, apiKey.trim());
                            await refresh();
                          } catch (e) {
                            setError(
                              e instanceof Error ? e.message : "Delete failed"
                            );
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

          <div className="mt-4 flex justify-between text-sm text-neutral-500">
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
        </section>
      </main>
    </div>
  );
}