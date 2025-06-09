import { create } from "zustand";
import { subscribeWithSelector } from "zustand/middleware";
import { debounce } from "lodash";
import React from "react";
import { 
  FilterValues, 
  FilterConfiguration, 
  FilterConfigurationManager,
  InheritedFilterValues,
  createInitialFiltersFromInherited,
  getDisabledFieldsFromInherited,
  hasActiveFilters,
  FilterInputType,
  clearAllFilters,
  ValidationResult
} from "@/lib/filters";

// =============================================================================
// ENHANCED TYPES
// =============================================================================

interface FieldLoadingState {
  [fieldId: string]: boolean;
}

interface PendingUpdate {
  fieldId: string;
  value: string;
  timestamp: number;
}

interface FilterState {
  // Current filter values
  filters: FilterValues;
  
  // Filter configuration and manager
  filterConfiguration: FilterConfiguration | null;
  filterManager: FilterConfigurationManager | null;
  
  // Inherited filters (from URL params, parent components, etc.)
  inheritedFilters: InheritedFilterValues;
  
  // Loading and error states
  isLoadingFilters: boolean;
  filterError: string | null;
  
  // Disabled fields (typically from inherited filters)
  disabledFields: string[];
  
  // Field-level loading states
  fieldLoadingStates: FieldLoadingState;
  
  // Initialization state
  isInitialized: boolean;
  isInitializing: boolean;
  initializationError: string | null;
  
  // Pending updates queue
  pendingUpdates: PendingUpdate[];
  
  // URL synchronization
  urlSyncEnabled: boolean;
  lastUrlSync: number;
  
  // Performance tracking
  lastUpdateTimestamp: number;
  updateCount: number;
}

interface FilterActions {
  // Filter value management
  setFilter: (fieldId: string, value: string) => Promise<boolean>;
  setFilters: (filters: FilterValues, skipValidation?: boolean) => Promise<boolean>;
  updateFilter: (fieldId: string, value: string) => Promise<boolean>;
  
  // Reset operations
  resetFilters: () => Promise<void>;
  resetSpecificFilters: (fieldIds: string[]) => Promise<void>;
  resetFiltersForDataset: (dataset?: string) => Promise<void>;
  
  // Configuration management
  setFilterConfiguration: (config: FilterConfiguration) => void;
  setFilterManager: (manager: FilterConfigurationManager) => void;
  
  // Inherited filters management
  setInheritedFilters: (inherited: InheritedFilterValues) => void;
  applyInheritedFilters: () => Promise<void>;
  
  // Loading and error states
  setIsLoadingFilters: (loading: boolean) => void;
  setFilterError: (error: string | null) => void;
  setFieldLoading: (fieldId: string, loading: boolean) => void;
  
  // Initialization
  initialize: (manager: FilterConfigurationManager, inheritedFilters?: InheritedFilterValues) => Promise<boolean>;
  setInitialized: (initialized: boolean) => void;
  setInitializing: (initializing: boolean) => void;
  setInitializationError: (error: string | null) => void;
  
  // Pending updates
  queueUpdate: (fieldId: string, value: string) => void;
  processPendingUpdates: () => Promise<void>;
  clearPendingUpdates: () => void;
  
  // URL synchronization - FIXED
  enableUrlSync: () => void;
  disableUrlSync: () => void;
  syncWithUrl: () => void;
  updateUrlFromFilters: () => void;
  
  // Utility functions
  hasActiveFilters: () => boolean;
  getFilterValue: (fieldId: string) => string;
  isFieldDisabled: (fieldId: string) => boolean;
  validateFilters: (filters?: FilterValues) => ValidationResult;
  
  // URL parameter handling
  resetFiltersFromUrl: () => Promise<void>;
  getDatasetFromUrl: () => string | null;
  
  // Persistence
  persistFilters: () => void;
  loadPersistedFilters: () => FilterValues | null;
  clearPersistedFilters: () => void;
  
  // Cleanup
  cleanup: () => void;
}

interface KesslerState extends FilterState {
  isLoggedIn: boolean;
  experimentalFeaturesEnabled: boolean;
  defaultState: string;
}

interface KesslerActions extends FilterActions {
  setIsLoggedIn: (isLoggedIn: boolean) => void;
  setExperimentalFeaturesEnabled: (enableExperimentalFeatures: boolean) => void;
  setDefaultState: (defaultState: string) => void;
}

type KesslerStore = KesslerState & KesslerActions;

// =============================================================================
// UTILITY FUNCTIONS
// =============================================================================

