"use client";
import axios from "axios";
import Link from "next/link";

import { useEffect, useState } from "react";
import { getClientRuntimeEnv } from "@/lib/env_variables/env_variables_hydration_script";
import { TableStyle } from "@/components/styles/Table";
import InfiniteScrollPlus from "@/components/InfiniteScroll/InfiniteScroll";
import { queryStringFromPagination } from "@/lib/types/new_search_types";

export interface OrganizationSearchSchema {
  query?: string;
}

const InstanitateOrganizationSearchSchema = (
  search_schema?: OrganizationSearchSchema,
): OrganizationSearchSchema => {
  return {
    ...Object.fromEntries(
      Object.keys(search_schema || {}).map((key) => [
        key,
        search_schema?.[key as keyof OrganizationSearchSchema] || "",
      ]),
    ),
  } as OrganizationSearchSchema;
};
const organizationsListGet = async (
  url: string,
  search_schema: OrganizationSearchSchema,
) => {
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
  const result = await axios.post(url, search_schema).then((res) => {
    if (res.status >= 400) {
      throw new Error(`Request failed with status code ${res.status}`);
    }
    return cleanData(res);
  });
  return result as OrganizationTableSchema[];
};

type OrganizationTableSchema = {
  name: string;
  id: string;
  aliases: string[];
  files_authored_count: number;
};

const OrganizationTable = ({
  orgList,
}: {
  orgList: OrganizationTableSchema[];
}) => {
  return (
    <table className={TableStyle}>
      <thead>
        <td className="w-[80%]">Name</td>
        <td className="w-[20%]">Documents Authored</td>
      </thead>
      <tbody>
        {orgList.map((org: any) => (
          <tr
            key={org.ID}
            className="border-base-300 hover:bg-base-200 transition duration-500 ease-out"
          >
            <td colSpan={2} className="p-0">
              <Link href={`/orgs/${org.id}`} className="flex w-full">
                <div className="w-[80%] px-4 py-3">{org.name}</div>
                <div className="w-[20%] px-4 py-3">
                  {org.files_authored_count}
                </div>
                {/* <div className="flex-1 px-4 py-3">{org.Description}</div> */}
              </Link>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

const OrganizationTableInfiniteScroll = ({
  lookup_data,
}: {
  lookup_data?: OrganizationSearchSchema;
}) => {
  const defaultPageSize = 40;
  const [page, setPage] = useState(0);
  const [hasMore, setHasMore] = useState(true);

  const [tableData, setTableData] = useState<OrganizationTableSchema[]>([]);
  const getPageResults = async (page: number, limit: number) => {
    const queryString = queryStringFromPagination({ page: page, limit: limit });
    const runtimeConfig = getClientRuntimeEnv();
    const searchData = InstanitateOrganizationSearchSchema(lookup_data);
    const url = `${runtimeConfig.public_api_url}/v2/search/organization${queryString}`;
    const result = await organizationsListGet(url, searchData);
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
    const numPageFetch = 3;
    setTableData([]);
    await getPageResults(0, defaultPageSize * numPageFetch);
    setPage(numPageFetch);
  };
  return (
    <>
      <InfiniteScrollPlus
        dataLength={tableData.length}
        hasMore={hasMore}
        getMore={getMore}
        loadInitial={getInitialData}
        reloadOnChange={0}
      >
        <OrganizationTable orgList={tableData} />
      </InfiniteScrollPlus>
    </>
  );
};

export default OrganizationTableInfiniteScroll;
