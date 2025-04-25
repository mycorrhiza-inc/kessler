import React, { useState } from "react";
import { ChevronDownIcon } from "@heroicons/react/24/solid";

interface FilterOption {
  label: string;
  value: string;
}

interface FilterDropdownProps {
  title: string;
  options: FilterOption[];
  selectedOptions: string[];
  onSelectOption: (value: string) => void;
}

const FilterDropdown: React.FC<FilterDropdownProps> = ({
  title,
  options,
  selectedOptions,
  onSelectOption,
}) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-full flex justify-between items-center px-4 py-2 bg-gray-100 rounded-md"
      >
        {title}
        <ChevronDownIcon
          className={`h-5 w-5 transition-transform ${isOpen ? "rotate-180" : ""}`}
        />
      </button>

      {isOpen && (
        <div className="absolute z-10 w-full mt-1 bg-white border rounded-md shadow-lg max-h-60 overflow-y-auto">
          {options.map((option) => (
            <div
              key={option.value}
              onClick={() => {
                onSelectOption(option.value);
                setIsOpen(false);
              }}
              className={`
                px-4 py-2 cursor-pointer hover:bg-gray-100
                ${selectedOptions.includes(option.value) ? "bg-blue-100" : ""}
              `}
            >
              {option.label}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default FilterDropdown;
