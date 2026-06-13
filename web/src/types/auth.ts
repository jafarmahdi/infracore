export type UserRole = "admin" | "operator" | "viewer";

export interface User {
  id: string;
  tenant_id: string;
  tenant_slug: string;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  avatar_url?: string;
  is_superuser: boolean;
  roles: string[] | null;
  permissions: string[];
  // convenience alias used in UI
  name?: string;
  role?: UserRole;
  company?: string;
}

export interface LoginCredentials {
  tenant_slug: string;
  email: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
  user: User;
}

export interface RefreshResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
  user: User;
}