const isSSR = () => typeof window === "undefined";

const safeJsonParse = <T>(json: string, fallback: T): T => {
  try {
    return JSON.parse(json);
  } catch {
    return fallback;
  }
};

const sanitizeString = (str: string): string => {
  return str.replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, "");
};

const createSafeURLSearchParams = (): URLSearchParams | null => {
  if (isSSR()) return null;
  
  try {
    return new URLSearchParams(window.location.search);
  } catch {
    return null;
  }
};

const safePersistence = {
  setItem: (key: string, value: string): boolean => {
    if (isSSR()) return false;
    
    try {
      if (window.localStorage) {
        localStorage.setItem(key, value);
        return true;
      }
    } catch (error) {
      console.warn('Failed to persist to localStorage:', error);
    }
    return false;
  },
  
  getItem: (key: string): string | null => {
    if (isSSR()) return null;
    
    try {
      if (window.localStorage) {
        return localStorage.getItem(key);
      }
    } catch (error) {
      console.warn('Failed to read from localStorage:', error);
    }
    return null;
  },
  
  removeItem: (key: string): boolean => {
    if (isSSR()) return false;
    
    try {
      if (window.localStorage) {
        localStorage.removeItem(key);
        return true;
      }
    } catch (error) {
      console.warn('Failed to remove from localStorage:', error);
    }
    return false;
  }
};

// =============================================================================
// DEBOUNCED OPERATIONS
// =============================================================================

const debouncedPersist = debounce((filters: FilterValues) => {
  safePersistence.setItem('kessler_filters', JSON.stringify(filters));
}, 1000);

const debouncedUrlUpdate = debounce((filters: FilterValues) => {
  if (isSSR()) return;
  
  try {
    const url = new URL(window.location.href);
    const params = url.searchParams;
    
    // Clear existing filter params
    const filterKeys = Array.from(params.keys()).filter(key => 
      key.startsWith('filter_') || ['dataset', 'case_number', 'filing_type'].includes(key)
    );
    filterKeys.forEach(key => params.delete(key));
    
    // Add current filters
    Object.entries(filters).forEach(([key, value]) => {
      if (value && value.trim() !== '') {
        params.set(`filter_${key}`, sanitizeString(value));
      }
    });
    
    // Update URL without triggering navigation
    window.history.replaceState({}, '', url.toString());
  } catch (error) {
    console.warn('Failed to update URL:', error);
  }
}, 500);

// =============================================================================
// ZUSTAND STORE WITH MIDDLEWARE
// =============================================================================

