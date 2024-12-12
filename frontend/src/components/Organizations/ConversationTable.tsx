"use client";
import axios from "axios";
import Link from "next/link";
import useSWRImmutable from "swr/immutable";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import { publicAPIURL } from "@/lib/env_variables";
import InfiniteScroll from "react-infinite-scroll-component";
import { useEffect, useState } from "react";

const conversationsListGet = async (url: string) => {
  const cleanData = (response: any) => {
    console.log(response.data);
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
};

const ConversationTable = ({
  convoList,
}: {
  convoList: ConversationTableSchema[];
}) => {
  return (
    <table className="table table-pin-rows">
      {/* disable pinned rows due to the top row overlaying the filter sidebar */}
      <thead>
        <tr>
          <td className="w-[60%]">Name</td>
          <td className="w-[20%]">ID</td>
          <td className="w-[10%]">Document Count</td>
          <td className="w-[10%]">State</td>
        </tr>
      </thead>
      <tbody>
        {convoList.map((convo: any) => (
          <tr
            key={convo.DocketID}
            className="border-base-300 hover:bg-base-200 transition duration-500 ease-out"
          >
            <td colSpan={4} className="p-0">
              <Link href={`/dockets/${convo.DocketID}`} className="flex w-full">
                <div className="w-[60%] px-4 py-3">{convo.Name}</div>
                <div className="w-[20%] px-4 py-3">{convo.DocketID}</div>
                <div className="w-[10%] px-4 py-3">{convo.DocumentCount}</div>
                <div className="w-[10%] px-4 py-3">{convo.State}</div>
              </Link>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
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
    <>
      <InfiniteScroll
        dataLength={tableData.length}
        hasMore={true}
        next={getMore}
        loader={<LoadingSpinner loadingText="Loading Conversations" />}
      >
        <ConversationTable convoList={tableData} />
      </InfiniteScroll>
    </>
  );
};

export default ConversationTableInfiniteScroll;
