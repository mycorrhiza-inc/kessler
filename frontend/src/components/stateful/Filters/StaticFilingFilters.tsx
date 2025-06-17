"use client";

import { encodeUrlParams, TypedUrlParams } from "@/lib/types/url_params";
import { useState } from "react";
import { DynamicMultiSelect } from "./FilterMultiSelect";
import { useRouter } from "next/router";
import { ne } from "@faker-js/faker";


// export interface FilterFieldDefinition {
//   id: string;
//   displayName: string;
//   description: string;
//   placeholder?: string;
//   options?: Array<{
//     value: string;
//     label: string;
//     disabled?: boolean;
//   }>;
// }
//
// export interface DynamicMultiSelectProps {
//   fieldDefinition: FilterFieldDefinition;
//   value: string;
//   onChange: (value: string) => void;
//   onFocus?: () => void;
//   onBlur?: () => void;
//   disabled?: boolean;
//   className?: string;
// }
// export interface DynamicSingleSelectProps {
//   fieldDefinition: FilterFieldDefinition;
//   value: string;
//   onChange: (value: string) => void;
//   onFocus?: () => void;
//   onBlur?: () => void;
//   disabled?: boolean;
//   className?: string;
//   allowClear?: boolean;
//   defaultValue?: string;
//   dynamicWidth?: boolean;
//   minWidth?: string;
//   maxWidth?: string;
// }

// Here are some type definitions from other files:
export enum FilterType {
  Single,
  Multi
}
export interface MinimalFilterDefinition {
  filterType: FilterType
  id: string,
  displayName: string;
  description?: string;
  placeholder?: string;
}

export default function HardcodedFileFilters({ urlParams, baseUrl }: { urlParams: TypedUrlParams, baseUrl: string }) {
  // add filter types for 
  // MultiSelect:
  // author_id
  // file_extension
  // SingleSelect:
  // conversation_id
  // file_class
  //
  //
  const fileFilterInfo: MinimalFilterDefinition[] = [];

  return <HardCodedFiltersFromInfo urlParams={urlParams} baseUrl={baseUrl} hardcodedFilterInfo={fileFilterInfo} />
}


export function HardCodedFiltersFromInfo({ urlParams, baseUrl, hardcodedFilterInfo }: { urlParams: TypedUrlParams, baseUrl: string, hardcodedFilterInfo: MinimalFilterDefinition[] }) {
  const initial_filters = urlParams.queryData.filters || {}
  const [filter_values, setFilterValues] = useState(initial_filters)
  const router = useRouter()
  const executeSearch = () => {
    const new_params: TypedUrlParams = {
      paginationData: {},
      queryData: {
        ...urlParams.queryData, filters: filter_values
      }
    }
    const endpoint = baseUrl + encodeUrlParams(new_params)
    router.push(endpoint)
  }
  // Figure out all the state management to get the filters to dynamically update the filter values whenever anything changes.
  return <>
    {/* Map the hardcoded filters into a list of dynamic filters */}
  </>
}
