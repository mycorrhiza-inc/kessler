import { Suspense } from "react";
import { HomeSearchBarClientBaseUrl } from "./HomeSearch";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import SearchResultsServerStandalone from "../Search/SearchResultsServerStandalone";


interface AIOServerProps {
  initialQuery: string;
}
const AllInOneServerSearch = ({ initialQuery }: AIOServerProps) => {

  return (
    <>
      <div className="flex flex-col items-center justify-center bg-base-100 p-4">
        <HomeSearchBarClientBaseUrl
          baseUrl="/search"
          initialState={initialQuery}
        />
      </div>
      <Suspense
        fallback={
          <LoadingSpinner loadingText="Fetching results from server." />
        }
      >
        <SearchResultsServerStandalone searchInfo={searchInfo} />
      </Suspense>
    </>
  )

}
