"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { listMyKeys, me, mintApiKey, type User } from "./auth-api";

type AuthState = {
  user: User | null;
  token: string | null;
  apiKey: string | null;
  loading: boolean;
  provisioning: boolean;
  provisioningError: string | null;
  setSession: (token: string, user: User, apiKey?: string) => void;
  setApiKey: (key: string) => void;
  ensureApiKey: () => Promise<string | null>;
  logout: () => void;
};

const STORAGE_TOKEN = "slugify_token";
const STORAGE_USER = "slugify_user";
const STORAGE_API_KEY = "slugify_api_key";

const AuthContext = createContext<AuthState | null>(null);

function isAuthError(err: unknown): boolean {
  if (!(err instanceof Error)) return false;
  const m = err.message.toLowerCase();
  return (
    m.includes("invalid token") ||
    m.includes("missing bearer") ||
    m.includes("unauthorized") ||
    m.includes("user not identified")
  );
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [apiKey, setApiKeyState] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [provisioning, setProvisioning] = useState(false);
  const [provisioningError, setProvisioningError] = useState<string | null>(null);

  useEffect(() => {
    const t = localStorage.getItem(STORAGE_TOKEN);
    const uRaw = localStorage.getItem(STORAGE_USER);
    const k = localStorage.getItem(STORAGE_API_KEY);

    if (k) setApiKeyState(k);

    if (t && uRaw) {
      setToken(t);
      try {
        setUser(JSON.parse(uRaw) as User);
      } catch {
        // bad cache; will be refreshed below
      }

      me(t)
        .then((fresh) => {
          setUser(fresh);
          localStorage.setItem(STORAGE_USER, JSON.stringify(fresh));
        })
        .catch((err) => {
          if (isAuthError(err)) {
            localStorage.removeItem(STORAGE_TOKEN);
            localStorage.removeItem(STORAGE_USER);
            localStorage.removeItem(STORAGE_API_KEY);
            setToken(null);
            setUser(null);
            setApiKeyState(null);
          }
        })
        .finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, []);

  const persistApiKey = useCallback((k: string) => {
    localStorage.setItem(STORAGE_API_KEY, k);
    setApiKeyState(k);
  }, []);

  const setSession = useCallback(
    (t: string, u: User, k?: string) => {
      localStorage.setItem(STORAGE_TOKEN, t);
      localStorage.setItem(STORAGE_USER, JSON.stringify(u));
      setToken(t);
      setUser(u);
      if (k) persistApiKey(k);
    },
    [persistApiKey]
  );

  const ensureApiKey = useCallback(async (): Promise<string | null> => {
    if (apiKey) return apiKey;
    if (!token) return null;

    setProvisioning(true);
    setProvisioningError(null);
    try {
      // 1) Re-use an existing key for this user if one exists.
      try {
        const existing = await listMyKeys(token);
        if (existing.length > 0) {
          persistApiKey(existing[0].key);
          return existing[0].key;
        }
      } catch {
        // listing failed — fall through to mint
      }

      // 2) Otherwise mint a fresh one.
      const { api_key } = await mintApiKey(token);
      persistApiKey(api_key);
      return api_key;
    } catch (err) {
      const msg = err instanceof Error ? err.message : "Could not provision an API key";
      setProvisioningError(msg);
      return null;
    } finally {
      setProvisioning(false);
    }
  }, [apiKey, token, persistApiKey]);

  const logout = useCallback(() => {
    localStorage.removeItem(STORAGE_TOKEN);
    localStorage.removeItem(STORAGE_USER);
    localStorage.removeItem(STORAGE_API_KEY);
    setToken(null);
    setUser(null);
    setApiKeyState(null);
    setProvisioningError(null);
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        apiKey,
        loading,
        provisioning,
        provisioningError,
        setSession,
        setApiKey: persistApiKey,
        ensureApiKey,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside AuthProvider");
  return ctx;
}