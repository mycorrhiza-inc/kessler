"use client";
import AllInOneServerSearch from "@/components/NewSearch/AllInOneServerSearch";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import { useState } from "react";

const OrgLookupPage = () => {
  const [queryString, setQueryString] = useState("");
  return (
    <>
      <h1 className="text-3xl font-bold">Organizations</h1>
      <div className="pr-4 w-full">
        <AllInOneServerSearch searchType={GenericSearchType.Organization} initialQuery="" initialFilters={[]} />
      </div>
    </>
  );
};

export default OrgLookupPage;
