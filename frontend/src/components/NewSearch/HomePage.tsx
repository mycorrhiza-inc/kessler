"use client";

import React, { useState } from "react";
import { GiMushroomsCluster } from "react-icons/gi";
import { MdKeyboardArrowDown } from "react-icons/md";

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

        {/* Search Container */}
        <div className="w-full space-y-4">
          {/* Search Input */}
          <div className="relative">
            <input
              type="text"
              placeholder="Grid of the Future"
              className="w-full px-4 py-3 pr-10 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6 text-gray-400 absolute right-3 top-1/2 transform -translate-y-1/2"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              />
            </svg>
          </div>

          {/* State Dropdown */}
          <div className="relative">
            <button
              onClick={() => setIsStateDropdownOpen(!isStateDropdownOpen)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg flex items-center justify-between"
            >
              <span>{selectedState}</span>
              <MdKeyboardArrowDown className="text-xl" />
            </button>

            {isStateDropdownOpen && (
              <div className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-y-auto">
                {states.map((state) => (
                  <div
                    key={state}
                    onClick={() => {
                      setSelectedState(state);
                      setIsStateDropdownOpen(false);
                    }}
                    className="px-4 py-2 hover:bg-gray-100 cursor-pointer"
                  >
                    {state}
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
