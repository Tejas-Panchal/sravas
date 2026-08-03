import { createContext, useState, useEffect, useCallback } from "react";
import type { ReactNode } from "react";
import type { User } from "../types/index.ts";
import * as authApi from "../lib/auth.ts";

interface AuthState {
  user: User | null;
  loading: boolean;
  login: (identifier: string, password: string) => Promise<void>;
  register: (formData: FormData) => Promise<void>;
  logout: () => Promise<void>;
}

export const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const u = await authApi.getCurrentUser();
        if (!cancelled) setUser(u);
      } catch {
        try {
          await authApi.refreshToken();
          if (!cancelled) {
            const u = await authApi.getCurrentUser();
            setUser(u);
          }
        } catch {
          if (!cancelled) setUser(null);
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();

    return () => { cancelled = true; };
  }, []);

  const login = useCallback(async (identifier: string, password: string) => {
    const u = await authApi.login(identifier, password);
    setUser(u);
  }, []);

  const register = useCallback(async (formData: FormData) => {
    await authApi.register(formData);
  }, []);

  const logout = useCallback(async () => {
    await authApi.logout();
    setUser(null);
  }, []);

  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}
