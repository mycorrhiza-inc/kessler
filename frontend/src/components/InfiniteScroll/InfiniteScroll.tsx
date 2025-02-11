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
  const [isTableReloading, setIsTableReloading] = useState(false);
  const wrappedLoadInitial = async () => {
    setIsTableReloading(true);
    await loadInitial();
    setIsTableReloading(false);
  };
  useEffect(() => {
    wrappedLoadInitial();
  }, [reloadOnChangeObj]);

  if (isTableReloading) {
    return <LoadingSpinner loadingText="Reloading..." />;
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
