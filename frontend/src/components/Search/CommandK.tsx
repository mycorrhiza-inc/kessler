import { useSearchState } from "@/lib/hooks/useSearchState";
import { SearchResultsComponent } from "./SearchResults";

function CommandKSearch() {}

function SearchCommand() {
  const {
    searchQuery,
    isSearching,
    filters,
    setSearchQuery,
    triggerSearch,
    getResultsCallback,
    searchTriggerIndicator,
    ...searchState
  } = useSearchState();

  const handleSearch = (query: string) => {
    setSearchQuery(query);
    triggerSearch();
  };

  return (
    <div className="p-4">
      <input
        value={searchQuery}
        onChange={(e) => handleSearch(e.target.value)}
        placeholder="Search..."
        className="input input-bordered w-full mb-4"
      />

      <SearchResultsComponent
        isSearching={isSearching}
        reloadOnChange={searchTriggerIndicator}
        searchGetter={getResultsCallback}
      ></SearchResultsComponent>
    </div>
  );
}
