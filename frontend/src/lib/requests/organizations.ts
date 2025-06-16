import axios from "axios";
import { getContextualAPIUrl } from "../env_variables";

export type OrganizationInfo = any;

export const getOrganizationInfo = async (orgID: string) => {
  const response = await axios.get(
    // `${runtimeConfig.public_api_url}/public/organizations/${orgID}`,
    `${getContextualAPIUrl()}/public/organizations/${orgID}`,
  );
  console.log("organization data", response.data);
  return response.data;
};