export const useKesslerStore = create<KesslerStore>()(
  subscribeWithSelector((set, get) => ({
    // === EXISTING STATE ===
    experimentalFeaturesEnabled: false,
    isLoggedIn: false,
    defaultState: "new-york",
    
    // === FILTER STATE ===
    filters: {},
    filterConfiguration: null,
    filterManager: null,
    inheritedFilters: [],
    isLoadingFilters: false,
    filterError: null,
    disabledFields: [],
    fieldLoadingStates: {},
    isInitialized: false,
    isInitializing: false,
    initializationError: null,
    pendingUpdates: [],
    urlSyncEnabled: true,
    lastUrlSync: 0,
    lastUpdateTimestamp: 0,
    updateCount: 0,
    
    // === EXISTING ACTIONS ===
    setIsLoggedIn: (isLoggedIn: boolean) => set({ isLoggedIn }),
    setExperimentalFeaturesEnabled: (experimentalFeaturesEnabled: boolean) =>
      set({ experimentalFeaturesEnabled }),
    setDefaultState: (defaultState: string) => set({ defaultState }),
    
    // === ENHANCED FILTER VALUE MANAGEMENT ===
    setFilter: async (fieldId: string, value: string): Promise<boolean> => {
      const state = get();
      
      // Guard: Check initialization
      if (!state.isInitialized) {
        console.warn('Store not initialized, queueing filter update');
        state.queueUpdate(fieldId, value);
        return false;
      }
      
      // Guard: Check if currently loading
      if (state.isLoadingFilters) {
        console.warn('Cannot set filter while loading configuration');
        return false;
      }
      
      // Guard: Check filter manager
      if (!state.filterManager) {
        console.error('FilterManager not initialized');
        return false;
      }
      
      // Sanitize input
      const sanitizedValue = sanitizeString(value);
      
      // Validate field exists and is enabled
      const field = state.filterManager.getField(fieldId);
      if (!field) {
        console.warn(`Unknown filter field: ${fieldId}`);
        return false;
      }
      
      if (!field.enabled) {
        console.warn(`Filter field is disabled: ${fieldId}`);
        return false;
      }
      
      // Validate value
      const testFilters = { ...state.filters, [fieldId]: sanitizedValue };
      const validation = state.filterManager.validateFilters(testFilters);
      
      if (!validation.isValid) {
        console.error('Filter validation failed:', validation.errors);
        set({ filterError: validation.errors[0]?.message || 'Validation failed' });
        return false;
      }
      
      // Update state
      set((currentState) => ({
        filters: {
          ...currentState.filters,
          [fieldId]: sanitizedValue,
        },
        filterError: null,
        lastUpdateTimestamp: Date.now(),
        updateCount: currentState.updateCount + 1,
      }));
      
      // Trigger side effects
      const newState = get();
      newState.persistFilters();
      
      if (newState.urlSyncEnabled) {
        debouncedUrlUpdate(newState.filters);
      }
      
      return true;
    },
    
    setFilters: async (filters: FilterValues, skipValidation = false): Promise<boolean> => {
      const state = get();
      
      if (!state.isInitialized && !skipValidation) {
        console.warn('Store not initialized, cannot set filters');
        return false;
      }
      
      // Sanitize all values
      const sanitizedFilters: FilterValues = {};
      Object.entries(filters).forEach(([key, value]) => {
        sanitizedFilters[key] = sanitizeString(value || '');
      });
      
      // Validate if not skipping
      if (!skipValidation && state.filterManager) {
        const validation = state.filterManager.validateFilters(sanitizedFilters);
        if (!validation.isValid) {
          console.error('Bulk filter validation failed:', validation.errors);
          set({ filterError: validation.errors[0]?.message || 'Validation failed' });
          return false;
        }
      }
      
      set({
        filters: { ...sanitizedFilters },
        filterError: null,
        lastUpdateTimestamp: Date.now(),
        updateCount: state.updateCount + 1,
      });
      
      // Trigger side effects
      const newState = get();
      newState.persistFilters();
      
      if (newState.urlSyncEnabled) {
        debouncedUrlUpdate(newState.filters);
      }
      
      return true;
    },
    
    updateFilter: async (fieldId: string, value: string): Promise<boolean> => {
      return get().setFilter(fieldId, value);
    },
    
    // === ENHANCED RESET OPERATIONS ===
    resetFilters: async (): Promise<void> => {
      const state = get();
      
      let emptyFilters: FilterValues = {};
      if (state.filterManager) {
        emptyFilters = state.filterManager.createEmptyFilters();
      }
      
      await state.setFilters(emptyFilters, true);
    },
    
    resetSpecificFilters: async (fieldIds: string[]): Promise<void> => {
      const state = get();
      const newFilters = { ...state.filters };
      
      fieldIds.forEach((fieldId) => {
        newFilters[fieldId] = "";
      });
      
      await state.setFilters(newFilters);
    },
    
    resetFiltersForDataset: async (dataset?: string): Promise<void> => {
      const state = get();
      const targetDataset = dataset || state.getDatasetFromUrl();
      
      if (targetDataset) {
        console.log(`Resetting filters for dataset: ${targetDataset}`);
      }
      
      await state.resetFilters();
    },
    
    // === CONFIGURATION MANAGEMENT ===
    setFilterConfiguration: (config: FilterConfiguration) => {
      set({ filterConfiguration: config });
    },
    
    setFilterManager: (manager: FilterConfigurationManager) => {
      set({ filterManager: manager });
    },
    
    // === INHERITED FILTERS MANAGEMENT ===
    setInheritedFilters: (inherited: InheritedFilterValues) => {
      const disabledFields = getDisabledFieldsFromInherited(inherited);
      set({ 
        inheritedFilters: inherited,
        disabledFields 
      });
    },
    
    applyInheritedFilters: async (): Promise<void> => {
      const state = get();
      
      if (state.filterManager && state.inheritedFilters.length > 0) {
        const initialFilters = createInitialFiltersFromInherited(
          state.inheritedFilters,
          state.filterManager
        );
        await state.setFilters(initialFilters, true);
      }
    },
    
    // === LOADING AND ERROR STATES ===
    setIsLoadingFilters: (loading: boolean) => {
      set({ isLoadingFilters: loading });
    },
    
    setFilterError: (error: string | null) => {
      set({ filterError: error });
    },
    
    setFieldLoading: (fieldId: string, loading: boolean) => {
      set(state => ({
        fieldLoadingStates: {
          ...state.fieldLoadingStates,
          [fieldId]: loading
        }
      }));
    },
    
    // === INITIALIZATION ===
    initialize: async (manager: FilterConfigurationManager, inheritedFilters: InheritedFilterValues = []): Promise<boolean> => {
      const state = get();
      
      // Prevent concurrent initialization
      if (state.isInitializing || state.isInitialized) {
        console.warn('Filter system already initialized or initializing');
        return state.isInitialized;
      }
      
      try {
        set({ 
          isInitializing: true, 
          initializationError: null,
          filterError: null 
        });
        
        // Load configuration
        const config = await manager.loadConfiguration();
        
        // Set up store
        set({
          filterManager: manager,
          filterConfiguration: config,
        });
        
        // Set inherited filters
        if (inheritedFilters.length > 0) {
          state.setInheritedFilters(inheritedFilters);
        }
        
        // Try to load persisted filters first
        const persistedFilters = state.loadPersistedFilters();
        if (persistedFilters) {
          await state.setFilters(persistedFilters, true);
        } else {
          // Apply inherited filters if no persisted data
          if (inheritedFilters.length > 0) {
            await state.applyInheritedFilters();
          } else {
            // Initialize with empty filters
            await state.resetFilters();
          }
        }
        
        // Process any pending updates
        await state.processPendingUpdates();
        
        // Enable URL sync
        state.enableUrlSync();
        
        set({ 
          isInitialized: true,
          isInitializing: false,
          initializationError: null
        });
        
        console.log('Filter system initialized successfully');
        return true;
        
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : "Unknown initialization error";
        set({ 
          initializationError: errorMessage,
          isInitializing: false,
          isInitialized: false
        });
        console.error('Failed to initialize filter system:', error);
        return false;
      }
    },
    
    setInitialized: (initialized: boolean) => {
      set({ isInitialized: initialized });
    },
    
    setInitializing: (initializing: boolean) => {
      set({ isInitializing: initializing });
    },
    
    setInitializationError: (error: string | null) => {
      set({ initializationError: error });
    },
    
    // === PENDING UPDATES ===
    queueUpdate: (fieldId: string, value: string) => {
      set(state => ({
        pendingUpdates: [
          ...state.pendingUpdates,
          { fieldId, value, timestamp: Date.now() }
        ]
      }));
    },
    
    processPendingUpdates: async (): Promise<void> => {
      const state = get();
      
      if (state.pendingUpdates.length === 0) return;
      
      console.log(`Processing ${state.pendingUpdates.length} pending updates`);
      
      // Process updates sequentially
      for (const update of state.pendingUpdates) {
        await state.setFilter(update.fieldId, update.value);
      }
      
      state.clearPendingUpdates();
    },
    
    clearPendingUpdates: () => {
      set({ pendingUpdates: [] });
    },
    
    // === URL SYNCHRONIZATION - FIXED ===
    enableUrlSync: () => {
      const state = get();
      
      if (isSSR() || state.urlSyncEnabled) return;
      
      set({ urlSyncEnabled: true });
      
      // Listen for browser navigation
      const handlePopState = () => {
        const currentState = get();
        currentState.syncWithUrl();
      };
      
      if (!isSSR()) {
        window.addEventListener('popstate', handlePopState);
      }
    },
    
    disableUrlSync: () => {
      set({ urlSyncEnabled: false });
    },
    
    syncWithUrl: () => {
      if (isSSR()) return;
      
      const state = get();
      const urlParams = createSafeURLSearchParams();
      
      if (!urlParams) return;
      
      const urlFilters: FilterValues = {};
      
      // Extract filter parameters from URL
      for (const [key, value] of urlParams.entries()) {
        if (key.startsWith('filter_')) {
          const fieldId = key.replace('filter_', '');
          urlFilters[fieldId] = sanitizeString(value);
        }
      }
      
      // Update filters if they differ from current state
      const hasChanges = Object.keys(urlFilters).some(
        key => urlFilters[key] !== state.filters[key]
      );
      
      if (hasChanges) {
        state.setFilters({ ...state.filters, ...urlFilters });
        set({ lastUrlSync: Date.now() });
      }
    },
    
    updateUrlFromFilters: () => {
      const state = get();
      if (state.urlSyncEnabled) {
        debouncedUrlUpdate(state.filters);
      }
    },
    
    // === UTILITY FUNCTIONS ===
    hasActiveFilters: () => {
      const { filters } = get();
      return hasActiveFilters(filters);
    },
    
    getFilterValue: (fieldId: string) => {
      const { filters } = get();
      return filters[fieldId] || "";
    },
    
    isFieldDisabled: (fieldId: string) => {
      const { disabledFields } = get();
      return disabledFields.includes(fieldId);
    },
    
    validateFilters: (filters?: FilterValues): ValidationResult => {
      const state = get();
      const filtersToValidate = filters || state.filters;
      
      if (!state.filterManager) {
        return {
          isValid: false,
          errors: [{ fieldId: '_system', message: 'Filter manager not initialized', type: 'system' }],
          warnings: []
        };
      }
      
      return state.filterManager.validateFilters(filtersToValidate);
    },
    
    // === URL PARAMETER HANDLING ===
    resetFiltersFromUrl: async (): Promise<void> => {
      const state = get();
      const dataset = state.getDatasetFromUrl();
      await state.resetFiltersForDataset(dataset || undefined);
    },
    
    getDatasetFromUrl: () => {
      if (isSSR()) return null;
      
      const urlParams = createSafeURLSearchParams();
      if (!urlParams) return null;
      
      const dataset = urlParams.get("dataset");
      return dataset ? sanitizeString(dataset) : null;
    },
    
    // === PERSISTENCE ===
    persistFilters: () => {
      const { filters } = get();
      debouncedPersist(filters);
    },
    
    loadPersistedFilters: (): FilterValues | null => {
      const stored = safePersistence.getItem('kessler_filters');
      if (!stored) return null;
      
      const parsed = safeJsonParse<FilterValues>(stored, {});
      
      // Sanitize persisted values
      const sanitized: FilterValues = {};
      Object.entries(parsed).forEach(([key, value]) => {
        if (typeof value === 'string') {
          sanitized[key] = sanitizeString(value);
        }
      });
      
      return sanitized;
    },
    
    clearPersistedFilters: () => {
      safePersistence.removeItem('kessler_filters');
    },
    
    // === CLEANUP ===
    cleanup: () => {
      const state = get();
      
      // Cancel debounced operations
      debouncedPersist.cancel();
      debouncedUrlUpdate.cancel();
      
      // Clear pending updates
      state.clearPendingUpdates();
      
      // Reset state
      set({
        isInitialized: false,
        isInitializing: false,
        filterManager: null,
        filterConfiguration: null,
        filters: {},
        inheritedFilters: [],
        disabledFields: [],
        fieldLoadingStates: {},
        pendingUpdates: [],
        filterError: null,
        initializationError: null,
        urlSyncEnabled: false,
      });
    },
  }))
);

