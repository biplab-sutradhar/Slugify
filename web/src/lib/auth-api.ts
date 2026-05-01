export type User = {
  id: string;
  email: string;
  name: string;
  created_at: string;
};

export type AuthResponse = {
  token: string;
  user: User;
  api_key?: string;
};

const BASE = "/auth-api";

async function parseError(res: Response): Promise<string> {
  try {
    const j = (await res.json()) as { error?: string };
    return j.error ?? res.statusText;
  } catch {
    return res.statusText;
  }
}

export async function register(body: {
  email: string;
  password: string;
  name?: string;
}): Promise<AuthResponse> {
  const res = await fetch(`${BASE}/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function login(body: {
  email: string;
  password: string;
}): Promise<AuthResponse> {
  const res = await fetch(`${BASE}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function me(token: string): Promise<User> {
  const res = await fetch(`${BASE}/me`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function mintApiKey(token: string): Promise<{ api_key: string }> {
  const res = await fetch(`${BASE}/api-key`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ name: "Default" }),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export type ApiKey = {
  id: string;
  key: string;
  name: string;
  scope: string;
  is_active: boolean;
  created_at: string;
  usage: number;
};

function jwtHeaders(token: string): HeadersInit {
  return {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };
}

export async function listMyKeys(token: string): Promise<ApiKey[]> {
  const res = await fetch(`${BASE}/keys`, { headers: jwtHeaders(token) });
  if (!res.ok) throw new Error(await parseError(res));
  const data = (await res.json()) as ApiKey[] | null;
  return data ?? [];
}

export async function createMyKey(
  token: string,
  body: { name: string; scope: string }
): Promise<ApiKey> {
  const res = await fetch(`${BASE}/keys`, {
    method: "POST",
    headers: jwtHeaders(token),
    body: JSON.stringify(body),
  });
  if (!res.ok) throw new Error(await parseError(res));
  return res.json();
}

export async function deleteMyKey(token: string, id: string): Promise<void> {
  const res = await fetch(`${BASE}/keys/${id}`, {
    method: "DELETE",
    headers: jwtHeaders(token),
  });
  if (!res.ok) throw new Error(await parseError(res));
}