"use client";
import { Conversation, NYConversation } from "@/lib/conversations";
import MarkdownRenderer from "../MarkdownRenderer";
import useSWRImmutable from "swr";
import axios from "axios";

const ConversationDescription = ({
  conversation,
}: {
  conversation: Conversation;
}) => {
  return (
    <div className="conversation-description">
      <div className="conversation-description__title">{conversation.name}</div>
      <div className="conversation-description__last-message">
        {conversation.description}
      </div>
    </div>
  );
};

const dummy: NYConversation = {
  docket_id: "Loading...",
  matter_type: "Loading...",
  matter_subtype: "Loading...",
  title: "Loading...",
  organization: "Loading...",
  date_filed: "Loading...",
};

// export const NYConversationDescription = ({ conversation }: { conversation: NYConversation }) => {

export default ConversationDescription;
