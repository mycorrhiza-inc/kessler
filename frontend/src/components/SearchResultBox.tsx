import SearchResult from "@/components/SearchResult";
import { Stack, Box } from "@mui/joy";

interface SearchResultBoxProps {
  searchResults?: any[];
  isSearching?: boolean;
}

const SearchResultBox = ({
  searchResults,
  isSearching,
}: SearchResultBoxProps) => {
  return (
    <div className="searchResults justify-center items-center ">
      <div className="searchResultsContent flex flex-col justify-center items-center pb-52 pt-24 space-y-2">
        {!isSearching &&
          searchResults &&
          searchResults.map((result, index) => (
            <SearchResult key={index} data={result} />
          ))}
        {isSearching && <span>loading...</span>}
      </div>
    </div>
  );
};

export default SearchResultBox;
