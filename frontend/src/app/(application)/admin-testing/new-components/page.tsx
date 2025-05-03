"use client";
import Card from "@/components/NewSearch/GenericResultCard";
import SideFilters from "@/components/NewSearch/SideFilters";
import { SearchResultsComponent } from "@/components/Search/SearchResults";
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
        <SearchResultsComponent
          isSearching={true}
          reloadOnChange={0}
          searchGetter={generateFakeResults}
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
