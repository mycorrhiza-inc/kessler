"use client";
import React, { useState, useMemo, useEffect } from "react";
import { getSearchResults } from "@/lib/requests/search";
import InfiniteScrollPlus from "@/components/InfiniteScroll/InfiniteScroll";
import SearchBox from "@/components/Search/SearchBox";
import { FileSearchBoxProps, PageContextMode } from "@/lib/types/SearchTypes";
import { adaptFilingToCard } from "@/lib/adapters/genericCardAdapters";
import Card, { CardSize } from "@/components/NewSearch/GenericResultCard";
import {
  QueryDataFile,
  InheritedFilterValues,
  initialFiltersFromInherited,
} from "@/lib/filters";
import { Filing } from "@/lib/types/FilingTypes";

interface FileSearchViewNewProps {
  /** Initial results for SSR */
  initialData?: Filing[];
  /** Initial page index for SSR */
  initialPage?: number;
  inheritedFilters: InheritedFilterValues;
  DocketColumn?: boolean;
}

const FileSearchViewNew: React.FC<FileSearchViewNewProps> = ({
  initialData = [],
  initialPage = 2,
  inheritedFilters,
}) => {};

export default FileSearchViewNew;
