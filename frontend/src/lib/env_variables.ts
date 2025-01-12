// This url is for hitting the kessler api from the client
export const publicAPIURL =
  process.env.NEXT_PUBLIC_KESSLER_API_URL || "https://api.kessler.xyz";
export const isLocalMode = publicAPIURL.indexOf("localhost") !== -1;

// This is the url for hitting the api, for internal use by the nextjs server runtime
export const internalAPIURL =
  process.env.INTERNAL_KESSLER_API_URL || "https://api.kessler.xyz";
// This url is for hitting the kessler api fVrom the client,
