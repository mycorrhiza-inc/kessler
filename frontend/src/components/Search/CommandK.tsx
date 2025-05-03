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
        console.log("Entered command k modal");
        // modalRef.current?.click(); // Toggle DaisyUI modal
        document.getElementById("command-k-modal-default").showModal();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    // Not sure what this does
    // return () => window.removeEventListener("keydown", keydownHandler);
  }, []);

  return (
    <>
      <dialog id="command-k-modal-default" className="modal">
        <div className="modal-box">
          <SearchCommand
            onClose={() => {
              resetToInitial();
              modalRef.current?.click(); // Close modal
            }}
          />
        </div>
      </dialog>
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
      <div className="flex-row gap-2 mb-4">
        <input
          value={searchQuery}
          onChange={(e) => handleSearch(e.target.value)}
          placeholder="Search..."
          className="input input-bordered grow"
          autoFocus
        />

        <div className="modal-action">
          <form method="dialog">
            {/* if there is a button in form, it will close the modal */}
            <button className="btn">Close</button>
          </form>
        </div>
      </div>

      <div className="grow overflow-auto">
        <SearchResultsComponent
          isSearching={searchState.isSearching}
          reloadOnChange={searchTriggerIndicator}
          searchGetter={getResultsCallback}
        />
      </div>
    </div>
  );
}
