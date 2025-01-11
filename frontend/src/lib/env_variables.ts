// This url is for hitting the kessler api from the client,
export const publicAPIURL =
  process.env.NEXT_PUBLIC_KESSLER_API_URL || "https://api.kessler.xyz";
export const isLocalMode = publicAPIURL.indexOf("localhost") !== -1;
// This is the url for hitting the api, for internal use by the nextjs server runtime, primarially for SSR pages.
// This defaults to the prod url for now, since setting it to localhost doesnt work because the api isnt accessible from the localhost of
// the nextjs docker container.
// TODO: Figure out why replacing this with http://backend-go:4041 doesnt work.
// which would improve security and prevent a network route outside the system for SSR.
// It does work for docker compose, k8s hates it though so it defaults to the api rn.
export const internalAPIURL =
  process.env.INTERNAL_KESSLER_API_URL || "https://api.kessler.xyz";

// export const internalAPIURL = "https://api.kessler.xyz";
