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
    <div className="searchResultsContent flex flex-col justify-center items-center pb-52 pt-24 space-y-2">
      <Stack
        className="searchResultsContent"
        style={{ justifyContent: "center" }}
        direction="column"
        justifyContent="center"
        alignItems="center"
        spacing={2}
      >
        {!isSearching &&
          searchResults &&
          searchResults.map((result, index) => (
            <SearchResult key={index} data={result} />
          ))}
        {!isSearching &&
          Array.isArray(searchResults) &&
          searchResults.length === 0 && <>No results found</>}
        {isSearching && <>loading...</>}
      </Stack>
    </div>
  );
};

export default SearchResultBox;
