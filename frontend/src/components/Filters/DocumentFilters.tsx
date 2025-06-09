import React, { useMemo, useCallback, useEffect, useState } from "react";
import { Dispatch, SetStateAction } from "react";
import {
  FilterFieldDefinition,
  FilterInputType,
  FilterValues,
  FilterConfigurationManager,
  FilterEndpoints,
  ValidationResult,
  FilterOption,
  createFilterManager,
} from "@/lib/filters";
import clsx from "clsx";

import {
  DynamicMultiSelect
} from "@/components/Filters/FilterMultiSelect";
import {
  DynamicDocumentFiltersProps,
  DynamicDocumentFilters
} from "@/components/Filters/DynamicFilters";


// =============================================================================
// TYPES AND INTERFACES
// =============================================================================

/**
 * Base props for all document filter components
 */
interface BaseDocumentFiltersProps {
  /** Current filter values */
  queryOptions: FilterValues;
  /** Function to update filter values */
  setQueryOptions: Dispatch<SetStateAction<FilterValues>>;
  /** Array of filter field IDs to display */
  showFields: string[];
  /** Optional array of filter field IDs to disable */
  disabledFields?: string[];
  /** Filter configuration manager instance */
  configManager: FilterConfigurationManager;
}

/**
 * Props for the main DynamicDocumentFilters component
 */

/**
 * Props for layout-specific wrapper components
 */
