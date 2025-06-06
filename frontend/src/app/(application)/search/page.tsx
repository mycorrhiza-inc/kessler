import React from "react";
import {
  GenericSearchType,
} from "@/lib/adapters/genericSearchCallback";
import AIOServerSearch from "@/components/NewSearch/AIOServerSearch";

interface SearchPageProps {
  searchParams: { q?: string };
}

export default function Page({ searchParams }: SearchPageProps) {
  const initialQuery = (searchParams.q || "").trim();

  return (
    <AIOServerSearch searchType={GenericSearchType.Filling} initialQuery={initialQuery} initialFilters={[]} />
  );
}
