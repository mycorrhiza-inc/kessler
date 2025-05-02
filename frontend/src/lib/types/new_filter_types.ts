export interface Filter<T = unknown> {
  key: string;
  value: T;
  label?: string;
}

export type FilterKey = string;

export type Filters = {
  // Private implementation detail - could be Map, Record, or custom structure
  _store: Record<FilterKey, Filter>;
  _version: number; // For change tracking
};

// Implementation-agnostic interface
export interface FiltersManager {
  getFilter: (key: FilterKey) => Filter | undefined;
  setFilter: <T>(key: FilterKey, value: T, label?: string) => void;
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
export const filtersManager: FiltersManager = {
  getFilter: (key) => {
    return globalFilters._store[key];
  },

  setFilter: (key, value, label) => {
    const newFilter: Filter = {
      key,
      value,
      label: label || key,
    };

    globalFilters._store[key] = newFilter;
    globalFilters._version++;
  },

  deleteFilter: (key) => {
    if (globalFilters._store[key]) {
      delete globalFilters._store[key];
      globalFilters._version++;
      return true;
    }
    return false;
  },

  getAllFilters: () => {
    return Object.values(globalFilters._store);
  },

  toArray: () => {
    return Object.values(globalFilters._store);
  },

  serialize: () => {
    return JSON.stringify(globalFilters._store);
  },
};

// Singleton instance (could make this configurable)
let globalFilters: Filters = createFilters();

// For testing/dev - reset filters
export const resetFilters = () => {
  globalFilters = createFilters();
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
