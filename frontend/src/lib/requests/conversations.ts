import axios from "axios";
import { apiURL } from "../env_variables";

import { Conversation } from "@/lib/conversations";
import { convertLength } from "@mui/material/styles/cssUtils";

export const GetConversationInformation = async (
	conversation_id: string
	) => {
	  // get the conversation information from the database
  const conversation = await axios.get(`${apiURL}/v2/public/conversations/${conversation_id}`);
  return conversation.data;
}