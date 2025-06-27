import { useEffect, useRef, useState } from "react";
import { useSearchState } from "@/lib/hooks/useSearchState";
import {
  GenericSearchInfo,
  GenericSearchType,
} from "@/lib/adapters/genericSearchCallback";

export function CommandKSearch() {
  const { resetToInitial } = useSearchState();
  const modalRef = useRef<HTMLDialogElement>(null);

  // Handle keyboard shortcut
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const isMac = navigator.userAgent.includes("Mac");
      const cmdK =
        (isMac && e.metaKey && e.key === "k") ||
        (!isMac && e.ctrlKey && e.key === "k");

      if (cmdK) {
        e.preventDefault();
        modalRef.current?.showModal();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  // Handle modal close events
  useEffect(() => {
    const dialog = modalRef.current;
    const handleClose = () => resetToInitial();

    dialog?.addEventListener("close", handleClose);
    return () => dialog?.removeEventListener("close", handleClose);
  }, [resetToInitial]);

  return (
    <>
      <dialog ref={modalRef} id="command-k-modal-default" className="modal">
        <div className="modal-box">
          <SearchCommand />
        </div>
        {/* DaisyUI modal backdrop click handler */}
        <form method="dialog" className="modal-backdrop">
          <button>close</button>
        </form>
      </dialog>
    </>
  );
}

function SearchCommand() {
  const {
    searchQuery,
    setSearchQuery,
    triggerSearch,
    getResultsCallback,
    searchTriggerIndicator,
    ...searchState
  } = useSearchState();

  const handleSearch = (query: string) => {
    triggerSearch({ query: query });
  };
  const searchInfo: GenericSearchInfo = {
    query: searchQuery,
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
            <button className="btn">Close</button>
          </form>
        </div>
      </div>

      <div className="grow overflow-auto">
        {/* <SearchResultsHomepageComponent */}
        {/*   searchInfo={searchInfo} */}
        {/*   isSearching={searchState.isSearching} */}
        {/*   reloadOnChange={searchTriggerIndicator} */}
        {/* /> */}
      </div>
    </div>
  );
}
