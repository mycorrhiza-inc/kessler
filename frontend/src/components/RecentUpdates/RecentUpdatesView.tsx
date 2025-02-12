"use client";
import { useState, useEffect } from "react";
import { Filing } from "@/lib/types/FilingTypes";
import { FilingTable } from "@/components/Tables/FilingTable";
import { getRecentFilings } from "@/lib/requests/search";

import InfiniteScroll from "react-infinite-scroll-component";
import LoadingSpinnerTimeout from "../styled-components/LoadingSpinnerTimeout";
import { fi } from "date-fns/locale";
import InfiniteScrollPlus from "../InfiniteScroll/InfiniteScroll";

// TODO: Break out Recent Updates into its own component seperate from all of the homepage logic
export default function RecentUpdatesView() {
  const [isSearching, setIsSearching] = useState(false);
  const [filing_ids, setFilingIds] = useState<string[]>([]);
  const [filings, setFilings] = useState<Filing[]>([]);
  const inital_page_load = 2;
  const [page, setPage] = useState(inital_page_load);

  const getInitial = async () => {
    if (filings.length > 0) {
      console.log("already have filings");
      return;
    }
    setIsSearching(true);
    console.log("getting recent updates");
    const new_filings = await getRecentFilings(0, 40 * inital_page_load);

    setFilings(new_filings);
    setIsSearching(false);
  };

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

  return (
    <InfiniteScrollPlus
      dataLength={filings.length}
      getMore={getMore}
      hasMore={true}
      loadInitial={getInitial}
      reloadOnChangeObj={[]}
    >
      <FilingTable filings={filings} DocketColumn />
    </InfiniteScrollPlus>
  );
}
