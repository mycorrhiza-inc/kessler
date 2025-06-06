"use client";
import React, { useState, useEffect } from "react";
import { GiMushroomsCluster } from "react-icons/gi";
import StateSelector from "./StateSelector";
import SearchBox from "./SearchBox";
import { useKesslerStore } from "@/lib/store";
import { DynamicSingleSelect } from "@/components/Filters/FilterSingleSelect";

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

const testFieldDefinition = {
  id: "state-selector",
  displayName: "State",
  description: "Select a US state",
  placeholder: "Choose a state...",
  options: states.map(state => ({
    value: state.toLowerCase().replace(/\s+/g, '-'),
    label: state,
    disabled: false
  }))
};

// Helper function to transform state name to expected value format
const transformStateToValue = (stateName: string) => {
  return stateName.toLowerCase().replace(/\s+/g, '-');
};

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
  const [showDropdown, setShowDropdown] = useState(false);
  const [dropdownVisible, setDropdownVisible] = useState(false);
  const globalStore = useKesslerStore();

  // Default state value for consistent server/client rendering
  const defaultStateValue = transformStateToValue("New York");

  // Load dropdown only after initial render is complete
  useEffect(() => {
    // Use requestAnimationFrame to ensure DOM is fully rendered
    const timer = requestAnimationFrame(() => {
      setShowDropdown(true);
      // Add a small delay for the fade-in effect
      setTimeout(() => {
        setDropdownVisible(true);
      }, 50);
    });

    return () => cancelAnimationFrame(timer);
  }, []);

  // Alternative: Use setTimeout for delayed loading
  // useEffect(() => {
  //   const timer = setTimeout(() => {
  //     setShowDropdown(true);
  //   }, 100); // Small delay to ensure page is rendered
  //   
  //   return () => clearTimeout(timer);
  // }, []);

  const selectedState = globalStore.defaultState ||
    (initialState ? transformStateToValue(initialState) : defaultStateValue);

  const setSelectedState = globalStore.setDefaultState;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      setTriggeredQuery(searchQuery.trim());
    }
  };

  const DropDownProps = {
    fieldDefinition: testFieldDefinition,
    value: selectedState,
    onChange: (value: string) => setSelectedState(value),
    onFocus: () => console.log("Component focused"),
    onBlur: () => console.log("Component blurred"),
    dynamicWidth: true,
    minWidth: "200px",
    defaultValue: defaultStateValue,
    disabled: false,
    className: "",
    allowClear: true
  }

  return (
    <div className="flex flex-col items-center space-y-6 w-full max-w-md overflow-visible relative z-10" >
      {/* Logo and Title */}
      <div className="flex flex-col items-center space-y-2" >
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

      {/* Conditional rendering with fade-in animation */}
      {showDropdown ? (
        <div
          className={`transition-opacity duration-500 ease-in-out ${dropdownVisible ? 'opacity-100' : 'opacity-0'
            }`}
        >
          <DynamicSingleSelect {...DropDownProps} />
        </div>
      ) : (
        /* Placeholder with same dimensions to prevent layout shift */
        <div className="w-full min-h-[3rem] p-3 border-2 border-gray-300 rounded-lg bg-gray-50 animate-pulse flex items-center justify-between opacity-100">
          <span className="text-gray-400">Loading state selector...</span>
          <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      )}
    </div>
  );
}
