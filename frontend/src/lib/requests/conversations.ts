import axios from "axios";
import { getClientRuntimeEnv } from "../env_variables/env_variables_hydration_script";

import type { Conversation } from "@/lib/conversations";

export const GetConversationInformation = async (
  conversation_id: string,
): Promise<Conversation> => {
  // get the conversation information from the database
  const runtimeConfig = getClientRuntimeEnv();
  const conversation = await axios.get(
    `${runtimeConfig.public_api_url}/v2/public/conversations/${conversation_id}`,
  );
  return conversation.data;
};
