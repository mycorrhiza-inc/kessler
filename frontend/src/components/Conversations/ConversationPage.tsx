import ConversationComponent from "@/components/Conversations/ConversationComponent";

import { FilterField, InheritedFilterValues } from "@/lib/filters";
import { PageContext } from "@/lib/page_context";
import axios from "axios";
import { BreadcrumbValues } from "../SitemapUtils";
import MarkdownRenderer from "../MarkdownRenderer";
import PageContainer from "../Page/PageContainer";
import { internalAPIURL } from "@/lib/env_variables";

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
const NYConversationDescription = ({ conversation }: { conversation: any }) => {
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
export const ConversationPage = async ({
  breadcrumbs,
}: {
  breadcrumbs: BreadcrumbValues;
}) => {
  const conversation_id =
    breadcrumbs.breadcrumbs[breadcrumbs.breadcrumbs.length - 1].value;
  const inheritedFilters: InheritedFilterValues = conversation_id
    ? [{ filter: FilterField.MatchDocketId, value: conversation_id }]
    : [];

  // const url = `${apiURL}/v2/public/conversations/named-lookup/${conversation_id}`;
  // USE THE PROD URL SINCE LOCALHOST ISNT REACHABLE FROM SERVER COMPONENTS
  const url = `${internalAPIURL}/v2/public/conversations/named-lookup/${conversation_id}`;
  const conversation = await getConversationData(url);
  // The title of the page looks weird with the really long title, either shorten
  const displayTitle =
    conversation.title.length > 50
      ? conversation.title.slice(0, 50) + "..."
      : conversation.title;
  var new_breadcrumbs = breadcrumbs;
  new_breadcrumbs.breadcrumbs[breadcrumbs.breadcrumbs.length - 1].title =
    displayTitle;

  return (
    <PageContainer breadcrumbs={new_breadcrumbs}>
      {conversation_id && (
        <NYConversationDescription conversation={conversation} />
      )}
      <ConversationComponent inheritedFilters={inheritedFilters} />
    </PageContainer>
  );
};
