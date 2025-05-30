"use client";
import React, { useState } from "react";
import { GiMushroomsCluster } from "react-icons/gi";
import StateSelector from "./StateSelector";
import SearchBox from "./SearchBox";
import { useKesslerStore } from "@/lib/store";
import { Filters } from "@/lib/types/new_filter_types";

const states = [
  "New York",
  "California",
  "Texas",
  "Florida",
  "Illinois",
  "Pennsylvania",
  "Ohio",
  "Georgia",
  "North Carolina",
  "Michigan",
];

export const HomeSearchBarClientBaseUrl = ({
  baseUrl,
  initialState,
}: {
  baseUrl: string;
  initialState?: string;
}) => {
  const handleSearch = (query: string) => {
    const q = query.trim();
    if (q) window.location.href = `${baseUrl}?q=${encodeURIComponent(q)}`;
  };
  return (
    <HomeSearchBar
      setTriggeredQuery={handleSearch}
      initialState={initialState}
    />
  );
};

export default function HomeSearchBar({
  setTriggeredQuery,
  initialState,
}: {
  setTriggeredQuery: (query: string) => void;
  initialState?: string;
}) {
  const [searchQuery, setSearchQuery] = useState("");
  const globalStore = useKesslerStore();
  const selectedState = globalStore.defaultState || initialState || "New York";
  const setSelectedState = globalStore.setDefaultState;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      setTriggeredQuery(searchQuery.trim());
    }
  };

  return (
    <div className="flex flex-col items-center space-y-6 w-full max-w-md">
      {/* Logo and Title */}
      <div className="flex flex-col items-center space-y-2">
        <div className="flex flex-row items-center space-x-9">
          <GiMushroomsCluster className="text-6xl lg:text-7xl xl:text-9xl text-base-content" />
          <h1 className="text-5xl lg:text-6xl xl:text-8xl font-bold font-serif tracking-tight">
            KESSLER
          </h1>
        </div>
        <p className="text-md xl:text-xl text-gray-600 text-center font-serif">
          Public Utility Commissions, Simplified.
        </p>
      </div>

      <form onSubmit={handleSubmit} className="w-full">
        <SearchBox
          value={searchQuery}
          onChange={setSearchQuery}
          placeholder="Grid of the Future"
        />
      </form>

      <StateSelector
        states={states}
        selectedState={selectedState}
        onSelect={setSelectedState}
      />
    </div>
  );
}
