"use client";
import axios from "axios";
import InfiniteScroll from "react-infinite-scroll-component";
import { useEffect, useState } from "react";
import LoadingSpinnerTimeout from "../styled-components/LoadingSpinnerTimeout";
import { useRouter } from "next/navigation";
import { getRuntimeEnv } from "@/lib/env_variables_hydration_script";

const conversationsListGet = async (url: string) => {
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
  Name: string;
  DocketID: string;
  DocumentCount: number;
  State: string;
  Description: string;
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
    <table className="table table-pin-rows">
      {/* disable pinned rows due to the top row overlaying the filter sidebar */}
      <thead>
        <tr>
          <td className="w-[40%]">Name</td>
          <td className="w-[10%]">Date Filed</td>
          <td className="w-[10%]">ID</td>
          <td className="w-[10%]">Matter Type</td>
          <td className="w-[10%]">Matter Subtype</td>
          {!truncate && (
            <>
              <td className="w-[10%]">Organization</td>
              <td className="w-[10%]">Document Count</td>
            </>
          )}
        </tr>
      </thead>

      <tbody>
        {convoList.map((convo: ConversationTableSchema) => {
          var description = null;
          const description_string = convo.Description;
          console.log(description_string);
          try {
            description = JSON.parse(description_string);
            console.log(description);
          } catch (e) {
            console.log("Error parsing JSON", e);
          }
          const matter_type = description?.matter_type || "unknown";
          const matter_subtype = description?.matter_subtype || "unknown";
          const organization = description?.organization || "unknown";
          const date_filed = description?.date_filed || "unknown";

          return (
            <tr
              key={convo.DocketID}
              className="border-base-300 hover:bg-base-200 transition duration-500 ease-out cursor-pointer"
              onClick={() => {
                router.push(`/dockets/${convo.DocketID}`);
              }}
            >
              <td className="w-[60%] px-4 py-3">{convo.Name}</td>
              <td className="w-[10%] px-4 py-3">{date_filed}</td>
              <td className="w-[10%] px-4 py-3">{convo.DocketID}</td>
              <td className="w-[10%] px-4 py-3">{matter_type}</td>
              <td className="w-[10%] px-4 py-3">{matter_subtype}</td>
              {!truncate && (
                <>
                  <td className="w-[10%] px-4 py-3">{organization}</td>
                  <td className="w-[10%] px-4 py-3">{convo.DocumentCount}</td>
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
    const result = await conversationsListGet(
      `${runtimeConfig.public_api_url}/v2/public/conversations/list?limit=${limit}&offset=${offset}`,
    );
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
          timeoutSeconds={3}
          loadingText="Loading Conversations"
        />
      }
    >
      <ConversationTable convoList={tableData} truncate={truncate} />
    </InfiniteScroll>
  );
};

export default ConversationTableInfiniteScroll;
