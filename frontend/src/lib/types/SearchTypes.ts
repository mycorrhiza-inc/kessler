import { ConvoSearchRequestData } from "@/components/LookupPages/SearchRequestData";
import { InheritedFilterValues, QueryDataFile } from "../filters";
import { Range } from "react-date-range";
import { InputType } from "@/components/Filters/FiltersInfo";

// Mock API call


export interface BasicSuggestion {
  id: string;
  type: InputType;
  label: string;
  value?: any;
  exclude?: boolean;
  excludable?: boolean;
}
export interface DateSuggestion extends BasicSuggestion {
  id: string;
  type: InputType.Date;
  label: string;
  value: Range;
}
export interface TextSuggestion extends BasicSuggestion {
  id: "";
}

export type Suggestion = BasicSuggestion ;

export interface Filter extends BasicSuggestion {
  id: string;
  type: InputType;
  label: string;
  value?: any;
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
