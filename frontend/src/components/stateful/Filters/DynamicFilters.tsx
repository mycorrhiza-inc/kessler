"use client"
import React, { useMemo, useCallback, useEffect, useState } from "react";
import {
  FilterFieldDefinition,
  FilterInputType,
  FilterValues,
  FilterConfigurationManager,
  FilterEndpoints,
  ValidationResult,
  createFilterManager,
} from "@/lib/filters";
import clsx from "clsx";

// Import Kessler store hooks
import {
  useFilters,
  useFilterField,
  useFilterFields,
  useFilterLoadingState,
  useFilterInitialization,
  useFilterSystemLifecycle,
  useUrlSync,
  useFilterPersistence,
  initializeFilterSystem,
  safeUpdateFilter,
  safeBulkUpdateFilters,
  getFilterSystemStatus,
} from "@/lib/store";

// Import the canonical multiselect component
import { DynamicMultiSelect } from "@/components/stateful/Filters/FilterMultiSelect";

// =============================================================================
// ENHANCED TYPES WITH KESSLER INTEGRATION
// =============================================================================

/**
 * Base props for Kessler-integrated filter components
 */
export interface KesslerFilterProps {
  /** Array of filter field IDs to display */
  showFields: string[];
  /** Optional array of filter field IDs to disable */
  disabledFields?: string[];
  /** Backend endpoints configuration */
  endpoints: FilterEndpoints;
  /** Whether to apply max-width-xs to inputs */
  maxWidthXs?: boolean;
  /** Custom input class names */
  inputClassNames?: {
    text?: string;
    select?: string;
    date?: string;
    number?: string;
  };
  /** Custom render functions for specific filter fields */
  customRenderers?: Record<string, FilterRenderer>;
  /** Event handlers */
  onFilterChange?: (fieldId: string, value: string) => void;
  onValidationChange?: (validation: ValidationResult) => void;
  /** Whether to show validation errors inline */
  showValidationErrors?: boolean;
  /** Auto-initialize the filter system */
  autoInitialize?: boolean;
  /** Enable URL synchronization */
  enableUrlSync?: boolean;
  /** Enable filter persistence */
  enablePersistence?: boolean;
}

/**
 * Props for the main Kessler-integrated component
 */
export interface KesslerDynamicFiltersProps extends KesslerFilterProps {
  /** CSS class name for the container */
  className?: string;
  /** Loading state override */
  isLoading?: boolean;
  /** Error state override */
  error?: string | null;
}

/**
 * Custom filter renderer function type
 */
type FilterRenderer = (
  fieldId: string,
  fieldDefinition: FilterFieldDefinition,
  currentValue: string,
  onChange: (value: string) => void,
  isDisabled: boolean
) => React.ReactElement;

// =============================================================================
// IMPROVED FILTER MANAGER INITIALIZATION
// =============================================================================

/**
 * Enhanced hook for managing filters with better error handling
 */
export const useKesslerFilters = (endpoints: FilterEndpoints, autoInitialize = true) => {
  const [filterManager, setFilterManager] = useState<FilterConfigurationManager | null>(null);
  const [initError, setInitError] = useState<string | null>(null);
  const [isCreatingManager, setIsCreatingManager] = useState(false);

  // Initialize filter manager with error handling
  useEffect(() => {
    let mounted = true;

    const createManager = async () => {
      if (isCreatingManager) return;

      try {
        setIsCreatingManager(true);
        setInitError(null);

        console.log('Creating filter manager with endpoints:', endpoints);

        // Validate endpoints
        if (!endpoints || typeof endpoints !== 'object') {
          throw new Error('Invalid endpoints configuration');
        }

        if (!endpoints.configuration) {
          throw new Error('Missing configuration endpoint');
        }

        // Create manager
        const manager = createFilterManager(endpoints);

        if (!mounted) return;

        console.log('Filter manager created successfully');
        setFilterManager(manager);

      } catch (error) {
        console.error('Failed to create filter manager:', error);
        if (mounted) {
          setInitError(error instanceof Error ? error.message : 'Failed to create filter manager');
        }
      } finally {
        if (mounted) {
          setIsCreatingManager(false);
        }
      }
    };

    createManager();

    return () => {
      mounted = false;
    };
  }, [endpoints, isCreatingManager]);

  // Use Kessler store lifecycle management with better error handling
  const lifecycle = useFilterSystemLifecycle(
    filterManager || undefined,
    [],
    autoInitialize && !!filterManager
  );

  // Get filter state and actions from Kessler store
  const filterState = useFilters();

  // Combine local and store errors
  const combinedError = initError || lifecycle.error;

  return {
    ...filterState,
    ...lifecycle,
    filterManager,
    error: combinedError,
    isCreatingManager,
    isReady: lifecycle.isReady && !combinedError && !isCreatingManager,
  };
};

