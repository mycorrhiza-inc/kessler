import React from "react";
import {
  GenericSearchType,
} from "@/lib/adapters/genericSearchCallback";
import AllInOneServerSearch from "@/components/NewSearch/AllInOneServerSearch";

interface SearchPageProps {
  searchParams: { q?: string };
}

export default function Page({ searchParams }: SearchPageProps) {
  const initialQuery = (searchParams.q || "").trim();

  return (
    <AllInOneServerSearch searchType={GenericSearchType.Filling} initialQuery={initialQuery} initialFilters={[]} />
  );
}
