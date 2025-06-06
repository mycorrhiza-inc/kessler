"use client";
import SearchBox from "@/components/Search/SearchBox";
import { PageContextMode } from "@/lib/types/SearchTypes";
import { useState } from "react";

const OrgLookupPage = () => {
  const [queryString, setQueryString] = useState("");
  return (
    <>
      <h1 className="text-3xl font-bold">Organizations</h1>
      <div className="pr-4 w-full">
        {/* <OrganizationTableInfiniteScroll lookup_data={{ query: queryString }} /> */}
      </div>
    </>
  );
};

export default OrgLookupPage;
