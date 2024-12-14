"use client";
import {
  InheritedFilterValues,
  initialFiltersFromInherited,
} from "@/lib/filters";
import { useMemo, useState } from "react";
import Modal from "../styled-components/Modal";
import {
  ChatBoxInternalsState,
  ChatBoxInternalsStateless,
  initialChatState,
} from "./ChatBoxInternals";

export const ChatModalClickDiv = ({
  inheritedFilters,
  children,
  className,
}: {
  inheritedFilters: InheritedFilterValues;
  children?: React.ReactNode;
  className?: string;
}) => {
  const initialFilterState = useMemo(() => {
    return initialFiltersFromInherited(inheritedFilters);
  }, [inheritedFilters]);
  const filterState = initialFilterState; // Maybe add a state hook for the llm setting its own filters or something.
  const [open, setOpen] = useState(false);
  const [citations, setCitations] = useState<any[]>([]);
  const clickButton = () => setOpen((prev) => !prev);
  const [chatState, setChatState] =
    useState<ChatBoxInternalsState>(initialChatState);

  return (
    <>
      <Modal setOpen={setOpen} open={open}>
        <ChatBoxInternalsStateless
          setCitations={setCitations}
          ragFilters={filterState}
          chatState={chatState}
          setChatState={setChatState}
        />
      </Modal>
      <div onClick={clickButton} className={className}>
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
