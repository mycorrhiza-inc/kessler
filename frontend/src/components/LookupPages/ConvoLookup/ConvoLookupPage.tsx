import ConversationTableInfiniteScroll from "@/components/Organizations/ConversationTable";
import SearchBox from "@/components/Search/SearchBox";
import { BreadcrumbValues } from "@/components/SitemapUtils";

const ConvoLookupPage = () => {
  return (
    <>
      <h1 className="text-3xl font-bold">Dockets</h1>
      <div className="overflow-x-hidden border-r pr-4">
        <SearchBox />
        <ConversationTableInfiniteScroll />
      </div>
    </>
  );
};

export default ConvoLookupPage;
