"use client";

import { encodeUrlParams, TypedUrlParams, UrlQueryParams } from "@/lib/types/url_params";
import { useState } from "react";
import { DynamicMultiSelect } from "./FilterMultiSelect";
import { DynamicSingleSelect } from "./FilterSingleSelect";
import { useRouter } from "next/navigation";
import {
  OrganizationsAutocompleteList,
  ConversationsAutocompleteList,
} from "./HardcodedAutocompletes";

// Minimal filter definition used to render dynamic filters
export enum FilterType {
  Single,
  Multi,
}
export interface MinimalFilterDefinition {
  filterType: FilterType;
  id: string;
  displayName: string;
  description?: string;
  placeholder?: string;
}

/**
 * HardcodedFileFilters defines which filters to show for filing.
 * - MultiSelect: author_id, file_extension
 * - SingleSelect: conversation_id, file_class
 */
export default function HardcodedFileFilters({
  urlParams,
  baseUrl,
  className = "",
}: {
  urlParams: UrlQueryParams;
  baseUrl: string;
  className?: string;
}) {
  const fileFilterInfo: MinimalFilterDefinition[] = [
    {
      filterType: FilterType.Multi,
      id: "author_id",
      displayName: "Author",
      description: "Filter by author",
      placeholder: "Select authors",
    },
    {
      filterType: FilterType.Multi,
      id: "file_extension",
      displayName: "File Extension",
      description: "Filter by file extension",
      placeholder: "Select file extensions",
    },
    {
      filterType: FilterType.Single,
      id: "conversation_id",
      displayName: "Conversation",
      description: "Filter by conversation",
      placeholder: "Select a conversation",
    },
    {
      filterType: FilterType.Single,
      id: "file_class",
      displayName: "File Class",
      description: "Filter by file class",
      placeholder: "Select file class",
    },
  ];

  return (
    <HardCodedFiltersFromInfo
      urlQueryParams={urlParams}
      baseUrl={baseUrl}
      hardcodedFilterInfo={fileFilterInfo}
      className={className}
    />
  );
}

export function HardCodedFiltersFromInfo({
  urlQueryParams,
  baseUrl,
  hardcodedFilterInfo,
  className = "",
}: {
  urlQueryParams: UrlQueryParams;
  baseUrl: string;
  hardcodedFilterInfo: MinimalFilterDefinition[];
  className?: string;
}) {
  const router = useRouter();
  const initialFilters = urlQueryParams.filters || {};
  const [filterValues, setFilterValues] = useState<Record<string, string>>(initialFilters);
  const [openDropdowns, setOpenDropdowns] = useState<Set<string>>(new Set());

  // Update filter and navigate
  const handleFilterChange = (id: string) => (value: string) => {
    const updated = { ...filterValues, [id]: value };
    setFilterValues(updated);

    const newParams: TypedUrlParams = {
      paginationData: {},
      namespace: "",
      queryData: { ...urlQueryParams, filters: updated },
    };
    const endpoint = baseUrl + encodeUrlParams(newParams);
    router.push(endpoint);
  };

  // Track which dropdowns are open to expand container accordingly
  const handleDropdownStateChange = (id: string, isOpen: boolean) => {
    setOpenDropdowns(prev => {
      const newSet = new Set(prev);
      if (isOpen) {
        newSet.add(id);
      } else {
        newSet.delete(id);
      }
      return newSet;
    });
  };

  // Helper to map filter IDs to static autocomplete options
  const getStaticOptions = (id: string) => {
    switch (id) {
      case "author_id":
        return OrganizationsAutocompleteList.map((org) => ({ value: org.value, label: org.label }));
      case "conversation_id":
        return ConversationsAutocompleteList.map((c) => ({ value: c.value, label: c.label }));
      default:
        return [];
    }
  };

  // Calculate additional height needed for open dropdowns
  const dropdownHeight = 300; // Standard dropdown height
  const additionalHeight = openDropdowns.size * dropdownHeight;
  const hasOpenDropdowns = openDropdowns.size > 0;

  return (
    <div
      className={`filter-container flex flex-col space-y-4 transition-all duration-300 ease-in-out overflow-visible ${className}`}
      style={{
        paddingBottom: hasOpenDropdowns ? `${additionalHeight}px` : '0px'
      }}
    >
      {hardcodedFilterInfo.map((def) => {
        const value = filterValues[def.id] || "";
        const options = getStaticOptions(def.id);
        const fieldDef = {
          id: def.id,
          displayName: def.displayName,
          description: def.description || "",
          placeholder: def.placeholder || "",
          options,
        };

        const isOpen = openDropdowns.has(def.id);

        return (
          <div key={def.id} className="flex-shrink-0 relative">
            {def.filterType === FilterType.Multi ? (
              <DynamicMultiSelect
                className=""
                fieldDefinition={fieldDef}
                value={value}
                onChange={handleFilterChange(def.id)}
                onDropdownStateChange={(open) => handleDropdownStateChange(def.id, open)}
              />
            ) : (
              <DynamicSingleSelect
                className=""
                fieldDefinition={fieldDef}
                value={value}
                onChange={handleFilterChange(def.id)}
                onDropdownStateChange={(open) => handleDropdownStateChange(def.id, open)}
              />
            )}
          </div>
        );
      })}
    </div>
  );
}
