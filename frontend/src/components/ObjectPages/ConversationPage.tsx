import { FilterField, InheritedFilterValues } from "@/lib/filters";
import axios from "axios";
import { BreadcrumbValues } from "../SitemapUtils";
import MarkdownRenderer from "../MarkdownRenderer";
import { getUniversalEnvConfig } from "@/lib/env_variables/env_variables";
import HeaderCard from "./HeaderCard";
import {
  GenericSearchInfo,
  GenericSearchType,
} from "@/lib/adapters/genericSearchCallback";
import SearchResultsServer from "../Search/SearchResultsServer";

const getConversationData = async (url: string) => {
  const response = await axios.get(url);
  if (response.status !== 200) {
    throw new Error(
      "Error fetching data at " + url + " with status " + response.status,
    );
  }
  console.log("organization data", response.data);
  const json_convo = response.data.Metadata;
  const convo = JSON.parse(json_convo);
  return convo;
};

const NYConversationDescription = ({ conversation }: { conversation: any }) => {
  return (
    <HeaderCard title={conversation.title}>
      <table className="table-auto">
        <colgroup>
          <col width="20%" />
          <col width="80%" />
        </colgroup>
        <tbody>
          <tr>
            <td>Case Number:</td>
            <td>
              {conversation.docket_id + "         "}
              <a
                href={`https://documents.dps.ny.gov/public/MatterManagement/CaseMaster.aspx?MatterCaseNo=${conversation.docket_id}`}
                className="btn btn-secondary btn-xs"
              >
                Browse on New York State Website
              </a>
            </td>
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
    </HeaderCard>
  );
};

export const generateConversationInfo = async (convoNamedID: string) => {
  const inheritedFilters: InheritedFilterValues = convoNamedID
    ? [{ filter: FilterField.MatchDocketId, value: convoNamedID }]
    : [];

  const url = `${getUniversalEnvConfig().internal_api_url}/v2/public/conversations/named-lookup/${convoNamedID}`;
  const conversation = await getConversationData(url);
  // The title of the page looks weird with the really long title, either shorten
  const displayTitle =
    conversation.title.length > 50
      ? conversation.title.slice(0, 50) + "..."
      : conversation.title;

  const breadcrumbs: BreadcrumbValues = {
    state: "",
    breadcrumbs: [
      { value: "dockets", title: "Dockets" },
      { value: convoNamedID, title: displayTitle },
    ],
  };
  return {
    inheritedFilters: inheritedFilters,
    conversation: conversation,
    breadcrumbs: breadcrumbs,
    displayTitle: displayTitle,
  };
};

export const ConversationPage = async ({
  conversation,
  inheritedFilters,
  breadcrumbs,
}: {
  conversation: any;
  inheritedFilters: InheritedFilterValues;
  breadcrumbs: BreadcrumbValues;
}) => {
  const searchInfo: GenericSearchInfo = {
    search_type: GenericSearchType.Filling,
    query: "",
  };
  return (
    <>
      <NYConversationDescription conversation={conversation} />
      <SearchResultsServer searchInfo={searchInfo} />
    </>
  );
};
