import api from "./api.ts";
import type { User, ApiResponse } from "../types/index.ts";

export async function getCurrentUser(): Promise<User> {
  const res = await api.get<ApiResponse<User>>("/users/get-user");
  return res.data.data;
}

export async function refreshToken(): Promise<void> {
  await api.post("/users/refresh-token");
}

export async function login(identifier: string, password: string): Promise<User> {
  const res = await api.post<ApiResponse<{ user: User }>>("/users/login", {
    username: identifier,
    email: identifier,
    password,
  });
  return res.data.data.user;
}

export async function register(formData: FormData): Promise<void> {
  await api.post("/users/register", formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });
}

export async function logout(): Promise<void> {
  await api.post("/users/logout");
}