/**
 * Hook for specific filter field with enhanced error handling
 */
export const useKesslerFilterField = (fieldId: string) => {
  const value = useFilterField(fieldId);
  const { setFilter, isFieldDisabled, filterManager } = useFilters();
  const isLoading = useFilterLoadingState(fieldId);

  const updateField = useCallback(async (newValue: string) => {
    try {
      const result = await safeUpdateFilter(fieldId, newValue);
      if (!result.success) {
        console.error(`Failed to update field ${fieldId}:`, result.error);
      }
      return result.success;
    } catch (error) {
      console.error(`Error updating field ${fieldId}:`, error);
      return false;
    }
  }, [fieldId]);

  const fieldDefinition = useMemo(() => {
    return filterManager?.getField(fieldId) || null;
  }, [filterManager, fieldId]);

  return {
    value,
    updateField,
    isLoading,
    isDisabled: isFieldDisabled(fieldId),
    fieldDefinition,
    setFilter: (value: string) => setFilter(fieldId, value),
  };
};

// =============================================================================
// ENHANCED BASE FILTER COMPONENT
// =============================================================================

/**
 * Base Kessler-integrated filter component with improved error handling
 */
function KesslerDynamicFilters(props: KesslerDynamicFiltersProps): React.ReactElement {
  const {
    className = "grid grid-flow-row auto-rows-max gap-4",
    showFields,
    disabledFields = [],
    endpoints,
    maxWidthXs = false,
    inputClassNames = {},
    customRenderers = {},
    onFilterChange,
    onValidationChange,
    showValidationErrors = true,
    autoInitialize = true,
    enableUrlSync = true,
    enablePersistence = true,
    isLoading: loadingOverride,
    error: errorOverride,
  } = props;

  // Initialize Kessler store with enhanced error handling
  const {
    filters,
    filterManager,
    filterError,
    validateFilters,
    setFilter,
    isReady,
    isInitializing,
    isCreatingManager,
    error: kesslerError,
  } = useKesslerFilters(endpoints, autoInitialize);

  // Enable URL sync if requested (only when ready)
  useUrlSync(enableUrlSync && isReady);

  // Enable persistence if requested (only when ready)
  useFilterPersistence(enablePersistence && isReady, enablePersistence && isReady);

  // Create disabled fields set
  const disabledFieldsSet = useMemo(() => new Set(disabledFields), [disabledFields]);

  // Get sorted field definitions
  const sortedFields = useMemo(() => {
    if (!filterManager) return [];

    try {
      return showFields
        .map((fieldId) => filterManager.getField(fieldId))
        .filter((field): field is FilterFieldDefinition => field !== null)
        .sort((a, b) => a.order - b.order);
    } catch (error) {
      console.error('Error getting field definitions:', error);
      return [];
    }
  }, [showFields, filterManager]);

  // Generate CSS class for input width constraint
  const maxWidthClass = useMemo(() => maxWidthXs ? "max-w-xs" : "", [maxWidthXs]);

  // Enhanced filter change handler with better error handling
  const handleFilterChange = useCallback(async (fieldId: string, value: string) => {
    try {
      const result = await safeUpdateFilter(fieldId, value);
      if (result.success) {
        onFilterChange?.(fieldId, value);
      } else {
        console.error(`Failed to update filter ${fieldId}:`, result.error);
      }
    } catch (error) {
      console.error(`Error updating filter ${fieldId}:`, error);
    }
  }, [onFilterChange]);

  // Validation effect with error handling
  useEffect(() => {
    if (onValidationChange && filterManager && isReady) {
      try {
        const validation = validateFilters();
        onValidationChange(validation);
      } catch (error) {
        console.error('Error validating filters:', error);
      }
    }
  }, [filters, validateFilters, onValidationChange, filterManager, isReady]);

  // =============================================================================
  // RENDER FUNCTIONS
  // =============================================================================

  const renderTextInput = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition,
    isDisabled: boolean
  ) => {
    const value = filters[fieldId] || "";

    return (
      <input
        className={clsx(
          "input input-bordered w-full",
          maxWidthClass,
          inputClassNames.text,
          isDisabled && "input-disabled"
        )}
        type="text"
        disabled={isDisabled}
        value={value}
        onChange={(e) => handleFilterChange(fieldId, e.target.value)}
        title={fieldDefinition.displayName}
        placeholder={fieldDefinition.placeholder}
        aria-label={fieldDefinition.displayName}
      />
    );
  }, [filters, handleFilterChange, maxWidthClass, inputClassNames.text]);

  const renderSelectInput = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition,
    isDisabled: boolean
  ) => {
    const value = filters[fieldId] || "";

    return (
      <select
        disabled={isDisabled}
        className={clsx(
          "select select-bordered w-full",
          maxWidthClass,
          inputClassNames.select,
          isDisabled && "select-disabled"
        )}
        value={value}
        onChange={(e) => handleFilterChange(fieldId, e.target.value)}
        aria-label={fieldDefinition.displayName}
      >
        <option value="">{fieldDefinition.placeholder || 'Select an option...'}</option>
        {fieldDefinition.options?.map((option) => (
          <option
            key={option.value}
            value={option.value}
            disabled={option.disabled}
          >
            {option.label}
          </option>
        ))}
      </select>
    );
  }, [filters, handleFilterChange, maxWidthClass, inputClassNames.select]);

  const renderDateInput = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition,
    isDisabled: boolean
  ) => {
    const value = filters[fieldId] || "";

    return (
      <input
        className={clsx(
          "input input-bordered w-full",
          maxWidthClass,
          inputClassNames.date,
          isDisabled && "input-disabled"
        )}
        type="date"
        disabled={isDisabled}
        value={value}
        onChange={(e) => handleFilterChange(fieldId, e.target.value)}
        title={fieldDefinition.displayName}
        aria-label={fieldDefinition.displayName}
      />
    );
  }, [filters, handleFilterChange, maxWidthClass, inputClassNames.date]);

  const renderMultiSelectInput = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition,
    isDisabled: boolean
  ) => {
    const value = filters[fieldId] || "";

    return (
      <DynamicMultiSelect
        fieldDefinition={fieldDefinition}
        value={value}
        onChange={(newValue) => handleFilterChange(fieldId, newValue)}
        disabled={isDisabled}
        className={clsx(maxWidthClass)}
      />
    );
  }, [filters, handleFilterChange, maxWidthClass]);

  const renderFilter = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition
  ): React.ReactElement | null => {
    const isDisabled = disabledFieldsSet.has(fieldId);

    // Check for custom renderer first
    if (customRenderers[fieldId]) {
      try {
        return customRenderers[fieldId](
          fieldId,
          fieldDefinition,
          filters[fieldId] || "",
          (value) => handleFilterChange(fieldId, value),
          isDisabled
        );
      } catch (error) {
        console.error(`Error rendering custom filter ${fieldId}:`, error);
        // Fall back to default renderer
      }
    }

    // Render based on input type
    switch (fieldDefinition.inputType) {
      case FilterInputType.Text:
        return renderTextInput(fieldId, fieldDefinition, isDisabled);
      case FilterInputType.Select:
        return renderSelectInput(fieldId, fieldDefinition, isDisabled);
      case FilterInputType.Date:
        return renderDateInput(fieldId, fieldDefinition, isDisabled);
      case FilterInputType.MultiSelect:
        return renderMultiSelectInput(fieldId, fieldDefinition, isDisabled);
      case FilterInputType.Hidden:
        return null;
      default:
        return renderTextInput(fieldId, fieldDefinition, isDisabled);
    }
  }, [
    disabledFieldsSet,
    customRenderers,
    filters,
    handleFilterChange,
    renderTextInput,
    renderSelectInput,
    renderDateInput,
    renderMultiSelectInput,
  ]);

  const renderFilterField = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition
  ): React.ReactElement => {
    const filterInput = renderFilter(fieldId, fieldDefinition);

    if (!filterInput) {
      return <React.Fragment key={fieldId} />;
    }

    return (
      <div key={fieldId} className="form-control w-full">
        <div className="label">
          <span className="label-text">
            <div className="tooltip tooltip-top" data-tip={fieldDefinition.description}>
              <span className="cursor-help">
                {fieldDefinition.displayName}
                {fieldDefinition.required && (
                  <span className="text-error ml-1">*</span>
                )}
              </span>
            </div>
          </span>
        </div>
        {filterInput}
      </div>
    );
  }, [renderFilter]);

  // =============================================================================
  // RENDER LOGIC WITH ENHANCED ERROR HANDLING
  // =============================================================================

  const isLoading = loadingOverride ?? isInitializing ?? isCreatingManager;
  const error = errorOverride ?? kesslerError ?? filterError;

  // Loading state
  if (isLoading) {
    return (
      <div className={clsx(className, "flex items-center justify-center p-8")}>
        <span className="loading loading-spinner loading-md mr-2"></span>
        <span>
          {isCreatingManager ? 'Creating filter manager...' :
            isInitializing ? 'Initializing filters...' :
              'Loading filters...'}
        </span>
      </div>
    );
  }

  // Error state with detailed information
  if (error) {
    return (
      <div className={className}>
        <div className="alert alert-error">
          <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <h3 className="font-bold">Filter System Error</h3>
            <div className="text-sm">{error}</div>
            <div className="mt-2">
              <button
                className="btn btn-sm btn-outline"
                onClick={() => window.location.reload()}
              >
                Retry
              </button>
            </div>
          </div>
        </div>

        {/* Debug Information */}
        <div className="mt-4 collapse collapse-arrow bg-base-200">
          <input type="checkbox" />
          <div className="collapse-title text-sm font-medium">
            🔍 Debug Information
          </div>
          <div className="collapse-content text-xs">
            <div className="space-y-2">
              <div><strong>Filter Manager:</strong> {filterManager ? 'Created' : 'Not created'}</div>
              <div><strong>Is Ready:</strong> {isReady ? 'Yes' : 'No'}</div>
              <div><strong>Endpoints:</strong> {JSON.stringify(endpoints, null, 2)}</div>
              <div><strong>Show Fields:</strong> {JSON.stringify(showFields)}</div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Not ready state
  if (!isReady) {
    return (
      <div className={className}>
        <div className="alert alert-warning">
          <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
          <span>Filter system is not ready. Please wait for initialization to complete...</span>
        </div>
      </div>
    );
  }

  // No fields to display
  if (sortedFields.length === 0) {
    return (
      <div className={className}>
        <div className="alert alert-info">
          <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <span>No filter fields configured for display.</span>
        </div>
      </div>
    );
  }

  // Render filters
  return (
    <div className={className} role="group" aria-label="Document filters">
      {sortedFields.map((fieldDefinition) =>
        renderFilterField(fieldDefinition.id, fieldDefinition)
      )}
    </div>
  );
}

// =============================================================================
// LAYOUT WRAPPER COMPONENTS WITH ERROR BOUNDARIES
// =============================================================================

/**
 * Error boundary wrapper for filter components
 */
function FilterErrorBoundary({
  children,
  fallback
}: {
  children: React.ReactNode;
  fallback?: React.ReactNode;
}) {
  const [hasError, setHasError] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const handleError = (event: ErrorEvent) => {
      setHasError(true);
      setError(new Error(event.message));
    };

    window.addEventListener('error', handleError);
    return () => window.removeEventListener('error', handleError);
  }, []);

  if (hasError) {
    return fallback || (
      <div className="alert alert-error">
        <span>Filter component error: {error?.message}</span>
        <button
          className="btn btn-sm btn-outline ml-2"
          onClick={() => {
            setHasError(false);
            setError(null);
          }}
        >
          Retry
        </button>
      </div>
    );
  }

  return <>{children}</>;
}

