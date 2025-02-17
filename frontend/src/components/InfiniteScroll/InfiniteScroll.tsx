"use client";
import { useEffect, useState } from "react";
import LoadingSpinnerTimeout from "../styled-components/LoadingSpinnerTimeout";
import InfiniteScroll from "react-infinite-scroll-component";
import LoadingSpinner from "../styled-components/LoadingSpinner";
import ErrorMessage from "../ErrorMessage";

const InfiniteScrollPlus = ({
  id,
  children,
  loadInitial,
  hasMore,
  getMore,
  reloadOnChangeObj,
  dataLength,
}: {
  id?: string;
  children: React.ReactNode;
  loadInitial: () => Promise<void>;
  getMore: () => Promise<void>;
  hasMore?: any;
  reloadOnChangeObj?: any;
  dataLength: number;
}) => {
  const [isTableReloading, setIsTableReloading] = useState(true);
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
  useEffect(() => {
    wrappedLoadInitial();
  }, [reloadOnChangeObj]);
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
