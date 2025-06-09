"use client";
import {
  InheritedFilterValues,
  initialFiltersFromInherited,
} from "@/lib/filters";
import { useMemo, useState } from "react";
import Modal from "../styled-components/Modal";
// import {
//   ChatBoxInternalsState,
//   ChatBoxInternalsStateless,
//   initialChatState,
// } from "./ChatBoxInternals";
import { useKesslerStore } from "@/lib/_store";

export const ExperimentalChatModalClickDiv = ({
  inheritedFilters,
  children,
  className,
}: {
  inheritedFilters: InheritedFilterValues;
  children?: React.ReactNode;
  className?: string;
}) => {
  const globalStore = useKesslerStore();
  const experimentalFilters = globalStore.experimentalFeaturesEnabled;
  return (
    <>
      {experimentalFilters && (
        <ChatModalClickDiv
          inheritedFilters={inheritedFilters}
          className={className}
        >
          {children}
        </ChatModalClickDiv>
      )}
    </>
  );
};

export const ChatModalClickDiv = ({
  inheritedFilters,
  children,
  className,
}: {
  inheritedFilters: InheritedFilterValues;
  children?: React.ReactNode;
  className?: string;
}) => {
  // const initialFilterState = useMemo(() => {
  //   return initialFiltersFromInherited(inheritedFilters);
  // }, [inheritedFilters]);
  // const filterState = initialFilterState; // Maybe add a state hook for the llm setting its own filters or something.
  const [open, setOpen] = useState(false);
  const [citations, setCitations] = useState<any[]>([]);
  const toggleOpen = () => setOpen((prev) => !prev);
  // const [chatState, setChatState] =
  //   useState<ChatBoxInternalsState>(initialChatState);

  return (
    <>
      <Modal setOpen={setOpen} open={open}>
        <div></div>
        {/* <ChatBoxInternalsStateless */}
        {/*   setCitations={setCitations} */}
        {/*   ragFilters={filterState} */}
        {/*   chatState={chatState} */}
        {/*   setChatState={setChatState} */}
      </Modal>
      <div onClick={toggleOpen} className={className}>
        {children}
      </div>
    </>
  );
};

export const ChatModalTestButton = ({
  inheritedFilters,
}: {
  inheritedFilters: InheritedFilterValues;
}) => {
  return (
    <ChatModalClickDiv inheritedFilters={inheritedFilters}>
      <button className="btn btn-primary">Test Chat</button>
    </ChatModalClickDiv>
  );
};
