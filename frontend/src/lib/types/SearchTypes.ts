import { ConvoSearchRequestData } from "@/components/LookupPages/SearchRequestData";
import { InheritedFilterValues, QueryDataFile } from "../filters";
import { Range } from "react-date-range";

// Mock API call

enum BasicSuggestions {
  Text = "text",
  Author = "author",
  Docket = "docket",
  Class = "class",
  Extension = "extension",
}

enum DateSuggestionType {
  Date = "date",
}

type SuggestionTypes = BasicSuggestions | DateSuggestionType;

export interface BasicSuggestion {
  id: string;
  type: SuggestionTypes;
  label: string;
  value: any;
}
export interface DateSuggestion extends BasicSuggestion {
  id: string;
  type: DateSuggestionType;
  label: string;
  value: Range;
}
export interface TextSuggestion extends BasicSuggestion {
  id: "";
}

export type Suggestion = BasicSuggestion | DateSuggestion;

export type Filter = {
  id: string;
  type: string;
  label: string;
  exclude?: boolean;
  excludable: boolean;
};
export enum PageContextMode {
  Files,
  Organizations,
  Conversations,
}
export interface FileSearchBoxProps {
  pageContext: PageContextMode.Files;
  setSearchData: React.Dispatch<React.SetStateAction<QueryDataFile>>;
  inheritedFileFilters: InheritedFilterValues;
}
export interface OrgSearchBoxProps {
  pageContext: PageContextMode.Organizations;
  setSearchQuery: React.Dispatch<React.SetStateAction<string>>;
}
export interface DocketSearchBoxProps {
  pageContext: PageContextMode.Conversations;
  setSearchData: React.Dispatch<React.SetStateAction<ConvoSearchRequestData>>;
}

export type SearchBoxInputProps =
  | FileSearchBoxProps
  | OrgSearchBoxProps
  | DocketSearchBoxProps;

export type FilterTypeDict = { [key: string]: Filter[] };
