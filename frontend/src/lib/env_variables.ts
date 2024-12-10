// This url is for hitting the kessler api from the client,
export const publicAPIURL =
  process.env.NEXT_PUBLIC_KESSLER_API_URL || "https://api.kessler.xyz";
// This is the url for hitting the api, for internal use by the nextjs server runtime, primarially for SSR pages.
// This defaults to the prod url for now, since setting it to localhost doesnt work because the api isnt accessible from the localhost of
// the nextjs docker container.
// TODO: Figure out why replacing this with http://backend-go:4041 doesnt work.
// which would improve security and prevent a network route outside the system for SSR.
export const internalAPIURL =
  process.env.INTERNAL_KESSLER_API_URL || "http://backend-go:4041";

// export const internalAPIURL = "https://api.kessler.xyz";
