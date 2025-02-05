import { QueryDataFile } from "../filters";

// Mock API call
export type Suggestion = {
  id: string;
  type: string;
  label: string;
  value: string;
};

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
}
export interface OrgSearchBoxProps {
  pageContext: PageContextMode.Organizations;
  setSearchQuery: React.Dispatch<React.SetStateAction<string>>;
}
export interface DocketSearchBoxProps {
  pageContext: PageContextMode.Conversations;
  setSearchQuery: React.Dispatch<React.SetStateAction<string>>;
}

export type SearchBoxInputProps =
  | FileSearchBoxProps
  | OrgSearchBoxProps
  | DocketSearchBoxProps;

export type FilterTypeDict = { [key: string]: Filter[] };
