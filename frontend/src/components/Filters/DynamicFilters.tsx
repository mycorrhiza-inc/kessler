"use client"
import React, { useMemo, useCallback, useEffect, useState, useRef } from "react";
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

// =============================================================================
// TYPES AND INTERFACES
// =============================================================================

/**
 * Base props for all document filter components
 */
export interface BaseDocumentFiltersProps {
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
export interface DynamicDocumentFiltersProps extends BaseDocumentFiltersProps {
  /** CSS class name for the container */
  className?: string;
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
  onFilterFocus?: (fieldId: string) => void;
  onFilterBlur?: (fieldId: string) => void;
  onValidationChange?: (validation: ValidationResult) => void;
  /** Loading and error states */
  isLoading?: boolean;
  error?: string | null;
  /** Whether to show validation errors inline */
  showValidationErrors?: boolean;
}

/**
 * Props for layout-specific wrapper components
 */
export interface LayoutDocumentFiltersProps extends Omit<BaseDocumentFiltersProps, 'configManager'> {
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
interface DynamicMultiSelectProps {
  fieldDefinition: FilterFieldDefinition;
  value: string;
  onChange: (value: string) => void;
  onFocus?: () => void;
  onBlur?: () => void;
  disabled?: boolean;
  className?: string;
}

function DynamicMultiSelect(props: DynamicMultiSelectProps): React.ReactElement {
  const {
    fieldDefinition,
    value,
    onChange,
    onFocus,
    onBlur,
    disabled = false,
    className,
  } = props;

  const [selectedValues, setSelectedValues] = useState<string[]>(
    value ? value.split(',').filter(Boolean) : []
  );
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  const handleSelectionChange = useCallback((selectedOptions: string[]) => {
    setSelectedValues(selectedOptions);
    onChange(selectedOptions.join(','));
  }, [onChange]);

  const toggleOption = useCallback((optionValue: string) => {
    const newSelection = selectedValues.includes(optionValue)
      ? selectedValues.filter(v => v !== optionValue)
      : [...selectedValues, optionValue];
    handleSelectionChange(newSelection);
  }, [selectedValues, handleSelectionChange]);

  const removeSelectedItem = useCallback((optionValue: string, event: React.MouseEvent) => {
    event.stopPropagation();
    const newSelection = selectedValues.filter(v => v !== optionValue);
    handleSelectionChange(newSelection);
  }, [selectedValues, handleSelectionChange]);

  // Filter options based on search term
  const filteredOptions = useMemo(() => {
    if (!searchTerm) return fieldDefinition.options || [];
    return (fieldDefinition.options || []).filter(option =>
      option.label.toLowerCase().includes(searchTerm.toLowerCase()) ||
      option.value.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [fieldDefinition.options, searchTerm]);

  // Get selected option labels for display
  const selectedLabels = useMemo(() => {
    return selectedValues.map(value => {
      const option = fieldDefinition.options?.find(opt => opt.value === value);
      return option ? option.label : value;
    });
  }, [selectedValues, fieldDefinition.options]);

  const handleDropdownToggle = useCallback(() => {
    if (!disabled) {
      setIsOpen(!isOpen);
      if (!isOpen) {
        onFocus?.();
      } else {
        onBlur?.();
      }
    }
  }, [disabled, isOpen, onFocus, onBlur]);

  const handleSearchChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  }, []);

  const clearSearch = useCallback(() => {
    setSearchTerm('');
  }, []);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as Element;
      if (isOpen && !target.closest('.dropdown')) {
        setIsOpen(false);
        onBlur?.();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isOpen, onBlur]);

