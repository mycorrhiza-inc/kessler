import { create } from "zustand";

interface KesslerState {
  isLoggedIn: boolean;
  experimentalFeaturesEnabled: boolean;
  setIsLoggedIn: (isLoggedIn: boolean) => void;
  setExperimentalFeaturesEnabled: (enableExperimentalFeatures: boolean) => void;
}

export const useKesslerStore = create<KesslerState>()((set) => ({
  experimentalFeaturesEnabled: false,
  isLoggedIn: false,
  setIsLoggedIn: (isLoggedIn: boolean) => set({ isLoggedIn }),
  setExperimentalFeaturesEnabled: (experimentalFeaturesEnabled: boolean) =>
    set({ experimentalFeaturesEnabled }),
}));
