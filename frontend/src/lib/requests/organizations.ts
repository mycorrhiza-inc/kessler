import axios from "axios";
import { apiURL, prodAPIURL } from "../env_variables";

export const getOrganizationInfo = async (orgID: string) => {
  const response = await axios.get(
    // `${apiURL}/v2/public/organizations/${orgID}`,
    `${prodAPIURL}/v2/public/organizations/${orgID}`,
    // "http://api.kessler.xyz/v2/recent_updates",
  );
  console.log("organization data", response.data);
  return response.data;
};
