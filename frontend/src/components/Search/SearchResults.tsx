import { motion, AnimatePresence } from "framer-motion";
import SearchResultsClient from "./SearchResultsClient";
import { GenericSearchInfo } from "@/lib/adapters/genericSearchCallback";

export function SearchResultsComponent({
  searchInfo,
  reloadOnChange,
  isSearching,
}: {
  searchInfo: GenericSearchInfo;
  isSearching: boolean;
  reloadOnChange: number;
}) {
  return (
    <AnimatePresence mode="wait">
      {isSearching && (
        <motion.div
          key="search-results"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -20 }}
          transition={{ duration: 0.3 }}
          className="w-full"
        >
          <SearchResultsClient
            initialPage={0}
            initialData={[]}
            genericSearchInfo={searchInfo}
            reloadOnChange={reloadOnChange}
          />
        </motion.div>
      )}
    </AnimatePresence>
  );
}
