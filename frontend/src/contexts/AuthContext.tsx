import React, { createContext, useContext, useEffect, useState, ReactNode } from "react";
import * as authApi from "../api/auth";

export type User = {
  id: number;
  username: string;
  email?: string;
  role: string;
};

type AuthContextType = {
  user: User | null;
  loading: boolean;
  login: (username: string, password: string) => Promise<void>;
  signup: (username: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refresh: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
};

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  
  const refresh = async () => {
    try {
      const data = await authApi.apiMe();
      const u = data && (data.user ?? data);
      setUser(u ?? null);
    } catch {
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    refresh(); // check session on mount
  }, []);

  // ---- login ----
  const login = async (username: string, password: string) => {
    await authApi.apiLogin({ username, password }); // sets cookie
    await refresh();
  };

  // ---- signup ----
  const signup = async (username: string, email: string, password: string) => {
    await authApi.apiSignup({ username, email, password }); // sets cookie
    await refresh();
  };

  // ---- logout ----
  const logout = async () => {
    try {
      await authApi.apiLogout(); // clears cookie
    } finally {
      setUser(null);
    }
  };

  const value = { user, loading, login, signup, logout, refresh };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
