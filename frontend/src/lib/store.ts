import { create } from "zustand";
import { isLocalMode, runtimeConfig.public_api_url } from "./env_variables";

interface KesslerState {
  isLoggedIn: boolean;
  experimentalFeaturesEnabled: boolean;
  setIsLoggedIn: (isLoggedIn: boolean) => void;
  setExperimentalFeaturesEnabled: (enableExperimentalFeatures: boolean) => void;
}

export const useKesslerStore = create<KesslerState>()((set) => ({
  experimentalFeaturesEnabled: isLocalMode,
  isLoggedIn: false,
  setIsLoggedIn: (isLoggedIn: boolean) => set({ isLoggedIn }),
  setExperimentalFeaturesEnabled: (experimentalFeaturesEnabled: boolean) =>
    set({ experimentalFeaturesEnabled }),
}));
