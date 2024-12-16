"use client";
import axios from "axios";
import Link from "next/link";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import { publicAPIURL } from "@/lib/env_variables";
import InfiniteScroll from "react-infinite-scroll-component";
import { useEffect, useState } from "react";
import LoadingSpinnerTimeout from "../styled-components/LoadingSpinnerTimeout";

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

const ConversationTableHeaderless = ({
  convoList,
}: {
  convoList: ConversationTableSchema[];
}) => {
  return (
    <tbody>
      {convoList.map((convo: ConversationTableSchema) => {
        var description = null;
        const description_string = convo.Description;
        // console.log(description_string);
        try {
          description = JSON.parse(description_string);
          // console.log(description);
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
            className="border-base-300 hover:bg-base-200 transition duration-500 ease-out"
          >
            <td colSpan={7} className="p-0">
              <Link href={`/dockets/${convo.DocketID}`} className="flex w-full">
                <div className="w-[40%] px-4 py-3">{convo.Name}</div>
                <div className="w-[10%] px-4 py-3">{convo.DocketID}</div>
                <div className="w-[10%] px-4 py-3">{convo.DocumentCount}</div>
                <div className="w-[10%] px-4 py-3">{matter_type}</div>
                <div className="w-[10%] px-4 py-3">{matter_subtype}</div>
                <div className="w-[10%] px-4 py-3">{organization}</div>
                <div className="w-[10%] px-4 py-3">{date_filed}</div>
              </Link>
            </td>
          </tr>
        );
      })}
    </tbody>
  );
};

const ConversationTableHeaders = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  return (
    <table className="table table-pin-rows">
      {/* disable pinned rows due to the top row overlaying the filter sidebar */}
      <thead>
        <tr>
          <td className="w-[40%]">Name</td>
          <td className="w-[10%]">ID</td>
          <td className="w-[10%]">Document Count</td>
          <td className="w-[10%]">Matter Type</td>
          <td className="w-[10%]">Matter Subtype</td>
          <td className="w-[10%]">Organization</td>
          <td className="w-[10%]">Date Filed</td>
        </tr>
      </thead>
      {children}
    </table>
  );
};

const ConversationTable = ({
  convoList,
}: {
  convoList: ConversationTableSchema[];
}) => {
  return (
    <ConversationTableHeaders>
      <ConversationTableHeaderless convoList={convoList} />
    </ConversationTableHeaders>
  );
};

const ConversationTableInfiniteScroll = () => {
  const [tableData, setTableData] = useState<ConversationTableSchema[]>([]);
  const defaultPageSize = 40;
  const [page, setPage] = useState(0);

  const getPageResults = async (page: number, limit: number) => {
    const offset = page * limit;
    const result = await conversationsListGet(
      `${publicAPIURL}/v2/public/conversations/list?limit=${limit}&offset=${offset}`,
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
      <ConversationTableHeaders>
        <ConversationTableHeaderless convoList={tableData} />
      </ConversationTableHeaders>
    </InfiniteScroll>
  );
};

export default ConversationTableInfiniteScroll;
