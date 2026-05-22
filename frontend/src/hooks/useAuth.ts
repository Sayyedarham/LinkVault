// src/hooks/useAuth.ts
import { useState, useCallback } from "react";
import { login as apiLogin, register as apiRegister } from "../api/auth";
import type { User } from "../types";

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(
    () => localStorage.getItem("token")
  );

  const signin = useCallback(async (email: string, password: string) => {
    const { data } = await apiLogin(email, password);
    localStorage.setItem("token", data.token);
    setToken(data.token);
    setUser(data.user);
  }, []);

  const signup = useCallback(async (email: string, password: string) => {
    const { data } = await apiRegister(email, password);
    localStorage.setItem("token", data.token);
    setToken(data.token);
    setUser(data.user);
  }, []);

  const signout = useCallback(() => {
    localStorage.removeItem("token");
    setToken(null);
    setUser(null);
  }, []);

  return { user, token, signin, signup, signout, isAuthed: !!token };
}