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
}: {
  urlParams: UrlQueryParams;
  baseUrl: string;
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
    />
  );
}

export function HardCodedFiltersFromInfo({
  urlQueryParams,
  baseUrl,
  hardcodedFilterInfo,
}: {
  urlQueryParams: UrlQueryParams;
  baseUrl: string;
  hardcodedFilterInfo: MinimalFilterDefinition[];
}) {
  const router = useRouter();
  const initialFilters = urlQueryParams.filters || {};
  const [filterValues, setFilterValues] = useState<Record<string, string>>(initialFilters);

  // Update filter and navigate
  const handleFilterChange = (id: string) => (value: string) => {
    const updated = { ...filterValues, [id]: value };
    setFilterValues(updated);

    const newParams: TypedUrlParams = {
      paginationData: {},
      queryData: { ...urlQueryParams, filters: updated },
    };
    const endpoint = baseUrl + encodeUrlParams(newParams);
    router.push(endpoint);
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

  return (
    <div className="filter-container flex flex-col space-y-24">
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
        return def.filterType === FilterType.Multi ? (
          <DynamicMultiSelect
            className="z-index-1"
            key={def.id}
            fieldDefinition={fieldDef}
            value={value}
            onChange={handleFilterChange(def.id)}
          />
        ) : (
          <DynamicSingleSelect
            className="z-index-1"
            key={def.id}
            fieldDefinition={fieldDef}
            value={value}
            onChange={handleFilterChange(def.id)}
          />
        );
      })}
    </div>
  );
}