/**
 * List layout with Kessler store integration and error boundary
 */
export function KesslerDocumentFiltersList(props: KesslerFilterProps): React.ReactElement {
  return (
    <FilterErrorBoundary>
      <KesslerDynamicFilters
        {...props}
        className="grid grid-flow-row auto-rows-max gap-4"
      />
    </FilterErrorBoundary>
  );
}

/**
 * Grid layout with Kessler store integration and error boundary
 */
export function KesslerDocumentFiltersGrid(props: KesslerFilterProps): React.ReactElement {
  return (
    <FilterErrorBoundary>
      <KesslerDynamicFilters
        {...props}
        className="grid grid-cols-4 gap-4"
      />
    </FilterErrorBoundary>
  );
}

/**
 * Responsive layout with Kessler store integration and error boundary
 */
export function KesslerResponsiveDynamicDocumentFilters(props: KesslerFilterProps): React.ReactElement {
  return (
    <FilterErrorBoundary>
      <KesslerDynamicFilters
        {...props}
        className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
      />
    </FilterErrorBoundary>
  );
}

/**
 * Inline layout with Kessler store integration and error boundary
 */
export function KesslerInlineDynamicDocumentFilters(props: KesslerFilterProps): React.ReactElement {
  return (
    <FilterErrorBoundary>
      <KesslerDynamicFilters
        {...props}
        className="flex flex-wrap gap-4"
        maxWidthXs={true}
      />
    </FilterErrorBoundary>
  );
}

