import { InheritedFilterValues, QueryDataFile } from "../filters";
import { Range } from "react-date-range";

// Mock API call

export enum FilterInputType {
  Text = "text",
  Select = "select",
  Datalist = "datalist",
  Date = "date",
  NotShown = "not_shown",
  OrgMultiselect = "org_multiselect",
  ConvoMultiselect = "convo_multiselect",
  Organization = "organization",
  Docket = "docket",
  FileClass = "file_class",
  FileExtension = "extension",
  NYDocket = "nypuc_docket_industry",
}

export interface ConvoSearchRequestData {
  query?: string;
  industry_type?: string;
}
export interface BasicSuggestion {
  id: string;
  type: FilterInputType;
  label: string;
  value?: any;
  exclude?: boolean;
  excludable?: boolean;
}
export interface DateSuggestion extends BasicSuggestion {
  id: string;
  type: FilterInputType.Date;
  label: string;
  value: Range;
}

export type Suggestion = BasicSuggestion;

export interface Filter extends BasicSuggestion {
  id: string;
  type: FilterInputType;
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
