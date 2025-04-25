import React from "react";
import { MdClose } from "react-icons/md";

interface SelectedFiltersProps {
  states: string[];
  authors: string[];
  dockets: string[];
  onRemoveFilter: (
    type: "states" | "authors" | "dockets",
    value: string,
  ) => void;
}

const SelectedFilters: React.FC<SelectedFiltersProps> = ({
  states,
  authors,
  dockets,
  onRemoveFilter,
}) => {
  // Combine all filters with their types
  const allFilters = [
    ...states.map((state) => ({
      type: "states" as const,
      value: state,
      color: "bg-blue-100 text-blue-800",
    })),
    ...authors.map((author) => ({
      type: "authors" as const,
      value: author,
      color: "bg-pink-100 text-pink-800",
    })),
    ...dockets.map((docket) => ({
      type: "dockets" as const,
      value: docket,
      color: "bg-green-100 text-green-800",
    })),
  ];

  if (allFilters.length === 0) return null;

  return (
    <div className="selected-filters space-y-2">
      <h3 className="text-sm font-semibold">Selected Filters</h3>
      <div className="flex flex-wrap gap-2">
        {allFilters.map(({ type, value, color }) => (
          <div
            key={`${type}-${value}`}
            className={`
              flex items-center px-2 py-1 rounded-full text-xs
              ${color}
            `}
          >
            {value}
            <button
              onClick={() => onRemoveFilter(type, value)}
              className="ml-2 hover:bg-opacity-75 rounded-full"
            >
              <MdClose className="h-4 w-4" />
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default SelectedFilters;
