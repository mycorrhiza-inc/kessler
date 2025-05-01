import { motion, AnimatePresence } from "framer-motion";

interface SearchResultsProps {
  isSearching: boolean;
  isLoading: boolean;
  error: string | null;
  searchResults: any[];
  children: React.ReactNode;
}

export function SearchResultsWrapper({
  isSearching,
  isLoading,
  error,
  searchResults,
  children,
}: SearchResultsProps) {
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
          {isLoading ? (
            <div className="flex justify-center p-8">
              <span className="loading loading-spinner loading-lg"></span>
            </div>
          ) : error ? (
            <div className="alert alert-error my-4">
              <span>{error}</span>
            </div>
          ) : searchResults.length > 0 ? (
            children
          ) : (
            <div className="text-center p-8">
              <h3 className="text-lg font-medium">No results found</h3>
              <p className="text-sm opacity-70 mt-1">
                Try adjusting your search terms
              </p>
            </div>
          )}
        </motion.div>
      )}
    </AnimatePresence>
  );
}
