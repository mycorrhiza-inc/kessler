"use client";
import { useEffect, useState } from "react";
import LoadingSpinnerTimeout from "../styled-components/LoadingSpinnerTimeout";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import ErrorMessage from "../messages/ErrorMessage";
import NoResultsMessage from "../messages/NoResultsMessage";
import InfiniteScroll from "./InfiniteScroll";

const InfiniteScrollPlus = ({
  children,
  loadInitial,
  hasMore,
  getMore,
  reloadOnChange,
  dataLength,
}: {
  children: React.ReactNode;
  loadInitial: () => Promise<void>;
  getMore: () => Promise<void>;
  hasMore: boolean;
  reloadOnChange?: number;
  dataLength: number;
}) => {
  const [isTableReloading, setIsTableReloading] = useState(false);
  const [err, setErr] = useState("");
  const wrappedLoadInitial = async () => {
    setIsTableReloading(true);
    try {
      await loadInitial();
    } catch (error) {
      setErr((error as Error).message);
    }
    setIsTableReloading(false);
  };
  const [previousReload, setPreviousReload] = useState(reloadOnChange);
  useEffect(() => {
    if (previousReload != reloadOnChange) {
      wrappedLoadInitial();
      setPreviousReload(reloadOnChange);
    }
  }, [reloadOnChange])
  if (err != "") {
    return (
      <ErrorMessage
        reload
        message="We encountered an error while loading the data. Please try refreshing the page"
        error={err}
      />
    );
  }

  if (isTableReloading) {
    return <LoadingSpinner loadingText="Loading..." />;
  }

  if (!isTableReloading && dataLength == 0) {
    return <NoResultsMessage />;
  }
  return (
    <InfiniteScroll dataLength={dataLength} hasMore={hasMore} getMore={getMore}>
      {children}
    </InfiniteScroll>
  );
};

export default InfiniteScrollPlus;
