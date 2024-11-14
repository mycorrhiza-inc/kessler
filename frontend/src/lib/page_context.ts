export type PageContext = {
  state?: string;
  slug: string[];
  final_identifier?: string;
};

export const getStateDisplayName = (state?: string) => {
  switch (state) {
    case "ny":
      return "New York State";
    case "ca":
      return "California";
    case "co":
      return "Colorado";
    case "":
      return "All States";
    case undefined:
      return "All States";
  }

  return "Unknown/Unsupported";
};