// =============================================================================
// GRANULAR SELECTOR HOOKS
// =============================================================================

/**
 * Hook for accessing a specific filter field value
 * Prevents unnecessary re-renders when other fields change
 */
export const useFilterField = (fieldId: string) => {
  return useKesslerStore(state => state.filters[fieldId] || "");
};

/**
 * Hook for accessing multiple specific filter fields
 */
export const useFilterFields = (fieldIds: string[]) => {
  return useKesslerStore(state => {
    const result: FilterValues = {};
    fieldIds.forEach(id => {
      result[id] = state.filters[id] || "";
    });
    return result;
  });
};

/**
 * Hook for filter loading states
 */
export const useFilterLoadingState = (fieldId?: string) => {
  return useKesslerStore(state => {
    if (fieldId) {
      return state.fieldLoadingStates[fieldId] || false;
    }
    return state.isLoadingFilters;
  });
};

/**
 * Hook for initialization state
 */
export const useFilterInitialization = () => {
  return useKesslerStore(state => ({
    isInitialized: state.isInitialized,
    isInitializing: state.isInitializing,
    initializationError: state.initializationError,
    initialize: state.initialize,
  }));
};

// =============================================================================
// ENHANCED HOOK VARIANTS
// =============================================================================

/**
 * Hook for filter-specific operations
 */