  return (
    <div className={clsx("dropdown", isOpen && "dropdown-open", className)}>
      {/* Main Button */}
      <div
        tabIndex={0}
        role="button"
        className={clsx(
          "btn btn-outline w-full justify-start min-h-fit h-auto py-2 px-3",
          disabled && "btn-disabled"
        )}
        onClick={handleDropdownToggle}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            handleDropdownToggle();
          }
        }}
      >
        <div className="w-full text-left">
          {selectedValues.length === 0 ? (
            <span className="text-gray-500">{fieldDefinition.placeholder || 'Select options...'}</span>
          ) : (
            <div className="space-y-1">
              <div className="font-medium text-sm text-gray-700">
                {fieldDefinition.displayName}: {selectedValues.length} selected
              </div>
              <div className="flex flex-wrap gap-1">
                {selectedLabels.map((label, index) => (
                  <span
                    key={selectedValues[index]}
                    className="inline-flex items-center gap-1 px-2 py-1 bg-primary/10 text-primary rounded-full text-xs"
                  >
                    {label}
                    <button
                      type="button"
                      onClick={(e) => removeSelectedItem(selectedValues[index], e)}
                      className="hover:bg-primary/20 rounded-full p-0.5 transition-colors"
                      aria-label={`Remove ${label}`}
                    >
                      <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  </span>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Dropdown Arrow */}
        <svg
          className={clsx("w-4 h-4 transition-transform flex-shrink-0", isOpen && "rotate-180")}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </div>

      {/* Dropdown Content */}
      {isOpen && (
        <div className="dropdown-content z-[1] menu p-0 shadow-lg bg-base-100 rounded-box w-full border border-base-300">
          {/* Search Box */}
          <div className="p-3 border-b border-base-300">
            <div className="relative">
              <input
                type="text"
                className="input input-bordered input-sm w-full pl-8 pr-8"
                placeholder={`Search ${fieldDefinition.displayName.toLowerCase()}...`}
                value={searchTerm}
                onChange={handleSearchChange}
                onClick={(e) => e.stopPropagation()}
              />
              <svg
                className="w-4 h-4 absolute left-2.5 top-1/2 transform -translate-y-1/2 text-gray-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              {searchTerm && (
                <button
                  type="button"
                  onClick={clearSearch}
                  className="absolute right-2.5 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              )}
            </div>
          </div>

          {/* Options List */}
          <div className="max-h-60 overflow-y-auto">
            {filteredOptions.length === 0 ? (
              <div className="p-3 text-center text-gray-500 text-sm">
                {searchTerm ? 'No options match your search' : 'No options available'}
              </div>
            ) : (
              <ul className="p-1">
                {filteredOptions.map((option) => {
                  const isSelected = selectedValues.includes(option.value);
                  const optionLabel = fieldDefinition.options?.find(opt => opt.value === option.value)?.label || option.value;

                  return (
                    <li key={option.value}>
                      <label
                        className={clsx(
                          "cursor-pointer flex items-center justify-between p-2 rounded hover:bg-base-200 transition-colors",
                          option.disabled && "opacity-50 cursor-not-allowed"
                        )}
                      >
                        <div className="flex items-center gap-3 flex-1 min-w-0">
                          <input
                            type="checkbox"
                            checked={isSelected}
                            onChange={() => !option.disabled && toggleOption(option.value)}
                            disabled={disabled || option.disabled}
                            className="checkbox checkbox-sm flex-shrink-0"
                            onClick={(e) => e.stopPropagation()}
                          />
                          <span className="text-sm truncate flex-1">{option.label}</span>
                        </div>

                        {/* Option Pill */}
                        {isSelected && (
                          <span className="inline-flex items-center gap-1 px-2 py-1 bg-primary/10 text-primary rounded-full text-xs flex-shrink-0 ml-2">
                            {optionLabel}
                            <button
                              type="button"
                              onClick={(e) => {
                                e.stopPropagation();
                                toggleOption(option.value);
                              }}
                              className="hover:bg-primary/20 rounded-full p-0.5 transition-colors"
                              aria-label={`Remove ${optionLabel}`}
                            >
                              <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                              </svg>
                            </button>
                          </span>
                        )}
                      </label>
                    </li>
                  );
                })}
              </ul>
            )}
          </div>

          {/* Footer with selection summary */}
          {selectedValues.length > 0 && (
            <div className="p-3 border-t border-base-300 bg-base-50 text-xs text-gray-600">
              {selectedValues.length} item{selectedValues.length !== 1 ? 's' : ''} selected
              {selectedValues.length > 0 && (
                <button
                  type="button"
                  onClick={() => handleSelectionChange([])}
                  className="ml-2 text-primary hover:text-primary-focus underline"
                >
                  Clear all
                </button>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

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
// MAIN COMPONENT
// =============================================================================

/**
 * Main dynamic document filters component
 * Renders filters based on backend configuration with extensible rendering
 */
export function DynamicDocumentFilters(props: DynamicDocumentFiltersProps): React.ReactElement {
  const {
    className = "grid grid-flow-row auto-rows-max gap-4",
    queryOptions,
    setQueryOptions,
    showFields,
    disabledFields = [],
    configManager,
    maxWidthXs = false,
    inputClassNames = {},
    customRenderers = {},
    onFilterChange,
    onFilterFocus,
    onFilterBlur,
    onValidationChange,
    isLoading = false,
    error,
    showValidationErrors = true,
  } = props;

  // =============================================================================
  // STATE AND MEMOIZED VALUES
  // =============================================================================

  const [validationResult, setValidationResult] = useState<ValidationResult>({
    isValid: true,
    errors: [],
    warnings: []
  });

  // Use ref to track if we've called onValidationChange to prevent loops
  const validationChangeRef = useRef<ValidationResult | null>(null);

  /**
   * Create a lookup set for disabled fields for O(1) access
   */
  const disabledFieldsSet = useMemo(() => {
    return new Set(disabledFields);
  }, [disabledFields]);

  /**
   * Get field definitions for the fields to show, sorted by order
   */
  const sortedFields = useMemo(() => {
    const fieldDefinitions = showFields
      .map((fieldId) => configManager.getField(fieldId))
      .filter((field): field is FilterFieldDefinition => field !== null)
      .sort((a, b) => a.order - b.order);

    return fieldDefinitions;
  }, [showFields, configManager]);

  /**
   * Generate CSS class for input width constraint
   */
  const maxWidthClass = useMemo(() => {
    return maxWidthXs ? "max-w-xs" : "";
  }, [maxWidthXs]);

  // =============================================================================
  // VALIDATION - FIXED TO PREVENT INFINITE LOOPS
  // =============================================================================

  /**
   * Validate current filter values with proper dependency management
   */
  useEffect(() => {
    const validation = configManager.validateFilters(queryOptions);

    // Only update if validation actually changed
    const hasChanged =
      validationResult.isValid !== validation.isValid ||
      validationResult.errors.length !== validation.errors.length ||
      validationResult.warnings.length !== validation.warnings.length;

    if (hasChanged) {
      setValidationResult(validation);

      // Only call onValidationChange if it's different from the last call
      if (
        onValidationChange &&
        (!validationChangeRef.current ||
          validationChangeRef.current.isValid !== validation.isValid ||
          validationChangeRef.current.errors.length !== validation.errors.length)
      ) {
        validationChangeRef.current = validation;
        onValidationChange(validation);
      }
    }
  }, [queryOptions, configManager]); // Removed onValidationChange from dependencies

  // =============================================================================
  // EVENT HANDLERS
  // =============================================================================

  /**
   * Handle input changes with validation
   */
  const handleInputChange = useCallback((
    event: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
    fieldId: string,
  ) => {
    const newValue = event.target.value;
    handleValueChange(fieldId, newValue);
  }, []);

  /**
   * Handle programmatic value changes
   */
  const handleValueChange = useCallback((fieldId: string, value: string) => {
    setQueryOptions((prevOptions) => ({
      ...prevOptions,
      [fieldId]: value,
    }));

    onFilterChange?.(fieldId, value);
  }, [setQueryOptions, onFilterChange]);

  /**
   * Handle focus events
   */
  const handleFocus = useCallback((fieldId: string) => {
    onFilterFocus?.(fieldId);
  }, [onFilterFocus]);

  /**
   * Handle blur events
   */
  const handleBlur = useCallback((fieldId: string) => {
    onFilterBlur?.(fieldId);
  }, [onFilterBlur]);

  // =============================================================================
  // RENDER FUNCTIONS
  // =============================================================================

  /**
   * Renders a text input filter
   */
  const renderTextInput = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition,
    isDisabled: boolean
  ) => (
    <input
      className={clsx(
        "input input-bordered w-full",
        maxWidthClass,
        inputClassNames.text,
        isDisabled && "input-disabled"
      )}
      type="text"
      disabled={isDisabled}
      value={queryOptions[fieldId] || ""}
      onChange={(e) => handleInputChange(e, fieldId)}
      onFocus={() => handleFocus(fieldId)}
      onBlur={() => handleBlur(fieldId)}
      title={fieldDefinition.displayName}
      placeholder={fieldDefinition.placeholder}
      aria-label={fieldDefinition.displayName}
      aria-describedby={`${fieldId}-description`}
    />
  ), [queryOptions, handleInputChange, handleFocus, handleBlur, maxWidthClass, inputClassNames.text]);

  /**
   * Renders a number input filter
   */
  const renderNumberInput = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition,
    isDisabled: boolean
  ) => (
    <input
      className={clsx(
        "input input-bordered w-full",
        maxWidthClass,
        inputClassNames.number,
        isDisabled && "input-disabled"
      )}
      type="number"
      disabled={isDisabled}
      value={queryOptions[fieldId] || ""}
      onChange={(e) => handleInputChange(e, fieldId)}
      onFocus={() => handleFocus(fieldId)}
      onBlur={() => handleBlur(fieldId)}
      title={fieldDefinition.displayName}
      placeholder={fieldDefinition.placeholder}
      aria-label={fieldDefinition.displayName}
      min={fieldDefinition.validation?.min}
      max={fieldDefinition.validation?.max}
    />
  ), [queryOptions, handleInputChange, handleFocus, handleBlur, maxWidthClass, inputClassNames.number]);

  /**
   * Renders a select dropdown filter
   */
  const renderSelectInput = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition,
    isDisabled: boolean
  ) => (
    <select
      disabled={isDisabled}
      className={clsx(
        "select select-bordered w-full",
        maxWidthClass,
        inputClassNames.select,
        isDisabled && "select-disabled"
      )}
      value={queryOptions[fieldId] || ""}
      onChange={(e) => handleInputChange(e, fieldId)}
      onFocus={() => handleFocus(fieldId)}
      onBlur={() => handleBlur(fieldId)}
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
  ), [queryOptions, handleInputChange, handleFocus, handleBlur, maxWidthClass, inputClassNames.select]);

  /**
   * Renders a date input filter
   */
  const renderDateInput = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition,
    isDisabled: boolean
  ) => (
    <input
      className={clsx(
        "input input-bordered w-full",
        maxWidthClass,
        inputClassNames.date,
        isDisabled && "input-disabled"
      )}
      type="date"
      disabled={isDisabled}
      value={queryOptions[fieldId] || ""}
      onChange={(e) => handleInputChange(e, fieldId)}
      onFocus={() => handleFocus(fieldId)}
      onBlur={() => handleBlur(fieldId)}
      title={fieldDefinition.displayName}
      aria-label={fieldDefinition.displayName}
    />
  ), [queryOptions, handleInputChange, handleFocus, handleBlur, maxWidthClass, inputClassNames.date]);

  /**
   * Renders a multi-select filter
   */
  const renderMultiSelectInput = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition,
    isDisabled: boolean
  ) => (
    <DynamicMultiSelect
      fieldDefinition={fieldDefinition}
      value={queryOptions[fieldId] || ""}
      onChange={(value) => handleValueChange(fieldId, value)}
      onFocus={() => handleFocus(fieldId)}
      onBlur={() => handleBlur(fieldId)}
      disabled={isDisabled}
      className={clsx(maxWidthClass)}
    />
  ), [queryOptions, handleValueChange, handleFocus, handleBlur, maxWidthClass]);

  /**
   * Renders a date range filter
   */
  const renderDateRangeInput = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition,
    isDisabled: boolean
  ) => (
    <DynamicDateRange
      fieldDefinition={fieldDefinition}
      value={queryOptions[fieldId] || ""}
      onChange={(value) => handleValueChange(fieldId, value)}
      onFocus={() => handleFocus(fieldId)}
      onBlur={() => handleBlur(fieldId)}
      disabled={isDisabled}
      className={clsx(maxWidthClass)}
    />
  ), [queryOptions, handleValueChange, handleFocus, handleBlur, maxWidthClass]);

  /**
   * Renders a boolean/checkbox filter
   */
  const renderBooleanInput = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition,
    isDisabled: boolean
  ) => (
    <div className="form-control">
      <label className="label cursor-pointer justify-start gap-2">
        <input
          type="checkbox"
          className="checkbox"
          disabled={isDisabled}
          checked={queryOptions[fieldId] === "true"}
          onChange={(e) => handleValueChange(fieldId, e.target.checked ? "true" : "false")}
          onFocus={() => handleFocus(fieldId)}
          onBlur={() => handleBlur(fieldId)}
          aria-label={fieldDefinition.displayName}
        />
        <span className="label-text">{fieldDefinition.placeholder || 'Enable'}</span>
      </label>
    </div>
  ), [queryOptions, handleValueChange, handleFocus, handleBlur]);

  /**
   * Renders a UUID input filter with validation
   */
  const renderUUIDInput = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition,
    isDisabled: boolean
  ) => (
    <input
      className={clsx(
        "input input-bordered w-full font-mono text-sm",
        maxWidthClass,
        inputClassNames.text,
        isDisabled && "input-disabled"
      )}
      type="text"
      disabled={isDisabled}
      value={queryOptions[fieldId] || ""}
      onChange={(e) => handleInputChange(e, fieldId)}
      onFocus={() => handleFocus(fieldId)}
      onBlur={() => handleBlur(fieldId)}
      title={fieldDefinition.displayName}
      placeholder={fieldDefinition.placeholder || "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}
      aria-label={fieldDefinition.displayName}
      pattern="[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
    />
  ), [queryOptions, handleInputChange, handleFocus, handleBlur, maxWidthClass, inputClassNames.text]);

  /**
   * Main filter rendering function with support for custom renderers
   */
  const renderFilter = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition
  ): React.ReactElement | null => {
    const isDisabled = disabledFieldsSet.has(fieldId);

    // Check for custom renderer first
    if (customRenderers[fieldId]) {
      return customRenderers[fieldId](
        fieldId,
        fieldDefinition,
        queryOptions[fieldId] || "",
        (value) => handleValueChange(fieldId, value),
        isDisabled
      );
    }

    // Handle disabled state with generic text input for most types
    if (isDisabled && fieldDefinition.inputType !== FilterInputType.Hidden) {
      return renderTextInput(fieldId, fieldDefinition, true);
    }

    // Render based on input type
    switch (fieldDefinition.inputType) {
      case FilterInputType.Text:
        return renderTextInput(fieldId, fieldDefinition, false);

      case FilterInputType.Number:
        return renderNumberInput(fieldId, fieldDefinition, false);

      case FilterInputType.Select:
        return renderSelectInput(fieldId, fieldDefinition, false);

      case FilterInputType.MultiSelect:
        return renderMultiSelectInput(fieldId, fieldDefinition, false);

      case FilterInputType.Date:
        return renderDateInput(fieldId, fieldDefinition, false);

      case FilterInputType.DateRange:
        return renderDateRangeInput(fieldId, fieldDefinition, false);

      case FilterInputType.Boolean:
        return renderBooleanInput(fieldId, fieldDefinition, false);

      case FilterInputType.UUID:
        return renderUUIDInput(fieldId, fieldDefinition, false);

      case FilterInputType.Custom:
        console.warn(`Custom input type for field ${fieldId} requires a custom renderer`);
        return renderTextInput(fieldId, fieldDefinition, false);

      case FilterInputType.Hidden:
        return null;

      default:
        console.warn(`Unknown input type for filter ${fieldId}:`, fieldDefinition.inputType);
        return renderTextInput(fieldId, fieldDefinition, false);
    }
  }, [
    disabledFieldsSet,
    customRenderers,
    queryOptions,
    handleValueChange,
    renderTextInput,
    renderNumberInput,
    renderSelectInput,
    renderMultiSelectInput,
    renderDateInput,
    renderDateRangeInput,
    renderBooleanInput,
    renderUUIDInput,
  ]);

  /**
   * Get validation errors for a specific field
   */
  const getFieldErrors = useCallback((fieldId: string) => {
    return validationResult.errors.filter(error => error.fieldId === fieldId);
  }, [validationResult.errors]);

  /**
   * Renders the complete filter field with label, input, and validation
   */
  const renderFilterField = useCallback((
    fieldId: string,
    fieldDefinition: FilterFieldDefinition
  ): React.ReactElement => {
    const filterInput = renderFilter(fieldId, fieldDefinition);
    const fieldErrors = showValidationErrors ? getFieldErrors(fieldId) : [];
    const hasErrors = fieldErrors.length > 0;

    // Don't render anything if the filter type is Hidden
    if (!filterInput) {
      return <React.Fragment key={fieldId} />;
    }

    return (
      <div key={fieldId} className="form-control w-full">
        {/* Label with tooltip and required indicator */}
        <div className="label">
          <span className="label-text">
            <div
              className="tooltip tooltip-top"
              data-tip={fieldDefinition.description}
            >
              <span className="cursor-help">
                {fieldDefinition.displayName}
                {fieldDefinition.required && (
                  <span className="text-error ml-1" aria-label="Required field">*</span>
                )}
              </span>
            </div>
          </span>
        </div>

        {/* Filter input with error styling */}
        <div className={clsx(hasErrors && "has-error")}>
          {filterInput}
        </div>

        {/* Validation errors */}
        {showValidationErrors && fieldErrors.length > 0 && (
          <div className="label">
            <span className="label-text-alt text-error">
              {fieldErrors[0].message}
            </span>
          </div>
        )}

        {/* Field description for screen readers */}
        <div id={`${fieldId}-description`} className="sr-only">
          {fieldDefinition.description}
        </div>
      </div>
    );
  }, [renderFilter, getFieldErrors, showValidationErrors]);

  // =============================================================================
  // RENDER
  // =============================================================================

  if (isLoading) {
    return (
      <div className={clsx(className, "flex items-center justify-center p-8")}>
        <span className="loading loading-spinner loading-md mr-2"></span>
        <span>Loading filters...</span>
      </div>
    );
  }

  if (error) {
    return (
      <div className={clsx(className)}>
        <div className="alert alert-error">
          <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>Error loading filters: {error}</span>
        </div>
      </div>
    );
  }

  return (
    <div className={className} role="group" aria-label="Document filters">
      {sortedFields.map((fieldDefinition) =>
        renderFilterField(fieldDefinition.id, fieldDefinition)
      )}

      {/* Global validation summary */}
      {showValidationErrors && !validationResult.isValid && (
        <div className="alert alert-warning mt-4">
          <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
          <div>
            <p className="font-semibold">Please fix the following issues:</p>
            <ul className="list-disc list-inside mt-2">
              {validationResult.errors.slice(0, 3).map((error, index) => (
                <li key={index} className="text-sm">{error.message}</li>
              ))}
              {validationResult.errors.length > 3 && (
                <li className="text-sm">...and {validationResult.errors.length - 3} more</li>
              )}
            </ul>
          </div>
        </div>
      )}
    </div>
  );
}

// =============================================================================
// UTILITY HOOKS
// =============================================================================

/**
 * Custom hook for managing filter state with validation and persistence
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
