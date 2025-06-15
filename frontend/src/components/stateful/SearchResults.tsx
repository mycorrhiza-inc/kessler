import React from 'react';
import { GenericSearchInfo } from '@/lib/adapters/genericSearchCallback';

export interface SearchResultsHomepageComponentProps {
  searchInfo: GenericSearchInfo;
  isSearching: boolean;
  reloadOnChange: number;
}

export const SearchResultsHomepageComponent: React.FC<SearchResultsHomepageComponentProps> = ({
  searchInfo,
  isSearching,
  reloadOnChange,
}) => {
  // TODO: Implement the search results UI
  return null;
};
