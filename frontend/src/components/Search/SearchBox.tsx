"use client";
import React, { useEffect, useRef, useState } from "react";
import { AngleDownIcon, AngleUpIcon } from "../Icons";
import { AuthorInfoPill, subdividedHueFromSeed } from "../Tables/TextPills";
import { QueryFileFilterFields, QueryDataFile } from "@/lib/filters";
import { Query } from "pg";
import {
  FileSearchBoxProps,
  Filter,
  PageContextMode,
  SearchBoxInputProps,
  Suggestion,
} from "@/lib/types/SearchTypes";
import { randomUUID } from "crypto";
import { ConvoSearchRequestData } from "../LookupPages/SearchRequestData";

const getRawSuggestions = (PageContext: PageContextMode): Suggestion[] => {
  if (PageContext === PageContextMode.Conversations) {
    return [
      {
        id: "fc001a23-4b5c-6d7e-8f9g-0h1i2j3k4l5",
        type: "nypuc_docket_industry",
        label: "Miscellaneous",
        value: "Miscellaneous",
      },
      {
        id: "fc002b34-5c6d-7e8f-9g0h-1i2j3k4l5m6",
        type: "nypuc_docket_industry",
        label: "Gas",
        value: "Gas",
      },
      {
        id: "fc003c45-6d7e-8f9g-0h1i-2j3k4l5m6n7",
        type: "nypuc_docket_industry",
        label: "Electric",
        value: "Electric",
      },
      {
        id: "fc004d56-7e8f-9g0h-1i2j-3k4l5m6n7o8",
        type: "nypuc_docket_industry",
        label: "Facility Gen.",
        value: "Facility Gen.",
      },
      {
        id: "fc005e67-8f9g-0h1i-2j3k-4l5m6n7o8p9",
        type: "nypuc_docket_industry",
        label: "Transmission",
        value: "Transmission",
      },
      {
        id: "fc006f78-9g0h-1i2j-3k4l-5m6n7o8p9q0",
        type: "nypuc_docket_industry",
        label: "Water",
        value: "Water",
      },
      {
        id: "fc007g89-0h1i-2j3k-4l5m-6n7o8p9q0r1",
        type: "nypuc_docket_industry",
        label: "Communication",
        value: "Communication",
      },
    ];
  }
  if (PageContext === PageContextMode.Organizations) {
    return [];
  }
  if (PageContext === PageContextMode.Files) {
    return [
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
        type: "docket",
        label: "24-E-0165: Commission Regarding the Grid of the Future",
        value: "bug-123",
      },

      {
        id: "fc001a23-5f7e-4b3c-9d2a-8f6e4c7d9e0b",
        type: "file_class",
        label: "Plans and Proposals",
        value: "Plans and Proposals",
      },
      {
        id: "fc002b34-6a8d-4c5f-ae3b-9g7f5d8e0f1c",
        type: "file_class",
        label: "Corrospondence",
        value: "Corrospondence",
      },
      {
        id: "fc003c45-7b9e-5d6g-bf4c-ah8g6e9f1g2d",
        type: "file_class",
        label: "Exhibits",
        value: "Exhibits",
      },
      {
        id: "fc004d56-8c0f-6e7h-cg5d-bi9h7f0g2h3e",
        type: "file_class",
        label: "Testimony",
        value: "Testimony",
      },
      {
        id: "fc005e67-9d1g-7f8i-dh6e-cj0i8g1h3i4f",
        type: "file_class",
        label: "Reports",
        value: "Reports",
      },
      {
        id: "fc006f78-0e2h-8g9j-ei7f-dk1j9h2i4j5g",
        type: "file_class",
        label: "Comments",
        value: "Comments",
      },
      {
        id: "fc007g89-1f3i-9h0k-fj8g-el2k0i3j5k6h",
        type: "file_class",
        label: "Attachment",
        value: "Attachment",
      },
    ];
  }
  console.error(
    "Unknown page context for generating raw suggestions",
    PageContext,
  );
  return [];
};
const mockFetchSuggestions = async (
  query: string,
  PageContext: PageContextMode,
): Promise<Suggestion[]> => {
  // Simulate API delay
  // await new Promise((resolve) => setTimeout(resolve, 300));

  const suggestions: Suggestion[] = getRawSuggestions(PageContext).filter(
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
    if (filter.type === "docket") {
      return subdividedHueFromSeed(filter.label);
    }
    if (filter.type === "file_class") {
      return subdividedHueFromSeed(filter.label);
    }
    if (filter.type === "text") {
      return "oklch(90% 0.01 30)";
    }
    return "oklch(90% 0.1 80)";
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
      props.setSearchQuery(getTextQueryFromFilterList(filterTypeDict));
      return;
    }
    if (props.pageContext === PageContextMode.Conversations) {
      props.setSearchData(generateConvoSearchData(filterTypeDict));
      return;
    }
  }
};

