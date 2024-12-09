"use client";
import axios from "axios";
import { useState, useEffect } from "react";
import { Filing } from "../../lib/types/FilingTypes";
import { FilingTable } from "@/components/Tables/FilingTable";
import { getFilingMetadata, getRecentFilings } from "@/lib/requests/search";

import LoadingSpinner from "../styled-components/LoadingSpinner";

export default function OrganizationFileTable({
  orgId,
  filing_ids,
}: {
  orgId: string;
  filing_ids: string[];
}) {
  const [isSearching, setIsSearching] = useState(false);
  // FIXME: this is horrible, please fix this right after the mid nov jvp meeting
  const [filings, setFilings] = useState<Filing[]>([]);
  // const [page, setPage] = useState(0);

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

  // useEffect(() => {
  //   getUpdates();
  // }, []);

  return isSearching ? (
    <LoadingSpinner />
  ) : (
    <FilingTable filings={filings} DocketColumn />
  );
}