// =============================================================================
// ENHANCED FILTER CONTROLS WITH BETTER ERROR HANDLING
// =============================================================================

/**
 * Filter control panel with enhanced error handling
 */
export function KesslerFilterControls(): React.ReactElement {
  const {
    hasActiveFilters,
    resetFilters,
    persistFilters,
    clearPersistedFilters,
    loadPersistedFilters,
    filters,
    filterError,
    isInitialized,
  } = useFilters();

  const [actionError, setActionError] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState<string | null>(null);

  const handleAction = async (action: string, fn: () => Promise<void> | void) => {
    try {
      setActionLoading(action);
      setActionError(null);
      await fn();
    } catch (error) {
      setActionError(error instanceof Error ? error.message : `Failed to ${action}`);
    } finally {
      setActionLoading(null);
    }
  };

  const handleClearAll = () => handleAction('clear filters', async () => {
    await resetFilters();
  });

  const handleSaveFilters = () => handleAction('save filters', () => {
    persistFilters();
  });

  const handleLoadSavedFilters = () => handleAction('load filters', async () => {
    const savedFilters = loadPersistedFilters();
    if (savedFilters) {
      const result = await safeBulkUpdateFilters(savedFilters);
      if (!result.success) {
        throw new Error(result.error || 'Failed to load filters');
      }
    }
  });

  const activeFilterCount = useMemo(() => {
    return Object.values(filters).filter(value => value !== "").length;
  }, [filters]);

  if (!isInitialized) {
    return (
      <div className="flex items-center gap-2 p-2">
        <span className="loading loading-spinner loading-sm"></span>
        <span className="text-sm">Initializing filters...</span>
      </div>
    );
  }

  return (
    <div className="flex items-center gap-2 p-4 bg-base-100 border rounded-lg">
      <div className="flex-1">
        <div className="text-sm font-medium">
          {hasActiveFilters() ? (
            <span className="text-primary">
              {activeFilterCount} active filter{activeFilterCount !== 1 ? 's' : ''}
            </span>
          ) : (
            <span className="text-gray-500">No active filters</span>
          )}
        </div>
        {(filterError || actionError) && (
          <div className="text-xs text-error mt-1">{filterError || actionError}</div>
        )}
      </div>

      <div className="flex gap-2">
        <button
          onClick={handleSaveFilters}
          className={`btn btn-sm btn-outline ${actionLoading === 'save filters' ? 'loading' : ''}`}
          disabled={!hasActiveFilters() || !!actionLoading}
          title="Save current filters"
        >
          💾 Save
        </button>

        <button
          onClick={handleLoadSavedFilters}
          className={`btn btn-sm btn-outline ${actionLoading === 'load filters' ? 'loading' : ''}`}
          disabled={!!actionLoading}
          title="Load saved filters"
        >
          📁 Load
        </button>

        <button
          onClick={handleClearAll}
          className={`btn btn-sm btn-outline btn-error ${actionLoading === 'clear filters' ? 'loading' : ''}`}
          disabled={!hasActiveFilters() || !!actionLoading}
          title="Clear all filters"
        >
          🗑️ Clear All
        </button>
      </div>
    </div>
  );
}

