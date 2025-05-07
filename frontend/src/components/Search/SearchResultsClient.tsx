import { motion, AnimatePresence } from "framer-motion";
import { useState } from "react";
import InfiniteScrollPlus from "../InfiniteScroll/InfiniteScroll";
import Card, { CardSize } from "../NewSearch/GenericResultCard";
import { generateFakeResults } from "../NewSearch/DummySearchResults";
import { sleep } from "@/utils/utils";
import {
  PaginationData,
  SearchResult,
  SearchResultsGetter,
} from "@/lib/types/new_search_types";

function RawSearchResults({
  searchResults,
}: {
  searchResults: SearchResult[];
}) {
  return (
    <div className="flex w-full">
      <div className="grid grid-cols-1 gap-4 p-8 w-full">
        {searchResults.map((data, index) => (
          <Card key={index} data={data} size={CardSize.Medium} />
        ))}
      </div>
    </div>
  );
}

interface GeneralInfiniteSearchParams {
  searchGetter: SearchResultsGetter;
  reloadOnChange: number;
  offset: number;
}

function SearchResultsInfiniteScroll({
  searchGetter,
  reloadOnChange,
  offset,
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
      reloadOnChange={reloadOnChange}
      dataLength={searchData.length}
      hasMore={hasMore}
    >
      <RawSearchResults searchResults={searchData} />
    </InfiniteScrollPlus>
  );
}

export function SearchResultsClientComponent({
  searchGetter,
  reloadOnChange,
  isSearching,
}: {
  isSearching: boolean;
  searchGetter: SearchResultsGetter;
  reloadOnChange: number;
}) {
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
          <SearchResultsInfiniteScroll
            searchGetter={searchGetter}
            reloadOnChange={reloadOnChange}
          />
        </motion.div>
      )}
    </AnimatePresence>
  );
}
