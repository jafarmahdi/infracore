import { create } from "zustand";
import { persist } from "zustand/middleware";

type Theme = "light" | "dark";
type Locale = "en" | "ar";

interface UIState {
  theme: Theme;
  locale: Locale;
  sidebarCollapsed: boolean;
  mobileSidebarOpen: boolean;
  toggleTheme: () => void;
  setLocale: (locale: Locale) => void;
  toggleSidebar: () => void;
  setMobileSidebarOpen: (open: boolean) => void;
}

export const useUIStore = create<UIState>()(
  persist(
    (set) => ({
      theme: "light",
      locale: "en",
      sidebarCollapsed: false,
      mobileSidebarOpen: false,
      toggleTheme: () =>
        set((state) => ({ theme: state.theme === "light" ? "dark" : "light" })),
      setLocale: (locale) => set({ locale }),
      toggleSidebar: () =>
        set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),
      setMobileSidebarOpen: (open) => set({ mobileSidebarOpen: open }),
    }),
    {
      name: "infracore-ui",
      partialize: ({ theme, locale, sidebarCollapsed }) => ({
        theme,
        locale,
        sidebarCollapsed,
      }),
    },
  ),
);
