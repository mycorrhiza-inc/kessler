import { motion, AnimatePresence } from "framer-motion";
import { useEffect, useState } from "react";
import InfiniteScrollPlus from "../InfiniteScroll/InfiniteScroll";
import DummyResults from "../NewSearch/DummySearchResults";

interface SearchResultsProps {
  isSearching: boolean;
  isLoading: boolean;
  error: string | null;
  searchResults: SearchResult[];
  children: React.ReactNode;
}

type SearchResult = any;

function RawSearchResults({
  searchResults,
}: {
  searchResults: SearchResult[];
}) {
  return <DummyResults />;
}

interface PaginationData {
  page: Number;
  limit: Number;
}

type SearchResultsGetter = (data: PaginationData) => Promise<SearchResult[]>;

interface GeneralInfiniteSearchParams {
  searchGetter: SearchResultsGetter;
  reloadOnChangeObj: any;
}

function SearchResultsInfiniteScroll({
  searchGetter,
  reloadOnChangeObj,
}: GeneralInfiniteSearchParams) {
  const [hasMore, setHasMore] = useState(true);
  const [searchData, setSearchData] = useState<SearchResult[]>([]);
  const [page, setPage] = useState(0);
  const pageSize = 40;
  const pushSearchResults = async (pagination: PaginationData) => {
    const newSearchResults = await searchGetter(pagination);
    if (newSearchResults.length != pagination.limit) {
      setHasMore(false);
    }
    setSearchData((prev) => prev.concat(newSearchResults));
  };
  const getInitialUpdates = async () => {
    setSearchData([]);
    setPage(0);
    setHasMore(true);
    const load_initial_pages = 2;
    const limit = pageSize * load_initial_pages;
    console.log("getting recent updates");
    await pushSearchResults({
      page: 0,
      limit: limit,
    });

    setPage(load_initial_pages);
  };

  const getMore = async () => {
    await pushSearchResults({
      page: page,
      limit: pageSize,
    });
    setPage((prev) => prev + 1);
  };

  // clear on reload
  // Unecessary already handled in infinite scroll
  // useEffect(() => {
  //   getInitialUpdates(); // Figure out a less confusing syntax for adding an async function to the execution queue
  // }, [reloadOnChangeObj]);

  return (
    <InfiniteScrollPlus
      loadInitial={getInitialUpdates}
      getMore={getMore}
      reloadOnChangeObj={reloadOnChangeObj}
      dataLength={searchData.length}
      hasMore={hasMore}
    >
      <RawSearchResults searchResults={searchData} />
    </InfiniteScrollPlus>
  );
}

export function SearchResultsComponent({
  isSearching,
  isLoading,
  error,
  searchResults,
}: SearchResultsProps) {
  return (
    <AnimatePresence mode="wait">
      {isSearching && (
        <motion.div
          key="search-results"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -20 }}
          transition={{ duration: 0.3 }}
          className="w-full"
        >
          {isLoading ? (
            <div className="flex justify-center p-8">
              <span className="loading loading-spinner loading-lg"></span>
            </div>
          ) : error ? (
            <div className="alert alert-error my-4">
              <span>{error}</span>
            </div>
          ) : searchResults.length > 0 ? (
            children
          ) : (
            <div className="text-center p-8">
              <h3 className="text-lg font-medium">No results found</h3>
              <p className="text-sm opacity-70 mt-1">
                Try adjusting your search terms
              </p>
            </div>
          )}
        </motion.div>
      )}
    </AnimatePresence>
  );
}
