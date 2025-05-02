import { useSearchState } from "@/lib/hooks/useSearchState";
import { SearchResultsWrapper } from "./SearchResults";

function CommandKSearch() {}

function SearchCommand() {
  const { searchQuery, isSearching, handleSearch, ...searchState } =
    useSearchState();

  return (
    <div className="p-4">
      <input
        value={searchQuery}
        onChange={(e) => handleSearch(e.target.value)}
        placeholder="Search..."
        className="input input-bordered w-full mb-4"
      />

      <SearchResultsWrapper
        isSearching={isSearching}
        {...searchState}
      ></SearchResultsWrapper>
    </div>
  );
}