export const useFilters = () => {
  const store = useKesslerStore();
  
  return {
    // State
    filters: store.filters,
    filterConfiguration: store.filterConfiguration,
    filterManager: store.filterManager,
    inheritedFilters: store.inheritedFilters,
    isLoadingFilters: store.isLoadingFilters,
    filterError: store.filterError,
    disabledFields: store.disabledFields,
    fieldLoadingStates: store.fieldLoadingStates,
    isInitialized: store.isInitialized,
    isInitializing: store.isInitializing,
    initializationError: store.initializationError,
    pendingUpdates: store.pendingUpdates,
    urlSyncEnabled: store.urlSyncEnabled,
    
    // Actions
    setFilter: store.setFilter,
    setFilters: store.setFilters,
    updateFilter: store.updateFilter,
    resetFilters: store.resetFilters,
    resetSpecificFilters: store.resetSpecificFilters,
    resetFiltersForDataset: store.resetFiltersForDataset,
    setFilterConfiguration: store.setFilterConfiguration,
    setFilterManager: store.setFilterManager,
    setInheritedFilters: store.setInheritedFilters,
    applyInheritedFilters: store.applyInheritedFilters,
    setIsLoadingFilters: store.setIsLoadingFilters,
    setFilterError: store.setFilterError,
    setFieldLoading: store.setFieldLoading,
    initialize: store.initialize,
    
    // URL Sync - FIXED
    enableUrlSync: store.enableUrlSync,
    disableUrlSync: store.disableUrlSync,
    syncWithUrl: store.syncWithUrl,
    updateUrlFromFilters: store.updateUrlFromFilters,
    
    // Utilities
    hasActiveFilters: store.hasActiveFilters,
    getFilterValue: store.getFilterValue,
    isFieldDisabled: store.isFieldDisabled,
    validateFilters: store.validateFilters,
    resetFiltersFromUrl: store.resetFiltersFromUrl,
    getDatasetFromUrl: store.getDatasetFromUrl,
    persistFilters: store.persistFilters,
    loadPersistedFilters: store.loadPersistedFilters,
    clearPersistedFilters: store.clearPersistedFilters,
    cleanup: store.cleanup,
  };
};

