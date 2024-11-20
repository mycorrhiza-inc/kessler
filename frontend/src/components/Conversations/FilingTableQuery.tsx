"use client";
import { Suspense, memo } from "react";
import useSWRImmutable from "swr";
import FilingTable from "./FilingTable";
import { QueryDataFile } from "@/lib/filters";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import getSearchResults from "@/lib/requests/search";

interface FilingTableQueryProps {
  queryData: QueryDataFile;
  scroll?: boolean;
}
const FilingTableQuery = memo(
  ({ queryData, scroll }: FilingTableQueryProps) => {
    const { data, error, isLoading } = useSWRImmutable(
      queryData,
      getSearchResults,
    );
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
    return <FilingTable filings={filings} scroll={scroll} />;
  },
);

export default FilingTableQuery;