interface LayoutDocumentFiltersProps extends Omit<BaseDocumentFiltersProps, 'configManager'> {
  /** Backend endpoints configuration */
  endpoints: FilterEndpoints;
  /** Whether to apply max-width-xs to inputs */
  maxWidthXs?: boolean;
  /** Custom input class names */
  inputClassNames?: DynamicDocumentFiltersProps['inputClassNames'];
  /** Custom render functions for specific filter fields */
  customRenderers?: DynamicDocumentFiltersProps['customRenderers'];
  /** Event handlers */
  onFilterChange?: DynamicDocumentFiltersProps['onFilterChange'];
  onValidationChange?: DynamicDocumentFiltersProps['onValidationChange'];
  /** Whether to show validation errors inline */
  showValidationErrors?: boolean;
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

/**
 * Layout configuration options
 */
export enum FilterLayout {
  List = "list",
  Grid = "grid",
  Flex = "flex",
  Custom = "custom",
}

// =============================================================================
// CUSTOM COMPONENTS FOR DYNAMIC INPUTS
// =============================================================================

/**
 * Dynamic multi-select component that loads options from configuration
 */


// =============================================================================
// EXAMPLE USAGE
// =============================================================================

/*
// Example field definition for testing
const exampleFieldDefinition = {
  id: "matter_type",
  displayName: "Matter Type",
  description: "Primary legal matter categories",
  inputType: FilterInputType.MultiSelect,
  placeholder: "Search and select matter types...",
  options: [
    { value: "litigation", label: "Litigation", disabled: false },
    { value: "corporate", label: "Corporate Law", disabled: false },
    { value: "criminal", label: "Criminal Law", disabled: false },
    { value: "family", label: "Family Law", disabled: false },
    { value: "real_estate", label: "Real Estate Law", disabled: false },
    { value: "employment", label: "Employment Law", disabled: false },
    { value: "intellectual_property", label: "Intellectual Property", disabled: false },
    { value: "tax", label: "Tax Law", disabled: false },
    { value: "bankruptcy", label: "Bankruptcy", disabled: false },
    { value: "immigration", label: "Immigration Law", disabled: false },
  ]
};

// Usage example
function ExampleUsage() {
  const [selectedValue, setSelectedValue] = useState("");

  return (
    <div className="p-6 max-w-md">
      <h3 className="text-lg font-semibold mb-4">Enhanced Multi-Select Example</h3>
      
      <DynamicMultiSelect
        fieldDefinition={exampleFieldDefinition}
        value={selectedValue}
        onChange={setSelectedValue}
        onFocus={() => console.log('Focused')}
        onBlur={() => console.log('Blurred')}
      />
      
      <div className="mt-4 p-3 bg-gray-100 rounded">
        <strong>Selected Value:</strong>
        <pre className="text-sm mt-1">{selectedValue || '(none)'}</pre>
      </div>
    </div>
  );
}
*/
/**
 * Dynamic date range picker component
 */
interface DynamicDateRangeProps {
  fieldDefinition: FilterFieldDefinition;
  value: string;
  onChange: (value: string) => void;
  onFocus?: () => void;
  onBlur?: () => void;
  disabled?: boolean;
  className?: string;
}

/**
 * Dynamic date range picker component
 */
function DynamicDateRange(props: DynamicDateRangeProps): React.ReactElement {
  const {
    fieldDefinition,
    value,
    onChange,
    onFocus,
    onBlur,
    disabled = false,
    className,
  } = props;
  const [startDate, endDate] = value.split('|');

  const handleStartDateChange = useCallback((newStartDate: string) => {
    const newValue = `${newStartDate}|${endDate || ''}`;
    onChange(newValue);
  }, [endDate, onChange]);

  const handleEndDateChange = useCallback((newEndDate: string) => {
    const newValue = `${startDate || ''}|${newEndDate}`;
    onChange(newValue);
  }, [startDate, onChange]);

  return (
    <div className={clsx("flex gap-2", className)}>
      <input
        type="date"
        value={startDate || ''}
        onChange={(e) => handleStartDateChange(e.target.value)}
        onFocus={onFocus}
        onBlur={onBlur}
        disabled={disabled}
        className="input input-bordered flex-1"
        placeholder="Start date"
      />
      <span className="self-center">to</span>
      <input
        type="date"
        value={endDate || ''}
        onChange={(e) => handleEndDateChange(e.target.value)}
        onFocus={onFocus}
        onBlur={onBlur}
        disabled={disabled}
        className="input input-bordered flex-1"
        placeholder="End date"
      />
    </div>
  );
}

// =============================================================================
// LAYOUT WRAPPER COMPONENTS
// =============================================================================

/**
 * Document filters displayed in a vertical list layout
 */
export function DynamicDocumentFiltersList(props: LayoutDocumentFiltersProps): React.ReactElement {
  return (
    <DynamicDocumentFiltersWithManager
      {...props}
      className="grid grid-flow-row auto-rows-max gap-4"
    />
  );
}

/**
 * Document filters displayed in a 4-column grid layout
 */
export function DynamicDocumentFiltersGrid(props: LayoutDocumentFiltersProps): React.ReactElement {
  return (
    <DynamicDocumentFiltersWithManager
      {...props}
      className="grid grid-cols-4 gap-4"
    />
  );
}

/**
 * Document filters displayed in a responsive grid layout
 */
export function ResponsiveDynamicDocumentFilters(props: LayoutDocumentFiltersProps): React.ReactElement {
  return (
    <DynamicDocumentFiltersWithManager
      {...props}
      className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
    />
  );
}

/**
 * Document filters displayed in a horizontal flex layout
 */
export function InlineDynamicDocumentFilters(props: LayoutDocumentFiltersProps): React.ReactElement {
  return (
    <DynamicDocumentFiltersWithManager
      {...props}
      className="flex flex-wrap gap-4"
      maxWidthXs={true}
    />
  );
}

// =============================================================================
// WRAPPER COMPONENT WITH CONFIGURATION MANAGER
// =============================================================================

/**
 * Wrapper component that manages the filter configuration loading
 */
function DynamicDocumentFiltersWithManager(props: LayoutDocumentFiltersProps & { className?: string }): React.ReactElement {
  const { endpoints, ...restProps } = props;
  const [configManager] = useState(() => createFilterManager(endpoints));
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadConfiguration = async () => {
      try {
        setIsLoading(true);
        setError(null);
        await configManager.loadConfiguration();
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load filter configuration');
      } finally {
        setIsLoading(false);
      }
    };

    loadConfiguration();
  }, [configManager]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-8">
        <span className="loading loading-spinner loading-md"></span>
        <span className="ml-2">Loading filter configuration...</span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-error">
        <span>Error loading filters: {error}</span>
      </div>
    );
  }

  return (
    <DynamicDocumentFilters
      {...restProps}
      configManager={configManager}
      isLoading={isLoading}
      error={error}
    />
  );
}


