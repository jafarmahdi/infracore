import type { LoginCredentials, LoginResponse } from "@/types/auth";
import { apiClient } from "./client";

export async function login(credentials: LoginCredentials): Promise<LoginResponse> {
  const { data } = await apiClient.post<LoginResponse>("/auth/login", credentials);
  return data;
}

export async function logout(): Promise<void> {
  await apiClient.post("/auth/logout").catch(() => {});
}

export async function getMe() {
  const { data } = await apiClient.get("/auth/me");
  return data;
}
