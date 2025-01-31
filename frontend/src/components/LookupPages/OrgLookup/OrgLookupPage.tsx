"use client";
import SearchBox from "@/components/Search/SearchBox";
import { PageContextMode } from "@/lib/types/SearchTypes";
import OrganizationTableInfiniteScroll from "./OrganizationTable";

const OrgLookupPage = () => {
  return (
    <>
      <h1 className="text-3xl font-bold">Organizations</h1>
      <div className="pr-4 w-full">
        <SearchBox input={{ pageContext: PageContextMode.Organizations }} />
        <OrganizationTableInfiniteScroll />
      </div>
    </>
  );
};

export default OrgLookupPage;
