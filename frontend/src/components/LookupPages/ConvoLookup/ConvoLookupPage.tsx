import ConversationTableInfiniteScroll from "@/components/Organizations/ConversationTable";
import OrganizationTableInfiniteScroll from "@/components/Organizations/OrganizationTable";
import PageContainer from "@/components/Page/PageContainer";
import SearchBox from "@/components/Search/SearchBox";
import { BreadcrumbValues } from "@/components/SitemapUtils";

const ConvoLookupPage = ({
  breadcrumbs,
}: {
  breadcrumbs: BreadcrumbValues;
}) => {
  return (
    <PageContainer breadcrumbs={breadcrumbs}>
      <h1 className="text-3xl font-bold">Dockets</h1>
      <div className="overflow-x-hidden border-r pr-4">
        <SearchBox />
        <ConversationTableInfiniteScroll />
      </div>
    </PageContainer>
  );
};

export default ConvoLookupPage;
