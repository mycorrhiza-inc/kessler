"use client";
import { memo } from "react";
import useSWR from "swr";
import FilingTable from "./FilingTable";
import getSearchResults from "./searchResultGet";
import { QueryDataFile } from "@/lib/filters";

const FilingTableQuery = memo(({ queryData }: { queryData: QueryDataFile }) => {
  const { data, error } = useSWR(queryData, getSearchResults, {
    suspense: true,
  });
  const filings = data;
  if (filings == undefined) {
    return <p>Filings returned from server is undefined.</p>;
  }
  return <FilingTable filings={filings} />;
  // } catch (error) {
  //   return (
  //     <p>
  //       Encountered an error getting files from the server. <br />
  //       {String(error)}
  //     </p>
  //   );
  // }
});

export default FilingTableQuery;
