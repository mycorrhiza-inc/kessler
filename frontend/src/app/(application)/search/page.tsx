"use client";
import React from "react";
import { useUrlParams } from "@/lib/hooks/useUrlParams";
import AllInOneClientSearch from "@/stateful_components/SearchBar/AllInOneClientSearch";
import DynamicFilters from "@/stateful_components/Filters/DynamicFilters";

export default function SearchPage() {
  const urlParams = useUrlParams();

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Search</h1>
      <AllInOneClientSearch urlParams={urlParams}
      />
      {/* <DynamicFilters filters={filters} dataset={dataset} /> */}
    </div>
  );
}
