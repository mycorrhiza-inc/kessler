export const stateFromHeaders = (headers: Headers) => {
  const host = headers.get("host") || "";
  const hostsplits = host.split(".");
  return hostsplits.length > 1 ? hostsplits[0] : undefined;
};
