import { useSearchState } from "@/lib/hooks/useSearchState";
import { SearchResults } from "./SearchResults";

export function SearchCommand() {
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

      <SearchResults isSearching={isSearching} {...searchState}></SearchResults>
    </div>
  );
}
