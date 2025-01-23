"use client";
import React, { useEffect, useRef, useState } from "react";
import { AngleDownIcon, AngleUpIcon } from "../Icons";
import { AuthorInfoPill, subdividedHueFromSeed } from "../Tables/TextPills";
import { QueryFileFilterFields, QueryDataFile } from "@/lib/filters";
import { Query } from "pg";

// Mock API call
type Suggestion = {
  id: string;
  type: string;
  label: string;
  value: string;
};

type Filter = {
  id: string;
  type: string;
  label: string;
  exclude?: boolean;
  excludable: boolean;
};

const mockFetchSuggestions = async (query: string): Promise<Suggestion[]> => {
  // Simulate API delay
  await new Promise((resolve) => setTimeout(resolve, 300));

  const suggestions: Suggestion[] = [
    {
      id: "00000000-0000-0000-0000-000000000000",
      type: "file_class",
      label: "Plans and Proposals",
      value: "Plans and Proposals",
    },
    {
      id: "00000000-0000-0000-0000-000000000000",
      type: "file_class",
      label: "Corrospondence",
      value: "Corrospondence",
    },
    {
      id: "00000000-0000-0000-0000-000000000000",
      type: "file_class",
      label: "Exhibits",
      value: "Exhibits",
    },
    {
      id: "00000000-0000-0000-0000-000000000000",
      type: "file_class",
      label: "Testimony",
      value: "Testimony",
    },
    {
      id: "00000000-0000-0000-0000-000000000000",
      type: "file_class",
      label: "Reports",
      value: "Reports",
    },
    {
      id: "00000000-0000-0000-0000-000000000000",
      type: "file_class",
      label: "Comments",
      value: "Comments",
    },
    {
      id: "00000000-0000-0000-0000-000000000000",
      type: "file_class",
      label: "Attachment",
      value: "Attachment",
    },

    {
      id: "0b544651-0226-4e0d-83af-184ef5aad4e5",
      type: "organization",
      label: "New York State Department of Public Service",
      value: "acme",
    },
    {
      id: "be6aa9d6-e03f-4f85-a2f4-ae7e14199ec4",
      type: "organization",
      label: "Protect Our Coast - LINY",
      value: "apple",
    },
    {
      id: "24-E-0165",
      type: "case",
      label: "24-E-0165: Commission Regarding the Grid of the Future",
      value: "bug-123",
    },
  ].filter(
    (s) =>
      s.label.toLowerCase().includes(query.toLowerCase()) ||
      s.type.toLowerCase().includes(query.toLowerCase()),
  );

  return suggestions;
};

const AdvancedSearch = () => {
  const [open, setOpen] = useState(false);

  const flip = () => {
    setOpen(!open);
  };

  return (
    <div className="p-4 text-base-content" onClick={flip}>
      <div className="tooltip" data-tip="Advanced Search">
        {open ? <AngleUpIcon /> : <AngleDownIcon />}
      </div>
    </div>
  );
};
type FiltersPoolProps = {
  selected: Filter[];
  handleFilterRemove: (filterId: string) => void;
  flipExclude: (filterId: string) => void;
};

const FiltersPool: React.FC<FiltersPoolProps> = ({
  selected,
  handleFilterRemove,
  flipExclude,
}) => {
  const [hoveredId, setHoveredId] = useState<string | null>(null);

  if (!selected || selected.length === 0) {
    return null;
  }
  const bgcolor = (filter: Filter) => {
    if (filter.exclude) {
      return "#fde8e8";
    }
    if (filter.type === "organization") {
      return subdividedHueFromSeed(filter.id);
    }
    if (filter.type === "case") {
      return subdividedHueFromSeed(filter.label);
    }
    if (filter.type === "text") {
      return "oklch(90% 0.01 30)";
    }
  };
  return (
    selected.length > 0 && (
      <div>
        <div className="divider pl-5 pr-5"></div>
        <div className="flex flex-wrap gap-2 p-2">
          {selected.map((filter) => (
            <div
              key={filter.id}
              className="flex items-center gap-1 px-3 py-1 rounded-full text-sm group"
              style={{
                backgroundColor: bgcolor(filter),
                color: filter.exclude ? "#f56565" : "black",
              }}
              onMouseEnter={() => setHoveredId(filter.id)}
              onMouseLeave={() => setHoveredId(null)}
            >
              <span className="font-medium flex items-center gap-1">
                {filter.excludable &&
                  !filter.exclude &&
                  hoveredId === filter.id && (
                    <button
                      onClick={() => flipExclude(filter.id)}
                      className="text-gray-500 hover:underline"
                    >
                      exclude
                    </button>
                  )}
                {filter.excludable && filter.exclude && (
                  <button
                    onClick={() => flipExclude(filter.id)}
                    className={`${hoveredId === filter.id ? "line-through" : ""}`}
                  >
                    exclude
                  </button>
                )}
                {filter.type}: {filter.label}
              </span>
              <button
                onClick={() => handleFilterRemove(filter.id)}
                className="ml-1 text-gray-500 hover:text-black font-bold"
              >
                ×
              </button>
            </div>
          ))}
        </div>
      </div>
    )
  );
};

