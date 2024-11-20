import axios from "axios";
import { apiURL } from "../env_variables";

export const getOrganizationInfo = async (orgID: string) => {
  const response = await axios.post(
    `${apiURL}v2/public/organizations/${orgID}`,
    // "http://api.kessler.xyz/v2/recent_updates",
  );
  console.log("recent data", response.data);
  if (response.data.length > 0) {
    return response.data;
  }
};
