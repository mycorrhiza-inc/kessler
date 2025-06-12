export const CLIENT_API_URL = process.env.NEXT_PUBLIC_KESSLER_API_URL || "https://api.kessler.xyz"

export const SERVER_API_URL = process.env.NEXT_PUBLIC_INTERNAL_KESSLER_API_URL || "http://backend-server:4041"


export const getContextualAPIUrl = (): string => {
  if (typeof (window) === "undefined") {
    return SERVER_API_URL
  }
  return CLIENT_API_URL
}
