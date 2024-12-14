"use client";
import { useState, useEffect } from "react";
import { Filing } from "@/lib/types/FilingTypes";
import { FilingTable } from "@/components/Tables/FilingTable";
import { getFilingMetadata, getRecentFilings } from "@/lib/requests/search";

import InfiniteScroll from "react-infinite-scroll-component";
import LoadingSpinnerTimeout from "../styled-components/LoadingSpinnerTimeout";

// TODO: Break out Recent Updates into its own component seperate from all of the homepage logic
export default function RecentUpdatesView() {
  const [isSearching, setIsSearching] = useState(false);
  const [filing_ids, setFilingIds] = useState<string[]>([]);
  const [filings, setFilings] = useState<Filing[]>([]);
  const [page, setPage] = useState(0);

  const getUpdates = async () => {
    setIsSearching(true);
    console.log("getting recent updates");
    const data = await getRecentFilings();
    console.log();

    const ids = data.map((item: any) => item.sourceID);
    console.log("ids", ids);
    setFilingIds(ids);
    setIsSearching(false);
  };

  // useEffect(() => {
  //   if (!filing_ids) {
  //     return;
  //   }
  //
  //   const fetchFilings = async () => {
  //     const newFilings = await Promise.all(
  //       filing_ids.map(async (id) => {
  //         const filing_data = await getFilingMetadata(id);
  //         console.log("new filings", filing_data);
  //         return filing_data;
  //       }),
  //     );
  //
  //     setFilings((previous: Filing[]): Filing[] => {
  //       const existingIds = new Set(previous.map((f: Filing) => f.id));
  //       const uniqueNewFilings = newFilings.filter(
  //         (f: Filing | null) => f !== null && !existingIds.has(f.id),
  //       ) as Filing[];
  //       console.log(" uniques: ", uniqueNewFilings);
  //       console.log("all data: ", [...previous, ...uniqueNewFilings]);
  //       return [...previous, ...uniqueNewFilings];
  //     });
  //   };
  //
  //   fetchFilings();
  // }, [filing_ids]);

  const getMore = async () => {
    setIsSearching(true);
    try {
      console.log("getting page ", page + 1);
      const new_filings = await getRecentFilings(page);
      setPage(page + 1);
      console.log(new_filings);
      if (filings.length > 0) {
        setFilings((old_filings: Filing[]) => [...old_filings, ...new_filings]);
      }
    } catch (error) {
      console.log(error);
    } finally {
      setIsSearching(false);
    }
  };

  useEffect(() => {
    getUpdates();
  }, []);

  return (
    <InfiniteScroll
      dataLength={filings.length}
      next={getMore}
      hasMore={true}
      loader={
        <div onClick={getMore}>
          <LoadingSpinnerTimeout
            loadingText="Loading Files"
            timeoutSeconds={3}
          />
        </div>
      }
    >
      <FilingTable filings={filings} DocketColumn />
    </InfiniteScroll>
  );
}
