"use client";
import React, { useEffect, useRef, useState } from "react";
import { AngleDownIcon, AngleUpIcon } from "../Icons";
import { subdividedHueFromSeed } from "../Tables/TextPills";
import {
  QueryDataFile,
  InheritedFilterValues,
  initialFiltersFromInherited,
} from "@/lib/filters";
import {
  FileSearchBoxProps,
  Filter,
  PageContextMode,
  SearchBoxInputProps,
  Suggestion,
} from "@/lib/types/SearchTypes";
import { ConvoSearchRequestData } from "../LookupPages/SearchRequestData";
import { mockFetchSuggestions } from "./SearchSuggestions";

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

const displayTypeDict = {
  nypuc_docket_industry: "Docket Industry",
};
const getDisplayType = (val: string): string => {
  if (val in displayTypeDict) {
    return displayTypeDict[val as keyof typeof displayTypeDict];
  }
  return val
    .split("_")
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
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
    return subdividedHueFromSeed(filter.label);
    // return "oklch(90% 0.1 80)";
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
                {getDisplayType(filter.type)}: {filter.label}
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
      fileProps.setSearchData(
        generateFileFiltersFromFilterList(
          props.inheritedFileFilters,
          filterTypeDict,
        ),
      );
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
  var convoSearchData: ConvoSearchRequestData = {};
  const setQuery = (value: string) => {
    convoSearchData.query = value;
  };
  const querySchema: filterExtractionSchema = {
    filters: filterTypeDict.text,
    valueProperty: "label",
    elseValue: "",
    setValueFunc: setQuery,
  };
  filterExtractionHelper(querySchema);
  const setIndustry = (value: string) => {
    convoSearchData.industry_type = value;
  };
  const industrySchema: filterExtractionSchema = {
    filters: filterTypeDict.nypuc_docket_industry,
    valueProperty: "label",
    elseValue: "",
    setValueFunc: setIndustry,
  };
  filterExtractionHelper(industrySchema);
  return convoSearchData;
};

const generateFileFiltersFromFilterList = (
  inheritedFileFilters: InheritedFilterValues,
  filterTypeDict: { [key: string]: Filter[] },
) => {
  const new_file_filters_metadata =
    initialFiltersFromInherited(inheritedFileFilters);
  const new_file_filters: QueryDataFile = {
    filters: new_file_filters_metadata,
    query: getTextQueryFromFilterList(filterTypeDict),
  };

  const filterConfigs = [
    {
      filterKey: "docket",
      targetPath: ["filters", "match_docket_id"],
      valueProperty: "id",
      elseValue: new_file_filters_metadata.match_docket_id,
    },
    {
      filterKey: "organization",
      targetPath: ["filters", "match_author"],
      valueProperty: "label",
      elseValue: new_file_filters_metadata.match_author,
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
      filterExtractionHelper({
        filters: filters,
        valueProperty: valueProperty,
        elseValue: elseValue,
        setValueFunc: setValue,
      });
    },
  );

  return new_file_filters;
};

interface filterExtractionSchema {
  filters: Filter[];
  valueProperty: string;
  elseValue: string;
  setValueFunc: (value: any) => void;
}
const filterExtractionHelper = (schema: filterExtractionSchema) => {
  const filters = schema.filters;
  const setValueFunc = schema.setValueFunc;
  const elseValue = schema.elseValue;
  const valueProperty = schema.valueProperty;
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

const SearchBox = ({
  input,
  ShowAdvancedSearch,
}: {
  input: SearchBoxInputProps;
  ShowAdvancedSearch?: boolean;
}) => {
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
    <div className="p-12 max-w-xl mx-auto">
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
            {ShowAdvancedSearch && <AdvancedSearch />}
          </div>

          {/* Suggestions dropdown - positioned relative to search container */}
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
                      <span className={`text-sm font-medium text-secondary`}>
                        {getDisplayType(suggestion.type)}:
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
