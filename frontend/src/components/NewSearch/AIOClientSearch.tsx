import { ReactNode, Suspense } from "react";
import { HomeSearchBarClientBaseUrl } from "./HomeSearch";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import SearchResultsServerStandalone from "../Search/SearchResultsServerStandalone";


const AllInOneClientSearchComponent = ({ serverSideResults, initialQuery }: { serverSideResults: ReactNode, initialQuery: }) => {

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
        {serverSideResults}
      </Suspense>
    </>
  )

}
