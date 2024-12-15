import axios from "axios";
import { internalAPIURL } from "../env_variables";

export type OrganizationInfo = any;

export const getOrganizationInfo = async (orgID: string) => {
  const response = await axios.get(
    // `${publicAPIURL}/v2/public/organizations/${orgID}`,
    `${internalAPIURL}/v2/public/organizations/${orgID}`,
    // "http://api.kessler.xyz/v2/recent_updates",
  );
  console.log("organization data", response.data);
  return response.data;
};
