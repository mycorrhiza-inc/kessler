"use client";
import axios from "axios";
import Link from "next/link";
import LoadingSpinner from "../styled-components/LoadingSpinner";


import { publicAPIURL } from "@/lib/env_variables";
import { useEffect, useState } from "react";
import InfiniteScroll from "react-infinite-scroll-component";
import { queryStringFromPageMaxHits } from "@/lib/pagination";

const organizationsListGet = async (url: string) => {
  const cleanData = (response: any) => {
    console.log(response.data);
    const return_data: any[] = response.data;
    if (return_data.length == 0 || return_data == undefined) {
      return [];
    }
    return return_data;
    // TODO: Fix this validator code at some point
    // const valid_data = return_data.map(
    //   (item): OrganizationSchemaComplete =>
    //     OrganizationSchemaCompleteValidator.parse(return_data),
    // );
    // return valid_data;
  };
  const result = await axios.get(url).then((res) => cleanData(res));
  return result as OrganizationTableSchema[];
};

type OrganizationTableSchema = {
  Name: string;
  DocumentCount: number;
};

const OrganizationTable = ({
  orgList,
}: {
  orgList: OrganizationTableSchema[];
}) => {
  return (
    <table className="table table-pin-rows">
      <thead>
        <tr>
          <td className="w-[80%]">Name</td>
          <td className="w-[20%]">Documents Authored</td>
          {/* <td>Description</td> */}
        </tr>
      </thead>
      <tbody>
        {orgList.map((org: any) => (
          <tr
            key={org.DocketID}
            className="border-base-300 hover:bg-base-200 transition duration-500 ease-out"
          >
            <td colSpan={2} className="p-0">
              <Link href={`/orgs/${org.ID}`} className="flex w-full">
                <div className="w-[80%] px-4 py-3">{org.Name}</div>
                <div className="w-[20%] px-4 py-3">{org.DocumentCount}</div>
                {/* <div className="flex-1 px-4 py-3">{org.Description}</div> */}
              </Link>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

const OrganizationTableInfiniteScroll = () => {
  const defaultPageSize = 40;
  const [page, setPage] = useState(0);

  const [tableData, setTableData] = useState<OrganizationTableSchema[]>([]);
  const getPageResults = async (page: number, limit: number) => {
    const queryString = queryStringFromPageMaxHits(page, limit);
    const result = await organizationsListGet(
      `${publicAPIURL}/v2/public/organizations/list${queryString}`,
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
        <OrganizationTable orgList={tableData} />
      </InfiniteScroll>
    </>
  );
};

export default OrganizationTableInfiniteScroll;
