"use client";
import ConversationTableInfiniteScroll from "@/components/LookupPages/ConvoLookup/ConversationTable";
import SearchBox from "@/components/Search/SearchBox";
import { BreadcrumbValues } from "@/components/SitemapUtils";
import { PageContextMode } from "@/lib/types/SearchTypes";
import { useState } from "react";

const ConvoLookupPage = () => {
  const [queryString, setQueryString] = useState("");

  return (
    <>
      <h1 className="text-3xl font-bold">Dockets</h1>
      <div className="pr-4 w-full">
        <SearchBox input={{ pageContext: PageContextMode.Conversations }} />
        <ConversationTableInfiniteScroll lookup_data={{ query: queryString }} />
      </div>
    </>
  );
};

export default ConvoLookupPage;
