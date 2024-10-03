import SearchResult from "@/components/SearchResult";
import DisplayCard from "@/components/DisplayCard";

interface SearchResultBoxProps {
  showCard?: string;
  searchResults?: any[];
  isSearching?: boolean;
}

const SearchResultBox = ({
  searchResults,
  isSearching,
  showCard,
}: SearchResultBoxProps) => {
  const noResultsLoaded =
    !isSearching && Array.isArray(searchResults) && searchResults.length === 0;
  return (
    <div className="searchResultsContent flex flex-col justify-center items-center pb-52 pt-24 space-y-2">
      <div
        className="searchResultsContent"
        style={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          gap: "2rem",
        }}
      >
        {!isSearching &&
          searchResults &&
          searchResults.map((result, index) => (
            <SearchResult key={index} data={result} />
          ))}

        {noResultsLoaded && showCard && <DisplayCard cardType={showCard} />}
        {noResultsLoaded && showCard == "" && <>No results found :(</>}
        {isSearching && (
          <>
            loading<span className="loading loading-dots loading-lg"></span>
          </>
        )}
      </div>
    </div>
  );
};

export default SearchResultBox;
