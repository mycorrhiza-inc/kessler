"use client";
import { useState } from "react";
import { ConvoSearchRequestData } from "../SearchRequestData";
import AIOServerSearch from "@/components/NewSearch/AIOServerSearch";
import { GenericSearchType } from "@/lib/adapters/genericSearchCallback";

const ConvoLookupPage = () => {
  return (
    <>
      <h1 className="text-3xl font-bold">Dockets</h1>
      <div className="pr-4 w-full">

        <AIOServerSearch searchType={GenericSearchType.Docket} initialQuery="" initialFilters={[]} />
      </div>
    </>
  );
};

export default ConvoLookupPage;
