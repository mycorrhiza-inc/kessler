"use client"
import ConversationTableInfiniteScroll from "@/components/Organizations/ConversationTable";
import SearchBox, { PageContextMode } from "@/components/Search/SearchBox";
import { BreadcrumbValues } from "@/components/SitemapUtils";

const ConvoLookupPage = () => {
  return (
    <>
      <h1 className="text-3xl font-bold">Dockets</h1>
      <div className="pr-4 w-full">
        <SearchBox input={{ page_context: PageContextMode.Conversations }} />
        <ConversationTableInfiniteScroll />
      </div>
    </>
  );
};

export default ConvoLookupPage;
