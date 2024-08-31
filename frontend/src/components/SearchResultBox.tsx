import SearchResult from "@/components/SearchResult";
import { Stack, Box } from "@mui/joy";
const SearchResultBox = ({
  searchResults,
  isSearching,
}: {
  searchResults: any[];
  isSearching: boolean;
}) => {
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
        {searchResults.map((result, index) => (
          <SearchResult key={index} data={result} />
        ))}
      </Stack>
    </div>
  );
};



export default SearchResultBox;
