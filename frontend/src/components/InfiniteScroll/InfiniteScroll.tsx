"use client";
import { useEffect, useState } from "react";
import LoadingSpinnerTimeout from "../styled-components/LoadingSpinnerTimeout";
import InfiniteScroll from "react-infinite-scroll-component";
import LoadingSpinner from "../styled-components/LoadingSpinner";

const InfiniteScrollPlus = ({
  children,
  loadInitial,
  hasMore,
  getMore,
  reloadOnChangeObj,
  dataLength,
}: {
  children: React.ReactNode;
  loadInitial: () => Promise<void>;
  getMore: () => Promise<void>;
  hasMore?: any;
  reloadOnChangeObj?: any;
  dataLength: number;
}) => {
  const [isTableReloading, setIsTableReloading] = useState(true);
  const [hasErrored, setHasErrored] = useState(false);
  const wrappedLoadInitial = async () => {
    setIsTableReloading(true);
    try {
      await loadInitial();
    } catch {
      setHasErrored(true);
    }
    setIsTableReloading(false);
  };
  useEffect(() => {
    wrappedLoadInitial();
  }, [reloadOnChangeObj]);
  if (hasErrored) {
    return (
      <div className="flex flex-col items-center justify-center p-8 m-4 rounded-lg bg-error/10 text-error">
        <div className="text-5xl mb-4">⚠️</div>
        <h3 className="text-xl font-bold mb-2">Oops! Something went wrong</h3>
        <p className="text-center mb-4">
          We encountered an error while loading the data. Please try refreshing
          the page.
        </p>
        <button
          onClick={() => window.location.reload()}
          className="btn btn-error btn-outline"
        >
          Reload Page
        </button>
      </div>
    );
  }

  if (isTableReloading) {
    return <LoadingSpinner loadingText="Loading..." />;
  }
  return (
    <InfiniteScroll
      dataLength={dataLength}
      hasMore={true}
      next={getMore}
      loader={
        <LoadingSpinnerTimeout timeoutSeconds={10} loadingText="Loading..." />
      }
    >
      {children}
    </InfiniteScroll>
  );
};

export default InfiniteScrollPlus;
