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
      <Stack
        className="searchResultsContent"
        style={{ justifyContent: "center", paddingBottom: "200px", paddingTop: "100px" }}
        direction="column"
        justifyContent="center"
        alignItems="center"
        spacing={2}
      >
        {!isSearching && searchResults && searchResults.map((result, index) => (
          <SearchResult key={index} data={result} />
        ))}
        {isSearching && (<>loading...</>)}
      </Stack>
    </div>
  );
};



export default SearchResultBox;
