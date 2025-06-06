import { ReactNode, Suspense, useState } from "react";
import HomeSearchBar, { HomeSearchBarClientBaseUrl } from "./HomeSearch";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import SearchResultsServerStandalone from "../Search/SearchResultsServerStandalone";
import { SearchResult } from "@/lib/types/new_search_types";
import SearchResultsClient from "../Search/SearchResultsClient";
import { BackendFilterObject } from "@/lib/filters";
import { GenericSearchInfo, GenericSearchType } from "@/lib/adapters/genericSearchCallback";


interface AIOClientProps {
  searchType: GenericSearchType,
  initialQuery: string;
  initialFilters: BackendFilterObject;
  initialData?: SearchResult[];
  children?: React.ReactNode; // SSR seed
}
const AIOClientSearchComponent = ({
  searchType,
  initialQuery,
  initialFilters,
  initialData,
  children, }: AIOClientProps) => {
  const [reloadOnChange, setReloadOnChange] = useState(0);

  const initialSearchInfo: GenericSearchInfo = {
    query: initialQuery,
    search_type: searchType,
    filters: initialFilters,
  }
  const [searchInfo, setSearchInfo] = useState(initialSearchInfo)


  const handleSearch = (query: string) => {
    setSearchInfo((prev) => { return { ...prev, query: query } })
    setReloadOnChange((prev) => (prev + 1) % 1024)
  }


  return (
    <>
      <div className="flex flex-col items-center justify-center bg-base-100 p-4">

        <HomeSearchBar
          setTriggeredQuery={handleSearch}
          initialState={initialQuery}
        />
      </div>
      <SearchResultsClient
        searchInfo={searchInfo}
        reloadOnChange={reloadOnChange}
        initialData={initialData}
      >
        {children}
      </SearchResultsClient>
    </>
  )

}

export default AIOClientSearchComponent;