// =============================================================================
// UTILITY HOOKS
// =============================================================================

/**
 * Custom hook for managing filter state with validation and persistence
 * @param initialFilters - Initial filter values
 * @param onFiltersChange - Optional callback when filters change
 * @returns Filter state and management functions
 */
export function useDocumentFilters(
  initialFilters: FilterValues,
  onFiltersChange?: (filters: FilterValues) => void
) {
  const [filters, setFilters] = useState<FilterValues>(initialFilters);

  const updateFilters = useCallback((
    newFilters: FilterValues | ((prev: FilterValues) => FilterValues)
  ) => {
    setFilters((prev) => {
      const updated = typeof newFilters === 'function' ? newFilters(prev) : newFilters;
      onFiltersChange?.(updated);
      return updated;
    });
  }, [onFiltersChange]);

  const clearFilters = useCallback(() => {
    updateFilters(initialFilters);
  }, [initialFilters, updateFilters]);

  const hasActiveFilters = useMemo(() => {
    return Object.values(filters).some(value => value !== "");
  }, [filters]);

  const getActiveFilterCount = useMemo(() => {
    return Object.values(filters).filter(value => value !== "").length;
  }, [filters]);

  const setFieldValue = useCallback((fieldId: string, value: string) => {
    updateFilters(prev => ({ ...prev, [fieldId]: value }));
  }, [updateFilters]);

  const clearField = useCallback((fieldId: string) => {
    updateFilters(prev => ({ ...prev, [fieldId]: "" }));
  }, [updateFilters]);

  return {
    filters,
    setFilters: updateFilters,
    clearFilters,
    hasActiveFilters,
    getActiveFilterCount,
    setFieldValue,
    clearField,
  };
}

/**
 * Custom hook for managing filter configuration loading and caching
 * @param endpoints - Backend endpoints configuration
 * @returns Configuration manager and loading state
 */
export function useFilterConfiguration(endpoints: FilterEndpoints) {
  const [configManager, setConfigManager] = useState<FilterConfigurationManager | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadConfiguration = async () => {
      try {
        setIsLoading(true);
        setError(null);

        const manager = createFilterManager(endpoints);
        await manager.loadConfiguration();

        setConfigManager(manager);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Failed to load filter configuration';
        setError(errorMessage);
        console.error('Filter configuration loading error:', err);
      } finally {
        setIsLoading(false);
      }
    };

    loadConfiguration();
  }, [endpoints]);

  const reloadConfiguration = useCallback(async () => {
    if (!configManager) return;

    try {
      setIsLoading(true);
      setError(null);
      await configManager.loadConfiguration(true); // Force reload
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to reload filter configuration';
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  }, [configManager]);

  return {
    configManager,
    isLoading,
    error,
    reloadConfiguration,
  };
}

/**
 * Custom hook for filter validation with debouncing
 * @param filters - Current filter values
 * @param configManager - Filter configuration manager
 * @param debounceMs - Debounce delay in milliseconds
 * @returns Validation result and validation state
 */
export function useFilterValidation(
  filters: FilterValues,
  configManager: FilterConfigurationManager | null,
  debounceMs: number = 300
) {
  const [validationResult, setValidationResult] = useState<ValidationResult>({
    isValid: true,
    errors: [],
    warnings: []
  });
  const [isValidating, setIsValidating] = useState(false);

  // Debounced validation effect
  useEffect(() => {
    if (!configManager) return;

    const timeoutId = setTimeout(async () => {
      setIsValidating(true);
      try {
        const result = configManager.validateFilters(filters);
        setValidationResult(result);
      } catch (error) {
        console.error('Validation error:', error);
        setValidationResult({
          isValid: false,
          errors: [{ fieldId: '_system', message: 'Validation failed', type: 'system' }],
          warnings: []
        });
      } finally {
        setIsValidating(false);
      }
    }, debounceMs);

    return () => clearTimeout(timeoutId);
  }, [filters, configManager, debounceMs]);

  return {
    validationResult,
    isValidating,
  };
}

// Default export for backward compatibility
export default DynamicDocumentFilters;
