"use client";
import axios from "axios";
import { useState, useEffect } from "react";
import { Filing } from "../../lib/types/FilingTypes";
import { FilingTable } from "../Conversations/FilingTable";
import Navbar from "../Navbar";
import { getFilingMetadata, getRecentFilings } from "@/lib/requests/search";

function ConvertToFiling(data: any): Filing {
  const newFiling: Filing = {
    id: data.sourceID,
  };

  return newFiling;
}

export default function RecentUpdatesView() {
  const [searchResults, setSearchResults] = useState([]);
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
    <>
      <Navbar user={null} />

      <div className="w-full h-full p-20">
        <div className="w-full h-full p-10 card grid grid-flow-rows box-border border-2 border-black ">
          <h1 className=" text-2xl font-bold">Recent Updates</h1>
          <FilingTable filings={filings} scroll={false} />
          <button onClick={() => getMore()}>Get More</button>
        </div>
      </div>
    </>
  );
}
