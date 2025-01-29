"use client";
import axios from "axios";
import InfiniteScroll from "react-infinite-scroll-component";
import { useEffect, useState } from "react";
import LoadingSpinnerTimeout from "../styled-components/LoadingSpinnerTimeout";
import { useRouter } from "next/navigation";
import { getRuntimeEnv } from "@/lib/env_variables_hydration_script";
import { ClassNames } from "@emotion/react";

type ConversationSearchSchema = {
  query: string;
  industry_type: string;
  date_from: string;
  date_to: string;
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
  const result = await axios.get(url).then((res) => cleanData(res));
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
}: {
  truncate?: boolean;
  convoList: ConversationTableSchema[];
}) => {
  const router = useRouter();
  return (
    <table className="table z-1">
      {/* disable pinned rows due to the top row overlaying the filter sidebar */}
      <thead>
        <tr>
          <td className="w-[30%]">Name</td>
          <td className="w-[10%]">Date Published</td>
          <td className="w-[10%]">ID</td>
          <td className="w-[10%]">Matter Type</td>
          <td className="w-[10%]">State</td>
          <td className="w-[10%]">Industry</td>
          {!truncate && (
            <>
              <td className="w-[10%]">Description</td>
              <td className="w-[10%]">Document Count</td>
            </>
          )}
        </tr>
      </thead>

      <tbody>
        {convoList.map((convo: ConversationTableSchema) => {
          const formattedDate = new Date(
            convo.date_published,
          ).toLocaleDateString();

          return (
            <tr
              key={convo.docket_gov_id}
              className="border-base-300 hover:bg-base-200 transition duration-500 ease-out cursor-pointer"
              onClick={() => {
                router.push(`/dockets/${convo.docket_gov_id}`);
              }}
            >
              <td className="w-[60%] px-4 py-3">{convo.name}</td>
              <td className="w-[10%] px-4 py-3">{formattedDate}</td>
              <td className="w-[10%] px-4 py-3">{convo.docket_gov_id}</td>
              <td className="w-[10%] px-4 py-3">{convo.matter_type}</td>
              <td className="w-[10%] px-4 py-3">{convo.state}</td>
              <td className="w-[10%] px-4 py-3">{convo.industry_type}</td>

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
  truncate,
}: {
  truncate?: boolean;
}) => {
  const [tableData, setTableData] = useState<ConversationTableSchema[]>([]);
  const defaultPageSize = 40;
  const [page, setPage] = useState(0);

  const getPageResults = async (page: number, limit: number) => {
    const offset = page * limit;
    const runtimeConfig = getRuntimeEnv();
    const searchData = {
      query: "",
      industry_type: "",
      date_from: "",
      date_to: "",
    };
    const url = `${runtimeConfig.public_api_url}/v2/search/conversation?offset=${offset}&limit=${limit}`;
    const result = await conversationSearchGet(searchData, url);
    return result;
  };
  const getMore = async () => {
    const result = await getPageResults(page, defaultPageSize);
    setTableData((prev) => [...prev, ...result]);
    setPage((prev) => prev + 1);
  };
  const getInitialData = async () => {
    const numPageFetch = 3;
    const result = await getPageResults(0, defaultPageSize * numPageFetch);
    setTableData(result);
    setPage(numPageFetch);
  };
  useEffect(() => {
    getInitialData();
  }, []);
  return (
    <InfiniteScroll
      dataLength={tableData.length}
      hasMore={true}
      next={getMore}
      loader={
        <LoadingSpinnerTimeout
          timeoutSeconds={10}
          loadingText="Loading Conversations"
        />
      }
    >
      <ConversationTable convoList={tableData} truncate={truncate} />
    </InfiniteScroll>
  );
};

export default ConversationTableInfiniteScroll;
