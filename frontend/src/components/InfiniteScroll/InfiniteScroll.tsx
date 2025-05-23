"use client";
import React, { useEffect, useState, useCallback, useRef } from "react";
import LoadingSpinnerTimeout from "../styled-components/LoadingSpinnerTimeout";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import ErrorMessage from "../messages/ErrorMessage";
import NoResultsMessage from "../messages/NoResultsMessage";

interface InfiniteScrollProps {
  children: React.ReactNode;
  getMore: () => Promise<void>;
  hasMore: boolean;
  dataLength: number;
  /**
   * Distance in pixels from bottom to trigger loading more (default: 100)
   */
  threshold?: number;
}

const InfiniteScroll: React.FC<InfiniteScrollProps> = ({
  children,
  getMore,
  hasMore,
  dataLength,
  threshold = 100,
}) => {
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const loadingRef = useRef(false);

  // Throttled scroll handler
  const handleScroll = useCallback(() => {
    if (loadingRef.current || loadingMore || !hasMore) {
      return;
    }
    const scrollPosition = window.innerHeight + window.scrollY;
    const bottomPosition = document.documentElement.scrollHeight - threshold;
    if (scrollPosition >= bottomPosition) {
      loadingRef.current = true;
      setLoadingMore(true);
      getMore()
        .catch((err) => setError((err as Error).message))
        .finally(() => {
          setLoadingMore(false);
          loadingRef.current = false;
        });
    }
  }, [getMore, hasMore, loadingMore, threshold]);

  // Attach and clean up scroll listener
  useEffect(() => {
    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  }, [handleScroll]);

  // Render error or loading or empty states
  if (error) {
    return (
      <ErrorMessage
        reload
        message={`We encountered an error while loading the data. Please try refreshing the page: ${JSON.stringify(error)}`}
        error={error}
      />
    );
  }
  if (dataLength === 0) {
    return <NoResultsMessage />;
  }

  return (
    <>
      {children}
      {loadingMore && (
        <LoadingSpinnerTimeout
          timeoutSeconds={10}
          loadingText="Loading more..."
        />
      )}
    </>
  );
};

export default InfiniteScroll;
