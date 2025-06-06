import SearchResultsServerStandalone from "@/components/Search/SearchResultsServer";
import { BreadcrumbValues } from "@/components/SitemapUtils";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import { stateFromHeaders } from "@/lib/nextjs_misc";
import { headers } from "next/headers";

export const metadata = {
  title: "Files - Kessler",
  description: "Search all availible files.",
};
export default function Page() {
  const headersList = headers();
  const state = stateFromHeaders(headersList);
  const breadcrumbs: BreadcrumbValues = {
    state: "ny",
    breadcrumbs: [{ title: "Files", value: "files" }],
  };
  return (
    <>
      <div className="flex">
        <h1 className="text-3xl font-bold">Files Search</h1>
      </div>
      <SearchResultsServerStandalone
        searchInfo={{ query: "", search_type: GenericSearchType.Filling }}
      />
    </>
  );
}
