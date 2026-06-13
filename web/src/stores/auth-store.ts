import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { User } from "@/types/auth";

interface AuthState {
  accessToken: string | null;
  user: User | null;
  isAuthenticated: boolean;
  setSession: (user: User, token: string) => void;
  clearSession: () => void;
  hasPermission: (perm: string) => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      accessToken: null,
      user: null,
      isAuthenticated: false,

      setSession: (user, token) => {
        set({ user, accessToken: token, isAuthenticated: true });
      },

      clearSession: () => {
        set({ user: null, accessToken: null, isAuthenticated: false });
      },

      hasPermission: (perm) => {
        const { user } = get();
        if (!user) return false;
        if (user.is_superuser) return true;
        return user.permissions?.includes(perm) ?? false;
      },
    }),
    {
      name: "infracore-auth",
      partialize: (state) => ({
        user: state.user,
        // accessToken intentionally excluded — re-auth on page reload
        isAuthenticated: false,
      }),
    },
  ),
);
