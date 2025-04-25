"use client";

import React, { useState } from "react";
import { MdKeyboardArrowDown } from "react-icons/md";

export default function StateSelector({
  states,
  selectedState,
  onSelect,
}: {
  states: string[];
  selectedState: string;
  onSelect: (state: string) => void;
}) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full px-4 py-3 border border-gray-300 rounded-lg flex items-center justify-between"
      >
        <span>{selectedState}</span>
        <MdKeyboardArrowDown className="text-xl" />
      </button>

      {isOpen && (
        <div className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-lg shadow-lg max-h-60 overflow-y-auto">
          {states.map((state) => (
            <div
              key={state}
              onClick={() => {
                onSelect(state);
                setIsOpen(false);
              }}
              className="px-4 py-2 hover:bg-gray-100 cursor-pointer"
            >
              {state}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