// =============================================================================
// FILTER STATUS INDICATOR WITH ENHANCED DIAGNOSTICS
// =============================================================================

/**
 * Enhanced filter system status indicator
 */
export function KesslerFilterStatus(): React.ReactElement {
  const [status, setStatus] = useState(() => getFilterSystemStatus());
  const [showDetails, setShowDetails] = useState(false);

  useEffect(() => {
    const interval = setInterval(() => {
      setStatus(getFilterSystemStatus());
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  const getStatusColor = () => {
    if (status.hasError) return "error";
    if (!status.isReady) return "warning";
    return "success";
  };

  const getStatusIcon = () => {
    if (status.hasError) return "❌";
    if (!status.isReady) return "⏳";
    return "✅";
  };

  return (
    <div className={`alert alert-${getStatusColor()} p-2`}>
      <div className="flex items-center gap-2 text-sm w-full">
        <span>{getStatusIcon()}</span>
        <div className="flex-1">
          <div className="font-medium">
            Filter System: {status.isReady ? 'Ready' : status.isInitializing ? 'Initializing' : 'Error'}
          </div>
          <div className="text-xs opacity-75">
            {status.activeFilterCount}/{status.totalFields} active •
            {status.hasInheritedFilters && ' Inherited • '}
            {status.urlSyncEnabled && ' URL Sync • '}
            Updates: {status.updateCount}
          </div>
        </div>

        <button
          onClick={() => setShowDetails(!showDetails)}
          className="btn btn-xs btn-ghost"
        >
          {showDetails ? '▲' : '▼'}
        </button>
      </div>

      {showDetails && (
        <div className="mt-2 p-2 bg-base-100 rounded text-xs">
          <div className="grid grid-cols-2 gap-2">
            <div><strong>Initialized:</strong> {status.isInitialized ? 'Yes' : 'No'}</div>
            <div><strong>Initializing:</strong> {status.isInitializing ? 'Yes' : 'No'}</div>
            <div><strong>Has Error:</strong> {status.hasError ? 'Yes' : 'No'}</div>
            <div><strong>URL Sync:</strong> {status.urlSyncEnabled ? 'On' : 'Off'}</div>
            <div><strong>Pending:</strong> {status.pendingUpdateCount}</div>
            <div><strong>Disabled:</strong> {status.disabledFieldCount}</div>
          </div>
          {status.errorMessage && (
            <div className="mt-2 text-error">
              <strong>Error:</strong> {status.errorMessage}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

// =============================================================================
// COMPLETE FILTER SYSTEM COMPONENT WITH ERROR BOUNDARIES
// =============================================================================

/**
 * Complete filter system with enhanced error handling
 */
interface KesslerFilterSystemProps extends KesslerFilterProps {
  /** Layout type */
  layout?: 'list' | 'grid' | 'responsive' | 'inline';
  /** Show filter controls */
  showControls?: boolean;
  /** Show status indicator */
  showStatus?: boolean;
  /** Container class name */
  containerClassName?: string;
}

export function KesslerFilterSystem(props: KesslerFilterSystemProps): React.ReactElement {
  const {
    layout = 'list',
    showControls = true,
    showStatus = false,
    containerClassName = "",
    ...filterProps
  } = props;

  const getLayoutComponent = () => {
    switch (layout) {
      case 'grid':
        return KesslerDocumentFiltersGrid;
      case 'responsive':
        return KesslerResponsiveDynamicDocumentFilters;
      case 'inline':
        return KesslerInlineDynamicDocumentFilters;
      default:
        return KesslerDocumentFiltersList;
    }
  };

  const LayoutComponent = getLayoutComponent();

  return (
    <FilterErrorBoundary>
      <div className={clsx("space-y-4", containerClassName)}>
        {showStatus && <KesslerFilterStatus />}
        {showControls && <KesslerFilterControls />}
        <LayoutComponent {...filterProps} />
      </div>
    </FilterErrorBoundary>
  );
}

// =============================================================================
// EXPORT ALL COMPONENTS WITH PROPER EXPORTS
// =============================================================================

// Primary exports
export {
  KesslerDynamicFilters as default,
  FilterErrorBoundary,
};


// Legacy compatibility exports for existing code
export const DynamicDocumentFiltersList = KesslerDocumentFiltersList;
export const DynamicDocumentFiltersGrid = KesslerDocumentFiltersGrid;
export const ResponsiveDynamicDocumentFilters = KesslerResponsiveDynamicDocumentFilters;
export const InlineDynamicDocumentFilters = KesslerInlineDynamicDocumentFilters;
