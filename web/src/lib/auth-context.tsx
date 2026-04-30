"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import { me, mintApiKey, type User } from "./auth-api";

type AuthState = {
  user: User | null;
  token: string | null;
  apiKey: string | null;
  loading: boolean;
  setSession: (token: string, user: User, apiKey?: string) => void;
  ensureApiKey: () => Promise<string | null>;
  logout: () => void;
};

const STORAGE_TOKEN = "slugify_token";
const STORAGE_USER = "slugify_user";
const STORAGE_API_KEY = "slugify_api_key";

const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [apiKey, setApiKey] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const t = localStorage.getItem(STORAGE_TOKEN);
    const u = localStorage.getItem(STORAGE_USER);
    const k = localStorage.getItem(STORAGE_API_KEY);
    if (k) setApiKey(k);

    if (t && u) {
      setToken(t);
      try {
        setUser(JSON.parse(u) as User);
      } catch {
        /* ignore */
      }
      me(t)
        .then((fresh) => {
          setUser(fresh);
          localStorage.setItem(STORAGE_USER, JSON.stringify(fresh));
        })
        .catch(() => {
          localStorage.removeItem(STORAGE_TOKEN);
          localStorage.removeItem(STORAGE_USER);
          localStorage.removeItem(STORAGE_API_KEY);
          setToken(null);
          setUser(null);
          setApiKey(null);
        })
        .finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, []);

  const setSession = useCallback(
    (t: string, u: User, k?: string) => {
      localStorage.setItem(STORAGE_TOKEN, t);
      localStorage.setItem(STORAGE_USER, JSON.stringify(u));
      setToken(t);
      setUser(u);
      if (k) {
        localStorage.setItem(STORAGE_API_KEY, k);
        setApiKey(k);
      }
    },
    []
  );

  const ensureApiKey = useCallback(async () => {
    if (apiKey) return apiKey;
    if (!token) return null;
    const { api_key } = await mintApiKey(token);
    localStorage.setItem(STORAGE_API_KEY, api_key);
    setApiKey(api_key);
    return api_key;
  }, [apiKey, token]);

  const logout = useCallback(() => {
    localStorage.removeItem(STORAGE_TOKEN);
    localStorage.removeItem(STORAGE_USER);
    localStorage.removeItem(STORAGE_API_KEY);
    setToken(null);
    setUser(null);
    setApiKey(null);
  }, []);

  return (
    <AuthContext.Provider
      value={{ user, token, apiKey, loading, setSession, ensureApiKey, logout }}
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