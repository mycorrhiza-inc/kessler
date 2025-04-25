import React, { useState } from "react";
import FilterDropdown from "./FilterDropdown";
import SelectedFilters from "./SelectedFilters";
import DateRangeFilter from "./DateRangeFilter";

// Define types for filter options
interface FilterOption {
  label: string;
  value: string;
}

// Define the structure of available filters
interface FiltersProps {
  states: FilterOption[];
  authors: FilterOption[];
  dockets: FilterOption[];
}

const SideFilters: React.FC<FiltersProps> = ({ states, authors, dockets }) => {
  // State to track selected filters
  const [selectedStates, setSelectedStates] = useState<string[]>([]);
  const [selectedAuthors, setSelectedAuthors] = useState<string[]>([]);
  const [selectedDockets, setSelectedDockets] = useState<string[]>([]);
  const [dateFrom, setDateFrom] = useState<string>("");
  const [dateTo, setDateTo] = useState<string>("");

  // Function to remove a filter
  const removeFilter = (
    type: "states" | "authors" | "dockets",
    value: string,
  ) => {
    switch (type) {
      case "states":
        setSelectedStates(selectedStates.filter((state) => state !== value));
        break;
      case "authors":
        setSelectedAuthors(
          selectedAuthors.filter((author) => author !== value),
        );
        break;
      case "dockets":
        setSelectedDockets(
          selectedDockets.filter((docket) => docket !== value),
        );
        break;
    }
  };

  // Function to add a filter
  const addFilter = (type: "states" | "authors" | "dockets", value: string) => {
    switch (type) {
      case "states":
        if (!selectedStates.includes(value)) {
          setSelectedStates([...selectedStates, value]);
        }
        break;
      case "authors":
        if (!selectedAuthors.includes(value)) {
          setSelectedAuthors([...selectedAuthors, value]);
        }
        break;
      case "dockets":
        if (!selectedDockets.includes(value)) {
          setSelectedDockets([...selectedDockets, value]);
        }
        break;
    }
  };

  return (
    <div className="side-filters space-y-4 p-4 bg-base-100 shadow-2xl rounded-lg">
      <DateRangeFilter
        dateFrom={dateFrom}
        dateTo={dateTo}
        onDateFromChange={setDateFrom}
        onDateToChange={setDateTo}
      />

      <FilterDropdown
        title="States"
        options={states}
        selectedOptions={selectedStates}
        onSelectOption={(value) => addFilter("states", value)}
      />

      <FilterDropdown
        title="Authors"
        options={authors}
        selectedOptions={selectedAuthors}
        onSelectOption={(value) => addFilter("authors", value)}
      />

      <FilterDropdown
        title="Dockets"
        options={dockets}
        selectedOptions={selectedDockets}
        onSelectOption={(value) => addFilter("dockets", value)}
      />

      <SelectedFilters
        states={selectedStates}
        authors={selectedAuthors}
        dockets={selectedDockets}
        onRemoveFilter={removeFilter}
      />
    </div>
  );
};

export default SideFilters;
