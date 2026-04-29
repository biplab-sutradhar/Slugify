"use client";

import { useCallback, useEffect, useState } from "react";
import {
  createApiKey,
  deleteApiKey,
  listApiKeys,
  type ApiKey,
} from "@/lib/slugify-api";
import { useApiKey } from "@/lib/use-api-key";

export default function KeysPage() {
  const { apiKey, setApiKey, save } = useApiKey();
  const [keys, setKeys] = useState<ApiKey[]>([]);
  const [name, setName] = useState("");
  const [scope, setScope] = useState("default");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [justCreated, setJustCreated] = useState<ApiKey | null>(null);

  const refresh = useCallback(async () => {
    if (!apiKey.trim()) return;
    setError(null);
    setLoading(true);
    try {
      const data = await listApiKeys(apiKey.trim());
      setKeys(data);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load keys");
    } finally {
      setLoading(false);
    }
  }, [apiKey]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  const onCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setJustCreated(null);
    if (!apiKey.trim()) {
      setError("API key required.");
      return;
    }
    setLoading(true);
    try {
      const created = await createApiKey(
        { name: name.trim(), scope: scope.trim() || "default" },
        apiKey.trim()
      );
      setJustCreated(created);
      setName("");
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Create failed");
    } finally {
      setLoading(false);
    }
  };

  const onRevoke = async (id: string) => {
    if (!confirm("Revoke this key? This cannot be undone.")) return;
    try {
      await deleteApiKey(id, apiKey.trim());
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Revoke failed");
    }
  };

  return (
    <main className="mx-auto max-w-2xl px-6 py-16">
      <p className="text-xs uppercase tracking-[0.2em] text-neutral-500">
        Slugify
      </p>
      <h1 className="mt-3 text-3xl font-medium tracking-tight">API keys</h1>
      <p className="mt-2 text-sm text-neutral-500">
        Manage keys that can call the Slugify API.
      </p>

      <section className="mt-10 space-y-4">
        <label className="block text-sm font-medium">Admin key</label>
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
          onClick={() => save(apiKey)}
          className="text-sm text-neutral-600 underline underline-offset-4 hover:text-foreground dark:text-neutral-400"
        >
          Save key in this tab
        </button>
      </section>

      <form onSubmit={onCreate} className="mt-10 space-y-4">
        <h2 className="text-lg font-medium">Create key</h2>
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <input
            required
            type="text"
            placeholder="Name (e.g. mobile app)"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="rounded-lg border border-black/10 bg-transparent px-3 py-2 text-sm outline-none ring-foreground/15 focus:ring-2 dark:border-white/15"
          />
          <input
            type="text"
            placeholder="Scope (default)"
            value={scope}
            onChange={(e) => setScope(e.target.value)}
            className="rounded-lg border border-black/10 bg-transparent px-3 py-2 text-sm outline-none ring-foreground/15 focus:ring-2 dark:border-white/15"
          />
        </div>
        <button
          type="submit"
          disabled={loading}
          className="rounded-full bg-foreground px-5 py-2 text-sm font-medium text-background transition hover:opacity-90 disabled:opacity-40"
        >
          {loading ? "Working…" : "Create key"}
        </button>
      </form>

      {error && (
        <p className="mt-6 rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-600 dark:text-red-400">
          {error}
        </p>
      )}

      {justCreated && (
        <div className="mt-8 rounded-xl border border-emerald-500/30 bg-emerald-500/10 p-4">
          <p className="text-xs font-medium text-emerald-700 dark:text-emerald-300">
            Save this key now — it won’t be shown again
          </p>
          <div className="mt-2 flex flex-wrap items-center gap-3">
            <code className="break-all font-mono text-sm">
              {justCreated.key}
            </code>
            <button
              type="button"
              onClick={() => navigator.clipboard.writeText(justCreated.key)}
              className="text-sm underline underline-offset-4"
            >
              Copy
            </button>
          </div>
        </div>
      )}

      <section className="mt-14">
        <h2 className="text-lg font-medium">All keys</h2>
        <div className="mt-6 overflow-hidden rounded-xl border border-black/10 dark:border-white/10">
          <table className="w-full text-left text-sm">
            <thead className="border-b border-black/10 bg-black/[0.03] dark:border-white/10 dark:bg-white/[0.04]">
              <tr>
                <th className="px-3 py-2 font-medium">Name</th>
                <th className="px-3 py-2 font-medium">Scope</th>
                <th className="px-3 py-2 font-medium">Usage</th>
                <th className="px-3 py-2 font-medium">Active</th>
                <th className="px-3 py-2 font-medium" />
              </tr>
            </thead>
            <tbody>
              {keys.length === 0 && !loading && (
                <tr>
                  <td
                    colSpan={5}
                    className="px-3 py-8 text-center text-neutral-500"
                  >
                    No keys yet.
                  </td>
                </tr>
              )}
              {keys.map((k) => (
                <tr
                  key={k.id}
                  className="border-t border-black/5 dark:border-white/10"
                >
                  <td className="max-w-[12rem] truncate px-3 py-2">
                    {k.name || "—"}
                  </td>
                  <td className="px-3 py-2">{k.scope}</td>
                  <td className="px-3 py-2 tabular-nums">{k.usage}</td>
                  <td className="px-3 py-2">
                    <span
                      className={
                        k.is_active
                          ? "text-emerald-600 dark:text-emerald-400"
                          : "text-neutral-500"
                      }
                    >
                      {k.is_active ? "On" : "Off"}
                    </span>
                  </td>
                  <td className="px-3 py-2 text-right">
                    <button
                      type="button"
                      onClick={() => void onRevoke(k.id)}
                      className="text-xs text-red-600 underline underline-offset-4 dark:text-red-400"
                    >
                      Revoke
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
    </main>
  );
}