"use client";
import { Conversation, NYConversation } from "@/lib/conversations";
import MarkdownRenderer from "../MarkdownRenderer";
import useSWRImmutable from "swr";
import axios from "axios";
import { apiURL } from "@/lib/env_variables";

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

const getConversationData = async (url: string) => {
  const response = await axios.get(
    url,
    // "http://api.kessler.xyz/v2/recent_updates",
  );
  if (response.status !== 200) {
    throw new Error("Error fetching data with status " + response.status);
  }
  console.log("organization data", response.data);
  const json_convo = response.data.Description;
  const convo = JSON.parse(json_convo);
  return convo;
};
// export const NYConversationDescription = ({ conversation }: { conversation: NYConversation }) => {
export const NYConversationDescription = ({
  docket_id,
}: {
  docket_id: string;
}) => {
  const url = `${apiURL}/v2/public/conversations/named-lookup/${docket_id}`;
  const { data, error, isLoading } = useSWRImmutable(url, getConversationData);
  const conversation = isLoading ? dummy : data;
  return (
    <div className="conversation-description">
      <h1 className="text-2xl font-bold">
        {conversation.title} <br />
      </h1>

      <table className="table-auto">
        <tbody>
          <tr>
            <td>Case Number:</td>
            <td> {conversation.docket_id}</td>
          </tr>
          <tr>
            <td>Title of Matter:</td>
            <td>
              <MarkdownRenderer>{conversation.title}</MarkdownRenderer>
            </td>
          </tr>
          <tr>
            <td>Company/Organization: </td>
            <td>{conversation.organization}</td>
          </tr>
          <tr>
            <td>Matter Type: </td>
            <td>{conversation.matter_type}</td>
          </tr>
          <tr>
            <td>Matter Subtype: </td>
            <td>{conversation.matter_subtype}</td>
          </tr>
          <tr>
            <td>Date Filed: </td>
            <td>{conversation.date_filed}</td>
          </tr>

          <tr></tr>
        </tbody>
      </table>
    </div>
  );
};

export default ConversationDescription;
