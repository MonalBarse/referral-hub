"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

import { api, tokenStore } from "@/lib/api";
import type { User } from "@/lib/types";

type AuthState = {
  user: User | null;
  token: string | null;
  ready: boolean;
  login: (email: string, name: string) => Promise<void>;
  logout: () => void;
};

const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [ready, setReady] = useState(false);

  // On first load, confirm any stored token is still valid before trusting it.
  useEffect(() => {
    let cancelled = false;
    const stored = tokenStore.get();
    const verify = stored ? api.me(stored) : Promise.resolve(null);

    verify
      .then((me) => {
        if (cancelled || !me || !stored) return;
        setUser(me);
        setToken(stored);
      })
      .catch(() => {
        tokenStore.clear();
      })
      .finally(() => {
        if (!cancelled) setReady(true);
      });

    return () => {
      cancelled = true;
    };
  }, []);

  const login = useCallback(async (email: string, name: string) => {
    const result = await api.login(email, name);
    tokenStore.set(result.token);
    setToken(result.token);
    setUser(result.user);
  }, []);

  const logout = useCallback(() => {
    tokenStore.clear();
    setToken(null);
    setUser(null);
  }, []);

  const value = useMemo(
    () => ({ user, token, ready, login, logout }),
    [user, token, ready, login, logout],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside AuthProvider");
  return ctx;
}
