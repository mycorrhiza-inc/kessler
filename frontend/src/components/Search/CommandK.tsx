import { useEffect, useRef } from "react";
import { useSearchState } from "@/lib/hooks/useSearchState";
import { SearchResultsComponent } from "./SearchResults";

export function CommandKSearch() {
  const { resetToInitial } = useSearchState();
  const modalRef = useRef<HTMLInputElement>(null);

  // Handle keyboard shortcut
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const isMac = navigator.userAgent.includes("Mac");
      const cmdK =
        (isMac && e.metaKey && e.key === "k") ||
        (!isMac && e.ctrlKey && e.key === "k");

      if (cmdK) {
        e.preventDefault();
        modalRef.current?.click(); // Toggle DaisyUI modal
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    // return () => window.removeEventListener("keydown", keydownHandler);
  }, []);

  return (
    <>
      {/* DaisyUI modal toggle */}
      <input
        ref={modalRef}
        type="checkbox"
        id="command-k-modal"
        className="modal-toggle"
      />

      {/* DaisyUI modal */}
      <div className="modal" role="dialog">
        <div className="modal-box w-11/12 max-w-5xl h-[80vh]">
          <SearchCommand
            onClose={() => {
              resetToInitial();
              modalRef.current?.click(); // Close modal
            }}
          />
        </div>

        {/* Click outside to close */}
        <label className="modal-backdrop" htmlFor="command-k-modal" />
      </div>
    </>
  );
}

function SearchCommand({ onClose }: { onClose: () => void }) {
  const {
    searchQuery,
    setSearchQuery,
    triggerSearch,
    getResultsCallback,
    searchTriggerIndicator,
    resetToInitial,
    ...searchState
  } = useSearchState();

  const handleSearch = (query: string) => {
    setSearchQuery(query);
    triggerSearch();
  };

  return (
    <div className="flex flex-col h-full">
      <div className="flex gap-2 mb-4">
        <input
          value={searchQuery}
          onChange={(e) => handleSearch(e.target.value)}
          placeholder="Search..."
          className="input input-bordered flex-grow"
          autoFocus
        />
        <button
          className="btn btn-ghost"
          onClick={onClose}
          aria-label="Close search"
        >
          ✕
        </button>
      </div>

      <div className="flex-grow overflow-auto">
        <SearchResultsComponent
          isSearching={searchState.isSearching}
          reloadOnChange={searchTriggerIndicator}
          searchGetter={getResultsCallback}
        />
      </div>
    </div>
  );
}
