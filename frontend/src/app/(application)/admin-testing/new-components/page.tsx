"use client";
import Card from "@/components/NewSearch/GenericResultCard";
import SideFilters from "@/components/NewSearch/SideFilters";
import { SearchResultsHomepageComponent } from "@/components/Search/SearchResults";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import { generateFakeResults } from "@/lib/search/search_utils";

const exampleFilters = {
  states: [
    { label: "New York", value: "NY" },
    { label: "California", value: "CA" },
    // ... more states
  ],
  authors: [
    { label: "New York State Department of Public Service", value: "NY_DPS" },
    { label: "Edison Water Co", value: "EDISON_WATER" },
    // ... more authors
  ],
  dockets: [
    { label: "18-M-0084", value: "18-M-0084" },
    // ... more dockets
  ],
};

export default function Page() {
  return (
    <div className="flex w-full">
      {/* Main search results content */}
      <div className="grid grid-cols-1 gap-4 p-8 w-full">
        <SearchResultsHomepageComponent
          isSearching={true}
          reloadOnChange={0}
          searchInfo={{ query: "", search_type: GenericSearchType.Filling }}
        />
      </div>
      <SideFilters
        states={exampleFilters.states}
        authors={exampleFilters.authors}
        dockets={exampleFilters.dockets}
      />
    </div>
  );
}
