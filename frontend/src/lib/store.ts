import { create } from "zustand";

interface KesslerState {
  isLoggedIn: boolean;
  enableExperimentalFeatures: boolean;
  setIsLoggedIn: (isLoggedIn: boolean) => void;
  setEnableExperimentalFeatures: (enableExperimentalFeatures: boolean) => void;
}

export const useKesslerStore = create<KesslerState>()((set) => ({
  enableExperimentalFeatures: false,
  isLoggedIn: false,
  setIsLoggedIn: (isLoggedIn: boolean) => set({ isLoggedIn }),
  setEnableExperimentalFeatures: (enableExperimentalFeatures: boolean) =>
    set({ enableExperimentalFeatures }),
}));
