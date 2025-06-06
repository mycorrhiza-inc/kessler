import Link from "next/link";
import { useState, useEffect, Suspense } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { ExperimentalChatModalClickDiv } from "../Chat/ChatModal";
import HomeSearchBar from "../NewSearch/HomeSearch";
import SearchResultsClient from "../Search/SearchResultsClient";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";
import HomePageClient from "./HomePageClient";
import SearchResultsServerStandalone from "../Search/SearchResultsServerStandalone";
import LoadingSpinner from "../styled-components/LoadingSpinner";

export const HomePageServer = () => {
  return (
    <Suspense fallback={<LoadingSpinner loadingText="Getting server data" />}>
      <HomePageClient
        serverRecentTables={
          <Suspense
            fallback={<LoadingSpinner loadingText="Getting server data" />}
          >
            <HomeRecentTables />
          </Suspense>
        }
      />
    </Suspense>
  );
};

const HomeRecentTables = () => {
  return (
    <>
      <div className="grid grid-cols-2 w-full z-1">
        <div>
          <Link className="text-3xl font-bold hover:underline" href="/dockets">
            Dockets
          </Link>
          <SearchResultsServerStandalone
            searchInfo={{
              query: "",
              search_type: GenericSearchType.Docket,
            }}
          />
          <div className="max-h-[600px] overflow-x-hidden border-r pr-4"></div>
        </div>
        <div className="z-1">
          <Link
            className="text-3xl font-bold hover:underline mb-5 p-10"
            href="/orgs"
          >
            Organizations
          </Link>

          <SearchResultsServerStandalone
            searchInfo={{
              query: "",
              search_type: GenericSearchType.Organization,
            }}
          />
          <div className="max-h-[600px] overflow-x-hidden pl-4"></div>
        </div>
      </div>
      <ExperimentalChatModalClickDiv
        className="btn btn-accent w-full"
        inheritedFilters={[]}
      >
        Unsure of what to do? Try chatting with the entire New York PUC
      </ExperimentalChatModalClickDiv>

      <h1 className=" text-2xl font-bold">Newest Docs</h1>
      <SearchResultsServerStandalone
        searchInfo={{
          query: "",
          search_type: GenericSearchType.Filling,
        }}
      />
    </>
  );
};