const getTextQueryFromFilterList = (filterTypeDict: {
  [key: string]: Filter[];
}) => {
  if (filterTypeDict.text) {
    if (filterTypeDict.text.length > 1) {
      console.log("This paramater shouldnt be more then length 1, ignoring ");
    }
    const first_filter_text = filterTypeDict.text[0].label;
    console.log("Filters are being updated with text");
    return first_filter_text;
  } else {
    return "";
  }
};

const generateConvoSearchData = (filterTypeDict: {
  [key: string]: Filter[];
}) => {
  const convoSearchData: ConvoSearchRequestData = {};
  const setQuery = (value: string) => {
    convoSearchData.query = value;
  };
  const setIndustry = (value: string) => {
    convoSearchData.industry_type = value;
  };
  filterExtractionHelper(filterTypeDict.text, "label", "", setQuery);
  filterExtractionHelper(
    filterTypeDict.nypuc_docket_industry,
    "label",
    "",
    setIndustry,
  );
  return convoSearchData;
};

const generateFileFiltersFromFilterList = (
  previous_file_filters: QueryDataFile,
  filterTypeDict: { [key: string]: Filter[] },
) => {
  const new_file_filters = { ...previous_file_filters };
  new_file_filters.query = getTextQueryFromFilterList(filterTypeDict);

  const filterConfigs = [
    {
      filterKey: "docket",
      targetPath: ["filters", "match_docket_id"],
      valueProperty: "id",
      elseValue: "",
    },
    {
      filterKey: "organization",
      targetPath: ["filters", "match_author"],
      valueProperty: "label",
      elseValue: "",
    },
    {
      filterKey: "file_class",
      targetPath: ["filters", "match_file_class"],
      valueProperty: "label",
      elseValue: "",
    },
  ];

  filterConfigs.forEach(
    ({ filterKey, targetPath, valueProperty, elseValue }) => {
      const filters = filterTypeDict[filterKey];

      const setValue = (value: string) => {
        setNestedValue(new_file_filters, targetPath, value);
      };
      filterExtractionHelper(filters, elseValue, valueProperty, setValue);
    },
  );

  return new_file_filters;
};
const filterExtractionHelper = (
  filters: Filter[],
  valueProperty: string,
  elseValue: string,
  setValueFunc: (value: any) => void,
) => {
  if (filters && filters.length > 0) {
    if (filters.length > 1) {
      console.log(
        `Parameter '${filters[0].type}' shouldn't have more than one filter, ignoring extras`,
      );
    }
    const value = filters[0][valueProperty as keyof Filter]; // Seems like a total hack.
    setValueFunc(value);
  } else {
    setValueFunc(elseValue);
  }
};

// Helper to set values in nested object paths
const setNestedValue = (obj: any, path: string[], value: any) => {
  let current = obj;
  for (let i = 0; i < path.length - 1; i++) {
    current = current[path[i]];
  }
  current[path[path.length - 1]] = value;
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
        await mockFetchSuggestions(newQuery, input.pageContext),
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
