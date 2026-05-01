"use client";

import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  type ApiKey,
  createMyKey,
  deleteMyKey,
  listMyKeys,
} from "@/lib/auth-api";
import { useAuth } from "@/lib/auth-context";

export default function KeysPage() {
  const { token, apiKey } = useAuth();
  const [keys, setKeys] = useState<ApiKey[]>([]);
  const [name, setName] = useState("");
  const [scope, setScope] = useState("default");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [justCreated, setJustCreated] = useState<ApiKey | null>(null);

  const refresh = useCallback(async () => {
    if (!token) return;
    setError(null);
    setLoading(true);
    try {
      setKeys(await listMyKeys(token));
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load keys");
    } finally {
      setLoading(false);
    }
  }, [token]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  const onCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) {
      setError("Not signed in.");
      return;
    }
    setError(null);
    setJustCreated(null);
    setLoading(true);
    try {
      const created = await createMyKey(token, {
        name: name.trim() || "API key",
        scope: scope.trim() || "default",
      });
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
    if (!token) return;
    if (!confirm("Revoke this key? This cannot be undone.")) return;
    try {
      await deleteMyKey(token, id);
      await refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Revoke failed");
    }
  };

  return (
    <div className="mx-auto max-w-5xl px-6 py-10">
      <p className="text-xs uppercase tracking-[0.2em] text-muted">Settings</p>
      <h1 className="mt-1 text-2xl font-semibold tracking-tight">API keys</h1>
      <p className="mt-2 text-sm text-muted">
        Use these to call the Slugify API from external scripts or services.
      </p>

      {apiKey && (
        <section className="mt-8 rounded-xl border border-[var(--border)] bg-[var(--surface)] p-5">
          <p className="text-xs uppercase tracking-[0.15em] text-muted">
            Default key (auto-provisioned)
          </p>
          <div className="mt-2 flex flex-wrap items-center gap-3">
            <code className="break-all font-mono text-sm">{apiKey}</code>
            <button
              type="button"
              onClick={() => navigator.clipboard.writeText(apiKey)}
              className="text-xs underline underline-offset-4"
            >
              Copy
            </button>
          </div>
        </section>
      )}

      <form onSubmit={onCreate} className="mt-10 space-y-4">
        <h2 className="text-lg font-medium">Create a new key</h2>
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <Input
            type="text"
            placeholder="Name (e.g. mobile app)"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <Input
            type="text"
            placeholder="Scope (default)"
            value={scope}
            onChange={(e) => setScope(e.target.value)}
          />
        </div>
        <Button type="submit" disabled={loading || !token}>
          {loading ? "Working…" : "Create key"}
        </Button>
      </form>

      {error && (
        <p className="mt-6 rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-600 dark:text-red-400">
          {error}
        </p>
      )}

      {justCreated && (
        <div className="mt-8 rounded-xl border border-emerald-500/30 bg-emerald-500/10 p-4">
          <p className="text-xs font-medium text-emerald-700 dark:text-emerald-300">
            Save this key now — you'll only see the full value here.
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

      <section className="mt-12">
        <h2 className="text-lg font-medium">All keys</h2>
        <div className="mt-4 overflow-hidden rounded-xl border border-[var(--border)]">
          <table className="w-full text-left text-sm">
            <thead className="border-b border-[var(--border)] bg-[var(--surface-2)]">
              <tr>
                <th className="px-3 py-2 font-medium">Name</th>
                <th className="px-3 py-2 font-medium">Scope</th>
                <th className="px-3 py-2 font-medium">Usage</th>
                <th className="px-3 py-2 font-medium">Active</th>
                <th className="px-3 py-2 font-medium" />
              </tr>
            </thead>
            <tbody>
              {!loading && keys.length === 0 && (
                <tr>
                  <td colSpan={5} className="px-3 py-8 text-center text-muted">
                    No keys yet.
                  </td>
                </tr>
              )}
              {keys.map((k) => (
                <tr key={k.id} className="border-t border-[var(--border)]">
                  <td className="max-w-[14rem] truncate px-3 py-2">
                    {k.name || "—"}
                  </td>
                  <td className="px-3 py-2">{k.scope}</td>
                  <td className="px-3 py-2 tabular-nums">{k.usage}</td>
                  <td className="px-3 py-2">
                    <span
                      className={
                        k.is_active
                          ? "text-[var(--success)]"
                          : "text-muted"
                      }
                    >
                      {k.is_active ? "Active" : "Off"}
                    </span>
                  </td>
                  <td className="px-3 py-2 text-right">
                    <button
                      type="button"
                      onClick={() => void onRevoke(k.id)}
                      className="text-xs text-[var(--danger)] underline underline-offset-4"
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
    </div>
  );
}