/**
 * Hook for non-filter operations (legacy compatibility)
 */
export const useKesslerCore = () => {
  const store = useKesslerStore();
  
  return {
    // State
    isLoggedIn: store.isLoggedIn,
    experimentalFeaturesEnabled: store.experimentalFeaturesEnabled,
    defaultState: store.defaultState,
    
    // Actions
    setIsLoggedIn: store.setIsLoggedIn,
    setExperimentalFeaturesEnabled: store.setExperimentalFeaturesEnabled,
    setDefaultState: store.setDefaultState,
  };
};

// =============================================================================
// ERROR BOUNDARY INTEGRATION
// =============================================================================

/**
 * Hook for error boundary integration
 */
export const useErrorBoundary = () => {
  const [error, setError] = React.useState<Error | null>(null);

  const resetError = React.useCallback(() => {
    setError(null);
  }, []);

  const throwError = React.useCallback((error: Error) => {
    setError(error);
  }, []);

  React.useEffect(() => {
    if (error) {
      throw error;
    }
  }, [error]);

  return { throwError, resetError };
};

// =============================================================================
// ENHANCED UTILITY FUNCTIONS
// =============================================================================

/**
 * Enhanced initialization function with comprehensive error handling
 */
export const initializeFilterSystem = async (
  filterManager: FilterConfigurationManager,
  inheritedFilters: InheritedFilterValues = []
): Promise<boolean> => {
  const store = useKesslerStore.getState();
  return await store.initialize(filterManager, inheritedFilters);
};

/**
 * Safe filter update with validation
 */
export const safeUpdateFilter = async (
  fieldId: string, 
  value: string
): Promise<{ success: boolean; error?: string }> => {
  const store = useKesslerStore.getState();
  
  try {
    const success = await store.setFilter(fieldId, value);
    return { success };
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : 'Unknown error';
    return { success: false, error: errorMessage };
  }
};

/**
 * Bulk filter update with rollback on failure
 */
export const safeBulkUpdateFilters = async (
  filters: FilterValues
): Promise<{ success: boolean; error?: string; rolledBack?: boolean }> => {
  const store = useKesslerStore.getState();
  const originalFilters = { ...store.filters };
  
  try {
    const success = await store.setFilters(filters);
    if (!success) {
      // Rollback to original state
      await store.setFilters(originalFilters, true);
      return { success: false, error: 'Validation failed', rolledBack: true };
    }
    return { success: true };
  } catch (error) {
    // Rollback to original state
    await store.setFilters(originalFilters, true);
    const errorMessage = error instanceof Error ? error.message : 'Unknown error';
    return { success: false, error: errorMessage, rolledBack: true };
  }
};

/**
 * Get comprehensive filter system status
 */
