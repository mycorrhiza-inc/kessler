"use client";
import React, { useEffect, useState, useCallback, useRef } from "react";
import ErrorMessage from "../messages/ErrorMessage";
import NoResultsMessage from "../messages/NoResultsMessage";
import LoadingSpinner from "../styled-components/LoadingSpinner";

interface InfiniteScrollProps {
  children: React.ReactNode;
  getMore: () => Promise<void>;
  hasMore: boolean;
  dataLength: number;
  /**
   * Distance in pixels from bottom to trigger loading more (default: 100)
   */
  threshold?: number;
  containerHeight?: number;
}

const InfiniteScroll: React.FC<InfiniteScrollProps> = ({
  children,
  getMore,
  hasMore,
  dataLength,
  threshold = 100,
  containerHeight = 800,
}) => {
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const loadingRef = useRef(false);
  const containerRef = useRef<HTMLDivElement>(null);

  const handleScroll = useCallback(() => {
    if (!containerRef.current || loadingRef.current || loadingMore || !hasMore)
      return;

    const { scrollTop, scrollHeight, clientHeight } = containerRef.current;
    const scrollPosition = scrollTop + clientHeight;
    const bottomPosition = scrollHeight - threshold;

    if (scrollPosition >= bottomPosition) {
      loadingRef.current = true;
      setLoadingMore(true);
      getMore()
        .catch((err) => setError(err.message || "unknown error encountered"))
        .finally(() => {
          setLoadingMore(false);
          loadingRef.current = false;
        });
    }
  }, [getMore, hasMore, loadingMore, threshold]);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    container.addEventListener("scroll", handleScroll, { passive: true });
    return () => container.removeEventListener("scroll", handleScroll);
  }, [handleScroll]);

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
    <div
      ref={containerRef}
      style={{
        height: `${containerHeight}px`,
        overflowY: "auto",
        position: "relative",
      }}
    >
      {children}
      {loadingMore && (
        <div style={{ padding: "20px 0" }}>
          <LoadingSpinner loadingText="Loading more..." />
        </div>
      )}
    </div>
  );
};

export default InfiniteScroll;
