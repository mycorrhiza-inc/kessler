import { useCallback, useState } from "react";

export type FilterData = any;

export type FilterKey = string;

export interface Filter {
  key: FilterKey;
  value: FilterData;
  label?: string;
}
export type Filters = {
  // Private implementation detail - could be Map, Record, or custom structure
  _store: Record<FilterKey, Filter>;
  _version: number; // For change tracking
};

// Implementation-agnostic interface
export interface FiltersManager {
  getFilter: (key: FilterKey) => Filter | undefined;
  setFilter: (key: FilterKey, value: FilterData, label?: string) => void;
  deleteFilter: (key: FilterKey) => boolean;
  getAllFilters: () => Filter[];
  toArray: () => Filter[];
  serialize: () => string;
}

// Concrete implementation using Record as backing store
export const createFilters = (initialFilters: Filter[] = []): Filters => {
  const store = initialFilters.reduce(
    (acc, filter) => {
      acc[filter.key] = filter;
      return acc;
    },
    {} as Record<FilterKey, Filter>,
  );

  return {
    _store: store,
    _version: 0,
  };
};

// Operation implementations
export const filtersGetFilter = (filters: Filters, key: FilterKey) => {
  return filters._store[key];
};

export const filtersSetFilter = (
  filters: Filters,
  key: FilterKey,
  value: FilterData,
  label?: string,
) => {
  const newFilter: Filter = {
    key,
    value,
    label: label || key,
  };

  filters._store[key] = newFilter;
  filters._version++;
};

export const filtersDeleteFilter = (filters: Filters, key: FilterKey) => {
  if (filters._store[key]) {
    delete filters._store[key];
    filters._version++;
    return true;
  }
  return false;
};

export const filtersGetAllFilters = (filters: Filters) => {
  return Object.values(filters._store);
};

export const filtersToArray = (filters: Filters) => {
  return Object.values(filters._store);
};

export const filtersSerialize = (filters: Filters) => {
  return JSON.stringify(filters._store);
};

// Reconstruction from serialized data
export const deserializeFilters = (serialized: string): Filters => {
  try {
    const parsed = JSON.parse(serialized);
    return createFilters(Object.values(parsed));
  } catch (e) {
    console.error("Failed to deserialize filters:", e);
    return createFilters();
  }
};

export const useFilterState = (initialFilters: Filter[] | Filters = []) => {
  const [filters, setFilters] = useState<Filters>(() =>
    Array.isArray(initialFilters)
      ? createFilters(initialFilters)
      : initialFilters,
  );
  // Memoized filter operations
  const setFilter = useCallback(
    (key: FilterKey, value: FilterData, label?: string) => {
      setFilters((current: Filters) => {
        // this might be an unecessary clone, in here for an abundance of caution, feel free to remove anytime
        const nextFilters = structuredClone(current);
        filtersSetFilter(nextFilters, key, value, label);
        return nextFilters;
      });
    },
    [],
  );

  const deleteFilter = useCallback((key: FilterKey) => {
    setFilters((current: Filters) => {
      const nextFilters = structuredClone(current);
      filtersDeleteFilter(nextFilters, key);
      return nextFilters;
    });
  }, []);

  const clearFilters = useCallback(() => {
    setFilters(createFilters([]));
  }, []);

  return {
    filters,
    setFilter,
    deleteFilter,
    clearFilters,
  };
};
