export type PageContext = {
  state?: string;
  slug: string[];
  final_identifier?: string;
};

export const getStateDisplayName = (state?: string) => {
  // switch (state) {
  //   case "ny":

  //   case "ca":
  //     return "California";
  //   case "co":
  //     return "Colorado";
  //   case "":
  //     return "All States";
  //   case undefined:
  //     return "All States";
  // }

  // return "Unknown/Unsupported";
  // TODO: have something better. we only do ny rn
  return "New York";
};
