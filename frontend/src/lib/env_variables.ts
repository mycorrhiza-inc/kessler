import getConfig from "next/config";

const { publicRuntimeConfig, serverRuntimeConfig } = getConfig();

// This url is for hitting the kessler api from the client
export const publicAPIURL =
  publicRuntimeConfig.NEXT_PUBLIC_KESSLER_API_URL || "https://api.kessler.xyz";
export const isLocalMode = publicAPIURL.indexOf("localhost") !== -1;

// This is the url for hitting the api, for internal use by the nextjs server runtime
export const internalAPIURL =
  serverRuntimeConfig.INTERNAL_KESSLER_API_URL || "https://api.kessler.xyz";
// This url is for hitting the kessler api fVrom the client,
