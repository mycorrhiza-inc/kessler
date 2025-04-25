"use client";

import React, { useState } from "react";
import { GiMushroomsCluster } from "react-icons/gi";
import { MdKeyboardArrowDown } from "react-icons/md";
import StateSelector from "./StateSelector";
import SearchBox from "./SearchBox";

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

export default function HomeSearchBar() {
  const [selectedState, setSelectedState] = useState("New York");
  const [isStateDropdownOpen, setIsStateDropdownOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-base-100 p-4">
      <div className="flex flex-col items-center space-y-6 w-full max-w-md">
        {/* Logo and Title */}
        <div className="flex flex-col items-center space-y-2">
          <div className="flex flex-row items-center space-x-9">
            <GiMushroomsCluster className="text-8xl text-base-content" />
            <h1 className="text-7xl font-bold font-serif tracking-tight">
              KESSLER
            </h1>
          </div>
          <p className="text-lg text-gray-600 text-center font-serif">
            Public Utility Commissions, Simplified.
          </p>
        </div>

        <SearchBox
          value={searchQuery}
          onChange={setSearchQuery}
          placeholder="Grid of the Future"
        />

        <StateSelector
          states={states}
          selectedState={selectedState}
          onSelect={setSelectedState}
        />
      </div>
    </div>
  );
}
