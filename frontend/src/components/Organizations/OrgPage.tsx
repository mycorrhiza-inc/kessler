"use client";
import axios from "axios";
import { useState, useEffect } from "react";
import { Filing } from "../../lib/types/FilingTypes";
import { FilingTable } from "../Conversations/FilingTable";
import Navbar from "../Navbar";
import { getFilingMetadata, getRecentFilings } from "@/lib/requests/search";

import InfiniteScroll from "react-infinite-scroll-component";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import { getOrganizationInfo } from "@/lib/requests/organizations";

function ConvertToFiling(data: any): Filing {
  const newFiling: Filing = {
    id: data.sourceID,
  };

  return newFiling;
}

export default function OrgPage({ orgId }: { orgId: string }) {
  const [searchResults, setSearchResults] = useState([]);
  const [isSearching, setIsSearching] = useState(false);
  // FIXME: this is horrible, please fix this right after the mid nov jvp meeting
  const [authorInfo, setAuthorInfo] = useState<any>({});
  const [filing_ids, setFilingIds] = useState<string[]>([]);
  const [filings, setFilings] = useState<Filing[]>([]);
  const [page, setPage] = useState(0);

  const getUpdates = async () => {
    setIsSearching(true);
    console.log("getting recent updates");
    const data = await getOrganizationInfo(orgId);
    console.log(data);
    setAuthorInfo(data);

    const ids = data.files_authored_ids;
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

      setFilings((previous) => {
        const existingIds = new Set(previous.map((f) => f.id));
        const uniqueNewFilings = newFilings.filter(
          (f) => !existingIds.has(f.id),
        );
        console.log(" uniques: ", uniqueNewFilings);
        console.log("all data: ", [...previous, ...uniqueNewFilings]);
        return [...previous, ...uniqueNewFilings];
      });
    };

    fetchFilings();
  }, [filing_ids]);

  // useEffect(() => {
  //   getUpdates();
  // }, []);

  return (
    <>
      <Navbar user={null} />

      <div className="w-full h-full p-20">
        <h1 className=" text-2xl font-bold">Recent Updates</h1>
        {isSearching ? (
          <LoadingSpinner />
        ) : (
          <FilingTable filings={filings} scroll={false} />
        )}
      </div>
    </>
  );
}
