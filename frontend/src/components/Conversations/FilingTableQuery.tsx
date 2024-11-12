"use client";
import { Suspense, memo } from "react";
import useSWRImmutable from "swr";
import FilingTable from "./FilingTable";
import { QueryDataFile } from "@/lib/filters";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import getSearchResults from "@/lib/requests/search";

const FilingTableQueryRaw = memo(
  ({ queryData }: { queryData: QueryDataFile }) => {
    const { data, error, isLoading } = useSWRImmutable(
      queryData,
      getSearchResults,
    );
    console.log("data: ", data);
    if (isLoading) {
      return <LoadingSpinner />;
    }
    if (error) {
      return (
        <p>
          Encountered an error getting files from the server. <br />
          {String(error)}
        </p>
      );
    }
    const filings = data;
    if (filings == undefined) {
      return <p>Filings returned from server is undefined.</p>;
    }
    return <FilingTable filings={filings} />;
  },
);

const FilingTableQuery = ({ queryData }: { queryData: QueryDataFile }) => {
  return <FilingTableQueryRaw queryData={queryData} />;
};
export default FilingTableQuery;
