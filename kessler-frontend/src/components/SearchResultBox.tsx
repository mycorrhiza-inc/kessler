import SearchResult from "@/components/SearchResult";

const SearchResultBox = ({
  searchResults,
  isSearching,
}: {
  searchResults: any[];
  isSearching: boolean;
}) => {
  return (
    <div className="searchResults">
      <div
        className="searchResultsContent"
        style={{ justifyContent: "center" }}
      >
        <div className="min-h-50vh">
          {isSearching ? (
            <div className="flex items-center justify-center mt-4">
              <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-white"></div>
            </div>
          ) : (
            <>
              {searchResults.map((result, index) => (
                <SearchResult key={index} data={result} />
              ))}
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default SearchResultBox;