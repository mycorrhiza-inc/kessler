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
  docket_id: "24-E-0138",
  matter_type: "Petition",
  matter_subtype:
    "Certificate of Public Convenience and Necessity - Electric Generation",
  title:
    "Petition of Bear Ridge Solar, LLC, for a Certificate of Public Convenience\n      and Necessity, Pursuant to Public Service Law Section 68, and for an Order\n      Granting Lightened Regulation.",
  organization: "Bear Ridge Solar, LLC",
  date_filed: "03/04/2024",
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
