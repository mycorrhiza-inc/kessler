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
import { PageContext } from "@/lib/page_context";
import { BreadcrumbValues } from "../SitemapUtils";
import { User } from "@supabase/supabase-js";

export default function OrganizationPage({
  user,
  breadcrumbs,
}: {
  user: User | null;
  breadcrumbs: BreadcrumbValues;
}) {
  const [isSearching, setIsSearching] = useState(false);
  // FIXME: this is horrible, please fix this right after the mid nov jvp meeting
  const [orgInfo, setOrgInfo] = useState<any>({});
  const [filing_ids, setFilingIds] = useState<string[]>([]);
  const [filings, setFilings] = useState<Filing[]>([]);
  const orgId =
    breadcrumbs.breadcrumbs[breadcrumbs.breadcrumbs.length - 1].value;
  const actual_breadcrumb_values = [
    ...breadcrumbs.breadcrumbs.slice(0, -1),
    { value: orgId, title: orgInfo.name || "Loading" },
  ];
  const actual_breadcrumbs: BreadcrumbValues = {
    breadcrumbs: actual_breadcrumb_values,
    state: breadcrumbs.state,
  };
  // const [page, setPage] = useState(0);

  const getUpdates = async () => {
    setIsSearching(true);
    console.log("getting recent updates");
    const data = await getOrganizationInfo(orgId || "");
    console.log(data);
    data.description = "Organization Descriptions Coming Soon!";
    setOrgInfo(data);

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

  useEffect(() => {
    getUpdates();
  }, []);

  return (
    <>
      <Navbar user={user} breadcrumbs={actual_breadcrumbs} />
      <div className="w-full h-full p-20">
        <h1 className=" text-2xl font-bold">Organization: {orgInfo.name}</h1>
        <p> {orgInfo.description || "Loading Organization Description"}</p>
        <h1 className=" text-2xl font-bold">Authored Documents</h1>
        {isSearching ? (
          <LoadingSpinner />
        ) : (
          <FilingTable filings={filings} DocketColumn />
        )}
      </div>
    </>
  );
}