export const getFilterSystemStatus = () => {
  const store = useKesslerStore.getState();
  
  return {
    isReady: store.isInitialized && !store.isInitializing && !store.initializationError,
    isInitialized: store.isInitialized,
    isInitializing: store.isInitializing,
    hasError: !!store.filterError || !!store.initializationError,
    errorMessage: store.filterError || store.initializationError,
    activeFilterCount: Object.values(store.filters).filter(v => v !== "").length,
    totalFields: store.filterConfiguration?.fields.length || 0,
    hasInheritedFilters: store.inheritedFilters.length > 0,
    disabledFieldCount: store.disabledFields.length,
    pendingUpdateCount: store.pendingUpdates.length,
    currentDataset: store.getDatasetFromUrl(),
    hasActiveFilters: store.hasActiveFilters(),
    urlSyncEnabled: store.urlSyncEnabled,
    lastUpdate: store.lastUpdateTimestamp,
    updateCount: store.updateCount,
  };
};

/**
 * Reset filters based on current URL parameters with error handling
 */
export const handleUrlBasedFilterReset = async (): Promise<{ success: boolean; error?: string }> => {
  try {
    const store = useKesslerStore.getState();
    await store.resetFiltersFromUrl();
    return { success: true };
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : 'Unknown error';
    return { success: false, error: errorMessage };
  }
};

/**
 * Enhanced filter cleanup with comprehensive teardown
 */
export const cleanupFilterSystem = (): void => {
  const store = useKesslerStore.getState();
  store.cleanup();
  console.log('Filter system cleaned up successfully');
};

/**
 * Performance monitoring utility
 */
export const getFilterPerformanceMetrics = () => {
  const store = useKesslerStore.getState();
  
  return {
    updateCount: store.updateCount,
    lastUpdateTimestamp: store.lastUpdateTimestamp,
    timeSinceLastUpdate: Date.now() - store.lastUpdateTimestamp,
    pendingUpdateCount: store.pendingUpdates.length,
    memoryUsage: {
      filtersSize: JSON.stringify(store.filters).length,
      configurationSize: store.filterConfiguration ? JSON.stringify(store.filterConfiguration).length : 0,
      totalStateSize: Object.keys(store).length,
    },
  };
};

// =============================================================================
// REACT HOOKS FOR COMPONENT INTEGRATION
// =============================================================================

/**
 * Hook for managing filter system lifecycle in components
 */
export const useFilterSystemLifecycle = (
  filterManager?: FilterConfigurationManager,
  inheritedFilters?: InheritedFilterValues,
  autoInitialize = true
) => {
  const { isInitialized, isInitializing, initializationError, initialize } = useFilterInitialization();
  const [localError, setLocalError] = React.useState<string | null>(null);

  React.useEffect(() => {
    if (autoInitialize && filterManager && !isInitialized && !isInitializing) {
      initialize(filterManager, inheritedFilters)
        .then((success) => {
          if (!success) {
            setLocalError('Failed to initialize filter system');
          }
        })
        .catch((error) => {
          setLocalError(error instanceof Error ? error.message : 'Unknown initialization error');
        });
    }
  }, [filterManager, inheritedFilters, autoInitialize, isInitialized, isInitializing, initialize]);

  // Cleanup on unmount
  React.useEffect(() => {
    return () => {
      if (isInitialized) {
        const store = useKesslerStore.getState();
        store.cleanup();
      }
    };
  }, [isInitialized]);

  return {
    isReady: isInitialized && !initializationError && !localError,
    isInitializing,
    error: initializationError || localError,
    retry: () => {
      setLocalError(null);
      if (filterManager) {
        initialize(filterManager, inheritedFilters);
      }
    },
  };
};

/**
 * Hook for URL synchronization in React components - FIXED
 */
export const useUrlSync = (enabled = true) => {
  const {
    urlSyncEnabled,
    enableUrlSync,
    disableUrlSync,
    syncWithUrl,
    updateUrlFromFilters
  } = useFilters();

  React.useEffect(() => {
    if (enabled && !urlSyncEnabled) {
      enableUrlSync();
    } else if (!enabled && urlSyncEnabled) {
      disableUrlSync();
    }
  }, [enabled, urlSyncEnabled, enableUrlSync, disableUrlSync]);

  // Sync on mount
  React.useEffect(() => {
    if (enabled && urlSyncEnabled) {
      syncWithUrl();
    }
  }, [enabled, urlSyncEnabled, syncWithUrl]);

  return {
    isEnabled: urlSyncEnabled,
    enable: enableUrlSync,
    disable: disableUrlSync,
    sync: syncWithUrl,
    updateUrl: updateUrlFromFilters,
  };
};

