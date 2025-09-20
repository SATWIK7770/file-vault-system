// src/api/auth.ts
import { postJSON, getJSON } from "./index";

export type SignupPayload = { username: string; email?: string; password: string };
export type LoginPayload = { username: string; password: string };

export async function apiSignup(p: SignupPayload) {
  return postJSON("/api/signup", p);
}

export async function apiLogin(p: LoginPayload) {
  return postJSON("/api/login", p);
}

export async function apiLogout() {
  return postJSON("/api/logout", {});
}

export async function apiMe() {
  return getJSON("/api/me");
}
