"use client";
import axios from "axios";
import InfiniteScroll from "react-infinite-scroll-component";
import { useEffect, useState } from "react";
import LoadingSpinnerTimeout from "../../styled-components/LoadingSpinnerTimeout";
import { useRouter } from "next/navigation";
import { getRuntimeEnv } from "@/lib/env_variables_hydration_script";
import { TextPill } from "@/components/Tables/TextPills";
import clsx from "clsx";
import { TableStyle } from "@/components/styles/Table";
import InfiniteScrollPlus from "@/components/InfiniteScroll/InfiniteScroll";

export type ConversationSearchSchema = {
  query?: string;
  industry_type?: string;
  date_from?: string;
  date_to?: string;
};

const InstanitateConversationSearchSchema = (
  search_schema?: ConversationSearchSchema,
): ConversationSearchSchema => {
  return {
    ...Object.fromEntries(
      Object.keys(search_schema || {}).map((key) => [
        key,
        search_schema?.[key as keyof ConversationSearchSchema] || "",
      ]),
    ),
  } as ConversationSearchSchema;
};
const conversationSearchGet = async (
  searchData: ConversationSearchSchema,
  url: string,
) => {
  const cleanData = (response: any) => {
    // console.log(response.data);
    const return_data: any[] = response.data;
    if (return_data.length == 0 || return_data == undefined) {
      return [];
    }
    return return_data;
  };
  const result = await axios
    .post(url, searchData)
    .then((res) => {
      if (res.status >= 400) {
        throw new Error(`Request failed with status code ${res.status}`);
      }
      return cleanData(res);
    })
    .catch((error) => {
      throw error;
    });
  return result as ConversationTableSchema[];
};

type ConversationTableSchema = {
  docket_gov_id: string;
  state: string;
  name: string;
  description: string;
  matter_type: string;
  industry_type: string;
  metadata: string;
  extra: string;
  documents_count: number;
  date_published: string;
  id: string;
};

const ConversationTable = ({
  truncate,
  convoList,
  state,
}: {
  truncate?: boolean;
  state?: boolean;
  convoList: ConversationTableSchema[];
}) => {
  const router = useRouter();
  return (
    <table className={TableStyle}>
      <thead>
        <td className="w-[30%]">Name</td>
        <td className="w-[10%]">Date Published</td>
        <td className="w-[10%]">ID</td>
        <td className="w-[10%]">Matter Type</td>
        {state && <td className="w-[10%]">State</td>}
        <td className="w-[10%]">Industry</td>
        {!truncate && (
          <>
            <td className="w-[10%]">Description</td>
            <td className="w-[10%]">Document Count</td>
          </>
        )}
      </thead>
      <tbody>
        {convoList.map((convo: ConversationTableSchema) => {
          // const formattedDate = new Date(
          //   convo.date_published,
          // ).toLocaleDateString();
          // I apologize for my sins, we were under constraints

          const hackDate = JSON.parse(convo.metadata)["date_filed"];
          const hackMatterType = JSON.parse(convo.metadata)["matter_type"];

          return (
            <tr
              key={convo.docket_gov_id}
              className="border-base-300 hover:bg-base-200 transition duration-500 ease-out cursor-pointer"
              onClick={() => {
                router.push(`/dockets/${convo.docket_gov_id}`);
              }}
            >
              <td className="w-[60%] px-4 py-3">{convo.name}</td>
              <td className="w-[10%] px-4 py-3">{hackDate}</td>
              <td className="w-[10%] px-4 py-3">{convo.docket_gov_id}</td>
              <td className="w-[10%] px-4 py-3">{hackMatterType}</td>
              {state && <td className="w-[10%] px-4 py-3">{convo.state}</td>}
              <td className="w-[10%] px-4 py-3">
                {convo.industry_type && <TextPill text={convo.industry_type} />}
              </td>
              {!truncate && (
                <>
                  <td className="w-[10%] px-4 py-3">{convo.description}</td>
                  <td className="w-[10%] px-4 py-3">{convo.documents_count}</td>
                </>
              )}
            </tr>
          );
        })}
      </tbody>
    </table>
  );
};

const ConversationTableInfiniteScroll = ({
  lookup_data,
  truncate,
}: {
  lookup_data?: ConversationSearchSchema;
  truncate?: boolean;
}) => {
  const [tableData, setTableData] = useState<ConversationTableSchema[]>([]);
  const defaultPageSize = 40;
  const [page, setPage] = useState(0);
  const [hasMore, setHasMore] = useState(true);

  const getPageResults = async (page: number, limit: number) => {
    const offset = page * limit;
    const runtimeConfig = getRuntimeEnv();
    const searchData = InstanitateConversationSearchSchema(lookup_data);
    const url = `${runtimeConfig.public_api_url}/v2/search/conversation?offset=${offset}&limit=${limit}`;
    const result = await conversationSearchGet(searchData, url);
    setTableData((prev) => [...prev, ...result]);
    if (result.length < limit) {
      setHasMore(false);
    }
    return result;
  };
  const getMore = async () => {
    await getPageResults(page, defaultPageSize);
    setPage((prev) => prev + 1);
  };
  const getInitialData = async () => {
    setHasMore(true);
    setTableData([]);
    const numPageFetch = 3;
    await getPageResults(0, defaultPageSize * numPageFetch);
    setPage(numPageFetch);
  };
  return (
    <InfiniteScrollPlus
      hasMore={hasMore}
      dataLength={tableData.length}
      getMore={getMore}
      loadInitial={getInitialData}
      reloadOnChangeObj={lookup_data}
    >
      <ConversationTable convoList={tableData} truncate={truncate} />
    </InfiniteScrollPlus>
  );
};

export default ConversationTableInfiniteScroll;
