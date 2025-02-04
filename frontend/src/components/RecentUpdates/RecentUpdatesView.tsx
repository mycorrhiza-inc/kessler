"use client";
import { useState, useEffect } from "react";
import { Filing } from "@/lib/types/FilingTypes";
import { FilingTable } from "@/components/Tables/FilingTable";
import { getRecentFilings } from "@/lib/requests/search";

import InfiniteScroll from "react-infinite-scroll-component";
import LoadingSpinnerTimeout from "../styled-components/LoadingSpinnerTimeout";
import { fi } from "date-fns/locale";

// TODO: Break out Recent Updates into its own component seperate from all of the homepage logic
export default function RecentUpdatesView() {
  const [isSearching, setIsSearching] = useState(false);
  const [filing_ids, setFilingIds] = useState<string[]>([]);
  const [filings, setFilings] = useState<Filing[]>([]);
  const inital_page_load = 2;
  const [page, setPage] = useState(inital_page_load);

  const getUpdates = async () => {
    if (filings.length > 0) {
      console.log("already have filings");
      return
    } 
    setIsSearching(true);
    console.log("getting recent updates");
    const new_filings = await getRecentFilings(0, 40 * inital_page_load);

    setFilings(new_filings);
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
      // console.log(new_filings);
      if (filings.length > 0) {
        setFilings((old_filings: Filing[]) => [...old_filings, ...new_filings]);
      }
    } catch (error) {
      console.log("error getting more filings");
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
            timeoutSeconds={10}
          />
        </div>
      }
    >
      <FilingTable filings={filings} DocketColumn />
    </InfiniteScroll>
  );
}