const suggestionToFilter = (suggestion: Suggestion): Filter => {
  if (suggestion.type === "text") {
    return { ...suggestion, exclude: false, excludable: false };
  }
  return { ...suggestion, exclude: false, excludable: true };
};

export enum PageContextMode {
  Files,
  Organizations,
  Dockets,
}
export interface FileSearchBoxProps {
  pageContext: PageContextMode.Files;
  setSearchData: React.Dispatch<React.SetStateAction<QueryDataFile>>;
}
export interface OrgSearchBoxProps {
  page_context: PageContextMode.Organizations;
}
export interface DocketSearchBoxProps {
  page_context: PageContextMode.Dockets;
}

export type SearchBoxInputProps =
  | FileSearchBoxProps
  | OrgSearchBoxProps
  | DocketSearchBoxProps;

export type FilterTypeDict = { [key: string]: Filter[] };
const setSearchFilters = (props: SearchBoxInputProps, filters: Filter[]) => {
  const filterTypeDict = filters.reduce(
    (acc: { [key: string]: Filter[] }, filter: Filter) => {
      if (!acc[filter.type]) {
        acc[filter.type] = [];
      }
      acc[filter.type].push(filter);
      return acc;
    },
    {},
  );

  if (props && "pageContext" in props) {
    if (props.pageContext === PageContextMode.Files) {
      const fileProps = props as FileSearchBoxProps;
      fileProps.setSearchData((previous_file_filters) => {
        const new_filters = generateFileFiltersFromFilterList(
          previous_file_filters,
          filterTypeDict,
        );
        return new_filters;
      });
      return;
    }
    if (props.pageContext === PageContextMode.Organizations) {
      return;
    }
    if (props.pageContext === PageContextMode.Dockets) {
      return;
    }
  }
};