/**
 * Hook for persisting filters across sessions
 */
export const useFilterPersistence = (autoLoad = true, autoPersist = true) => {
  const { 
    persistFilters, 
    loadPersistedFilters, 
    clearPersistedFilters,
    isInitialized,
    setFilters 
  } = useFilters();

  // Load persisted filters on mount
  React.useEffect(() => {
    if (autoLoad && isInitialized) {
      const persisted = loadPersistedFilters();
      if (persisted && Object.keys(persisted).length > 0) {
        setFilters(persisted, true).catch(console.error);
      }
    }
  }, [autoLoad, isInitialized, loadPersistedFilters, setFilters]);

  // Subscribe to filter changes for auto-persistence
  React.useEffect(() => {
    if (!autoPersist) return;

    const unsubscribe = useKesslerStore.subscribe(
      (state) => state.filters,
      (filters) => {
        if (Object.keys(filters).length > 0) {
          persistFilters();
        }
      }
    );

    return unsubscribe;
  }, [autoPersist, persistFilters]);

  return {
    persist: persistFilters,
    load: loadPersistedFilters,
    clear: clearPersistedFilters,
  };
};

// =============================================================================
// TESTING UTILITIES
// =============================================================================

/**
 * Mock store for testing
 */
export const createMockFilterStore = (initialState?: Partial<KesslerState>) => {
  const mockStore = create<KesslerStore>()((set, get) => ({
    ...useKesslerStore.getState(),
    ...initialState,
  }));

  return mockStore;
};

/**
 * Test utilities for filter validation
 */
export const testUtils = {
  /**
   * Create test filter configuration
   */
  createTestConfiguration: (): FilterConfiguration => ({
    fields: [
      {
        id: 'test_field',
        backendKey: 'test_field',
        displayName: 'Test Field',
        description: 'Test field for unit tests',
        inputType: FilterInputType.Text,
        required: false,
        order: 1,
        category: 'test',
        validation: {},
        defaultValue: '',
        enabled: true,
      }
    ],
    categories: [
      {
        id: 'test',
        name: 'Test Category',
        order: 1,
      }
    ],
    config: {
      version: '1.0.0',
      lastUpdated: new Date().toISOString(),
      defaultCategory: 'test',
    },
  }),

  /**
   * Simulate network delay for testing
   */
  delay: (ms: number) => new Promise(resolve => setTimeout(resolve, ms)),

  /**
   * Reset store to clean state for testing
   */
  resetStore: () => {
    const store = useKesslerStore.getState();
    store.cleanup();
  },

  /**
   * Validate store state consistency
   */
  validateStoreConsistency: () => {
    const store = useKesslerStore.getState();
    const issues: string[] = [];

    if (store.isInitialized && !store.filterManager) {
      issues.push('Store is initialized but filterManager is null');
    }

    if (store.isInitialized && !store.filterConfiguration) {
      issues.push('Store is initialized but filterConfiguration is null');
    }

    if (store.pendingUpdates.length > 0 && store.isInitialized) {
      issues.push('Pending updates exist after initialization');
    }

    if (store.isInitializing && store.isInitialized) {
      issues.push('Store cannot be both initializing and initialized');
    }

    return {
      isConsistent: issues.length === 0,
      issues,
    };
  },
};

// =============================================================================
// PERFORMANCE MONITORING
// =============================================================================

/**
 * Performance monitoring hook
 */
export const useFilterPerformance = () => {
  const [metrics, setMetrics] = React.useState(() => getFilterPerformanceMetrics());

  React.useEffect(() => {
    const interval = setInterval(() => {
      setMetrics(getFilterPerformanceMetrics());
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  return metrics;
};

/**
 * Memory usage monitor
 */
export const useMemoryMonitor = (threshold = 1024 * 1024) => { // 1MB threshold
  const [warning, setWarning] = React.useState(false);

  React.useEffect(() => {
    const checkMemory = () => {
      const metrics = getFilterPerformanceMetrics();
      const totalSize = metrics.memoryUsage.filtersSize + metrics.memoryUsage.configurationSize;
      
      if (totalSize > threshold) {
        setWarning(true);
        console.warn('Filter system memory usage high:', totalSize, 'bytes');
      } else {
        setWarning(false);
      }
    };

    const interval = setInterval(checkMemory, 5000);
    return () => clearInterval(interval);
  }, [threshold]);

  return { memoryWarning: warning };
};

export default useKesslerStore;