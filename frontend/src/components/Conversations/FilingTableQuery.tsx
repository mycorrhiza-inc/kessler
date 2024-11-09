"use client";
import { Suspense, memo } from "react";
import useSWR from "swr";
import FilingTable from "./FilingTable";
import getSearchResults from "./searchResultGet";
import { QueryDataFile } from "@/lib/filters";
import LoadingSpinner from "../styled-components/LoadingSpinner";

const FilingTableQueryRaw = memo(
  ({ queryData }: { queryData: QueryDataFile }) => {
    const { data, error } = useSWR(queryData, getSearchResults, {
      suspense: true,
    });
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
  return (
    <Suspense
      fallback={<LoadingSpinner loadingText="Loading Search Results..." />}
    >
      <FilingTableQueryRaw queryData={queryData} />
    </Suspense>
  );
};
export default FilingTableQuery;
