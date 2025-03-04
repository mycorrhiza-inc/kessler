import axios from "axios";
import { getRuntimeEnv } from "../env_variables_hydration_script";

export const GetConversationInformation = async (conversation_id: string) => {
  // get the conversation information from the database
  const runtimeConfig = getRuntimeEnv();
  const conversation = await axios.get(
    `${runtimeConfig.public_api_url}/v2/public/conversations/${conversation_id}`,
  );
  return conversation.data;
};
