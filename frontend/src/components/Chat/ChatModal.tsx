import {
  InheritedFilterValues,
  initialFiltersFromInherited,
} from "@/lib/filters";
import { useMemo, useState } from "react";
import Modal from "../styled-components/Modal";
import ChatBoxInternals from "./ChatBoxInternals";
import clsx from "clsx";

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

  return (
    <>
      <Modal setOpen={setOpen} open={open}>
        <ChatBoxInternals
          setCitations={setCitations}
          ragFilters={filterState}
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
