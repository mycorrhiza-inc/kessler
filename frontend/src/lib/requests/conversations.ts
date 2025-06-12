import axios from "axios";

import type { Conversation } from "@/lib/conversations";
import { getContextualAPIUrl } from "../env_variables";

export const GetConversationInformation = async (
  conversation_id: string,
): Promise<Conversation> => {
  // get the conversation information from the database
  const conversation = await axios.get(
    `${getContextualAPIUrl()}/v2/public/conversations/${conversation_id}`,
  );
  return conversation.data;
};
