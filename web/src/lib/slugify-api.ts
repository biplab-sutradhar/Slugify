export type Link = {
  id: string;
  short_code: string;
  long_url: string;
  is_active: boolean;
  clicks: number;
  created_at: string;
};

export type ApiKey = {
  id: string;
  key: string;
  name: string;
  scope: string;
  is_active: boolean;
  created_at: string;
  usage: number;
};

const BASE = "/backend";

function headers(apiKey: string): HeadersInit {
  return {
    "Content-Type": "application/json",
    "X-API-Key": apiKey,
  };
}

async function parseError(res: Response): Promise<string> {
  try {
    const j = (await res.json()) as { error?: string };
    return j.error ?? res.statusText;
  } catch {
    return res.statusText;
  }
}

/**
 * Cleans up sloppy URL input and returns a normalized URL string,
 * or null if the input doesn't represent a real URL.
 *
 * Handles:
 *  - leading/trailing whitespace
 *  - wrapping <…>, "…", '…', `…`
 *  - missing scheme (defaults to https://)
 *  - zero-width / BOM characters from rich-text copy-paste
 */
export function normalizeUrl(input: string): string | null {
  if (!input) return null;

  let s = input.trim();
  // Strip invisible characters that often come from copy-paste
  s = s.replace(/[\\u200B-\\u200D\\uFEFF]/g, "");
  // Strip wrapping quotes / brackets one layer at a time
  s = s.replace(/^[<"'`(]+/, "").replace(/[>"'`)]+$/, "");
  s = s.trim();

  if (!s) return null;

  const withScheme = /^https?:\/\//i.test(s) ? s : `https://${s}`;

  try {
    const u = new URL(withScheme);
    if (!u.hostname || !u.hostname.includes(".")) return null;
    return u.toString();
  } catch {
    return null;
  }
}

export async function shorten(longUrl: string, apiKey: string) {
  const res = await fetch(`${BASE}/shorten`, {
    method: "POST",
    headers: headers(apiKey),
    body: JSON.stringify({ long_url: longUrl }),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return (await res.json()) as { short_url: string };
}

export async function listLinks(
  apiKey: string,
  opts: { limit?: number; offset?: number } = {}
) {
  const q = new URLSearchParams();
  if (opts.limit != null) q.set("limit", String(opts.limit));
  if (opts.offset != null) q.set("offset", String(opts.offset));
  const res = await fetch(`${BASE}/links?${q.toString()}`, {
    headers: headers(apiKey),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return (await res.json()) as Link[];
}

export async function patchLinkActive(
  id: string,
  isActive: boolean,
  apiKey: string
) {
  const res = await fetch(`${BASE}/links/${id}`, {
    method: "PATCH",
    headers: headers(apiKey),
    body: JSON.stringify({ is_active: isActive }),
  });
  if (!res.ok) throw new Error(await parseError(res));
}

export async function deleteLink(id: string, apiKey: string) {
  const res = await fetch(`${BASE}/links/${id}`, {
    method: "DELETE",
    headers: headers(apiKey),
  });
  if (!res.ok) throw new Error(await parseError(res));
}

export async function listApiKeys(apiKey: string) {
  const res = await fetch(`${BASE}/keys`, { headers: headers(apiKey) });
  if (!res.ok) throw new Error(await parseError(res));
  const data = (await res.json()) as ApiKey[] | null;
  return data ?? [];
}

export async function deleteApiKey(id: string, apiKey: string) {
  const res = await fetch(`${BASE}/keys/${id}`, {
    method: "DELETE",
    headers: headers(apiKey),
  });
  if (!res.ok) throw new Error(await parseError(res));
}

export async function createApiKey(
  body: { name: string; scope: string },
  apiKey: string
) {
  const res = await fetch(`${BASE}/keys`, {
    method: "POST",
    headers: headers(apiKey),
    body: JSON.stringify(body),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return (await res.json()) as ApiKey;
}