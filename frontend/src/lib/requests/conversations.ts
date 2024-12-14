import axios from "axios";
import { publicAPIURL } from "../env_variables";


export const GetConversationInformation = async (conversation_id: string) => {
  // get the conversation information from the database
  const conversation = await axios.get(
    `${publicAPIURL}/v2/public/conversations/${conversation_id}`,
  );
  return conversation.data;
};
