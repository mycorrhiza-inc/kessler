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
    // FIXME: FIGURE OUT WHY IT ISNT RETURNING NEW YORK FOR THESE BAD CASES
    case "":
      return "New York";
    case undefined:
      return "New York";
  }

  return "New York";
};
