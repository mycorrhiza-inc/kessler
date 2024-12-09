"use client";
import axios from "axios";
import { useState, useEffect } from "react";
import { Filing } from "@/lib/types/FilingTypes";
import { FilingTable } from "@/components/Tables/FilingTable";
import { getFilingMetadata, getRecentFilings } from "@/lib/requests/search";

import InfiniteScroll from "react-infinite-scroll-component";
import ConversationTable from "../Organizations/ConversationTable";
import OrganizationTable from "../Organizations/OrganizationTable";
import Link from "next/link";
import LoadingSpinnerTimeout from "../styled-components/LoadingSpinnerTimeout";
import PageContainer from "../Page/PageContainer";

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

  useEffect(() => {
    if (!filing_ids) {
      return;
    }

    const fetchFilings = async () => {
      const newFilings = await Promise.all(
        filing_ids.map(async (id) => {
          const filing_data = await getFilingMetadata(id);
          console.log("new filings", filing_data);
          return filing_data;
        }),
      );

      setFilings((previous: Filing[]): Filing[] => {
        const existingIds = new Set(previous.map((f: Filing) => f.id));
        const uniqueNewFilings = newFilings.filter(
          (f: Filing | null) => f !== null && !existingIds.has(f.id),
        ) as Filing[];
        console.log(" uniques: ", uniqueNewFilings);
        console.log("all data: ", [...previous, ...uniqueNewFilings]);
        return [...previous, ...uniqueNewFilings];
      });
    };

    fetchFilings();
  }, [filing_ids]);

  const getMore = async () => {
    setIsSearching(true);
    try {
      console.log("getting page ", page + 1);
      const data = await getRecentFilings(page);
      setPage(page + 1);
      console.log(data);
      if (data.length > 0) {
        setFilingIds([
          ...filing_ids,
          ...data.map((item: any) => item.sourceID),
        ]);
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
    <PageContainer breadcrumbs={{ breadcrumbs: [] }}>
      <div className="grid grid-cols-2 w-full">
        <div className="max-h-[600px] overflow-x-hidden border-r pr-4">
          <Link
            className="text-3xl font-bold hover:underline"
            href="/proceedings"
          >
            Proceedings
          </Link>
          <ConversationTable />
        </div>
        <div className="max-h-[600px] overflow-x-hidden pl-4">
          <Link className="text-3xl font-bold hover:underline" href="/orgs">
            Organizations
          </Link>
          <OrganizationTable />
        </div>
      </div>
      <div className="border-t my-8"></div>
      <h1 className=" text-2xl font-bold">Newest Docs</h1>
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
    </PageContainer>
  );
}