const generateFileFiltersFromFilterList = (
  previous_file_filters: QueryDataFile,
  filterTypeDict: { [key: string]: Filter[] },
) => {
  const new_file_filters = { ...previous_file_filters };

  if (filterTypeDict.text) {
    if (filterTypeDict.text.length > 1) {
      console.log("This paramater shouldnt be more then length 1, ignoring ");
    }
    const first_filter_text = filterTypeDict.text[0].label;
    new_file_filters.query = first_filter_text;
    console.log("Filters are being updated with text");
  } else {
    new_file_filters.query = "";
  }
  return new_file_filters;
};
const SearchBox = ({ input }: { input: SearchBoxInputProps }) => {
  const [query, setQuery] = useState("");
  const [suggestions, setSuggestions] = useState<Suggestion[]>([]);
  const [selectedFilters, setSelectedFilters] = useState<Filter[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [highlightedIndex, setHighlightedIndex] = useState(0);
  const searchContainerRef = useRef<HTMLDivElement>(null);

  // Handle clicks outside of the search container
  useEffect(() => {
    const handleClickOutside = (event: any) => {
      // Check if the click was outside and we have suggestions open
      if (
        searchContainerRef.current &&
        !searchContainerRef.current.contains(event.target) &&
        suggestions.length > 0
      ) {
        console.log("Click outside detected, closing suggestions");
        setSuggestions([]);
        setQuery("");
      }
    };

    // Use mousedown and touchstart for better mobile support
    document.addEventListener("click", handleClickOutside, true);
    document.addEventListener("touchstart", handleClickOutside, true);

    return () => {
      document.removeEventListener("click", handleClickOutside, true);
      document.removeEventListener("touchstart", handleClickOutside, true);
    };
  }, [suggestions.length]); // Add suggestions.length as dependency
  useEffect(() => {
    setSearchFilters(input, selectedFilters);
  }, [selectedFilters]);

  const wrapReturnedSuggestions = (
    suggestions: Suggestion[],
    new_query: string,
  ) => {
    const text_query_suggestion = {
      id: "00000000-0000-0000-0000-000000000000",
      type: "text",
      label: new_query,
      value: new_query,
    };
    return [text_query_suggestion, ...suggestions];
  };

  const handleInputChange = async (e: any) => {
    const newQuery = e.target.value;
    setQuery(newQuery);

    if (newQuery.trim()) {
      setIsLoading(true);
      const results = wrapReturnedSuggestions(
        await mockFetchSuggestions(newQuery),
        newQuery,
      );
      setSuggestions(results);
      setHighlightedIndex(0); // Reset highlight to first option
      setIsLoading(false);
    } else {
      setSuggestions([]);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (!suggestions.length) return;

    switch (e.key) {
      case "ArrowDown":
        e.preventDefault();
        setHighlightedIndex((prev) =>
          prev < suggestions.length - 1 ? prev + 1 : prev,
        );
        break;
      case "ArrowUp":
        e.preventDefault();
        setHighlightedIndex((prev) => (prev > 0 ? prev - 1 : prev));
        break;
      case "Enter":
        e.preventDefault();
        if (suggestions[highlightedIndex]) {
          handleSuggestionClick(suggestions[highlightedIndex]);
        }
        break;
    }
  };

  const handleSuggestionClick = (suggestion: Suggestion) => {
    if (!selectedFilters.some((f) => f.id === suggestion.id)) {
      setSelectedFilters([...selectedFilters, suggestionToFilter(suggestion)]);
    }
    setQuery("");
    setSuggestions([]);
  };

  const handleFilterRemove = (filterId: string) => {
    setSelectedFilters(selectedFilters.filter((f) => f.id !== filterId));
  };

  const flipExclude = (filterId: string) => {
    setSelectedFilters(
      selectedFilters.map((f) => {
        if (f.id === filterId) {
          return { ...f, exclude: !f.exclude };
        }
        return f;
      }),
    );
  };

  return (
    <div className="p-4 max-w-xl mx-auto">
      <div className="flex flex-col gap-2">
        {/* Search container */}
        <div className="relative">
          {/* Search input */}
          <div className="relative flex flex-row">
            <input
              type="text"
              value={query}
              onChange={handleInputChange}
              onKeyDown={handleKeyDown}
              placeholder="Search anything..."
              className="w-full p-3 border rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none bg-base-100"
            />

            {isLoading && (
              <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                <div className="animate-spin h-5 w-5 border-2 border-blue-500 border-t-transparent rounded-full" />
              </div>
            )}
            <AdvancedSearch />
          </div>

          {/* Suggestions dropdown - Now positioned relative to search container */}
          {suggestions.length > 0 && (
            <div className="absolute left-0 right-0 top-full mt-1 z-50 h-auto bg-base-100 border rounded-lg shadow-lg">
              <ul className=" max-h-60 overflow-auto">
                {suggestions.map((suggestion, index) => (
                  <li key={suggestion.id}>
                    <button
                      onClick={() => handleSuggestionClick(suggestion)}
                      onMouseEnter={() => setHighlightedIndex(index)}
                      className={`w-full px-4 py-3 text-left transition-colors ${
                        index === highlightedIndex
                          ? "bg-primary/10 text-primary"
                          : "hover:secondary-content"
                      }`}
                    >
                      <span className={`text-sm font-medium text-primary`}>
                        {suggestion.type}:
                      </span>{" "}
                      <span className="text-base-content">
                        {suggestion.label}
                      </span>
                    </button>
                  </li>
                ))}
              </ul>
              <FiltersPool
                selected={selectedFilters}
                handleFilterRemove={handleFilterRemove}
                flipExclude={flipExclude}
              />
            </div>
          )}
        </div>

        <FiltersPool
          selected={selectedFilters}
          handleFilterRemove={handleFilterRemove}
          flipExclude={flipExclude}
        />
      </div>
    </div>
  );
};

export default SearchBox;
