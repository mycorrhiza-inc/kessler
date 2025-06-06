/**
 * Dynamic Filter System with Backend Configuration
 * This module provides a flexible filtering system where filter fields
 * and their configurations are loaded from a backend endpoint.
 */

import { contextualApiUrl, getEnvConfig } from "./env_variables/env_variables";

// =============================================================================
// CORE TYPES
// =============================================================================

/**
 * Filter field definition from backend
 */
export interface FilterFieldDefinition {
  /** Unique identifier for the filter field */
  id: string;
  /** Backend field name/key */
  backendKey: string;
  /** Human-readable display name */
  displayName: string;
  /** Detailed description for tooltips */
  description: string;
  /** Input type for rendering */
  inputType: FilterInputType;
  /** Whether this field is required */
  required?: boolean;
  /** Placeholder text for inputs */
  placeholder?: string;
  /** Display order/priority */
  order: number;
  /** Category for grouping */
  category: string;
  /** Validation rules */
  validation?: FilterValidation;
  /** Options for select/multiselect inputs */
  options?: FilterOption[];
  /** Default value */
  defaultValue?: string;
  /** Whether field is currently enabled */
  enabled: boolean;
}

/**
 * Input types supported by the filter system
 */
export enum FilterInputType {
  Text = "text",
  Select = "select",
  MultiSelect = "multiselect",
  Date = "date",
  DateRange = "daterange",
  Number = "number",
  Boolean = "boolean",
  UUID = "uuid",
  Custom = "custom",
  Hidden = "hidden",
}

/**
 * Validation configuration for filter fields
 */
export interface FilterValidation {
  /** Minimum length for text inputs */
  minLength?: number;
  /** Maximum length for text inputs */
  maxLength?: number;
  /** Regular expression pattern */
  pattern?: string;
  /** Minimum value for number inputs */
  min?: number;
  /** Maximum value for number inputs */
  max?: number;
  /** Custom validation function name */
  customValidator?: string;
}

/**
 * Option for select/multiselect inputs
 */
export interface FilterOption {
  /** Option value */
  value: string;
  /** Display label */
  label: string;
  /** Whether option is disabled */
  disabled?: boolean;
  /** Additional metadata */
  metadata?: Record<string, any>;
}

/**
 * Category grouping for filters
 */
export interface FilterCategory {
  /** Category identifier */
  id: string;
  /** Display name */
  name: string;
  /** Description */
  description?: string;
  /** Display order */
  order: number;
  /** Whether category is collapsible in UI */
  collapsible?: boolean;
}

/**
 * Complete filter configuration from backend
 */
export interface FilterConfiguration {
  /** All available filter fields */
  fields: FilterFieldDefinition[];
  /** Filter categories */
  categories: FilterCategory[];
  /** Global configuration */
  config: {
    /** API version */
    version: string;
    /** Last updated timestamp */
    lastUpdated: string;
    /** Default category for uncategorized fields */
    defaultCategory: string;
  };
}

/**
 * Current filter values (dynamic based on loaded configuration)
 */
export type FilterValues = Record<string, string>;

/**
 * Inherited filter value with dynamic field ID
 */
export interface InheritedFilterValue {
  fieldId: string;
  value: any;
}

/**
 * Array of inherited filter values
 */
export type InheritedFilterValues = InheritedFilterValue[];

/**
 * Query data structure
 */
export interface QueryDataFile {
  query: string;
  filters: FilterValues;
}

/**
 * Backend filter object (structure determined by backend)
 */
export interface BackendFilterObject {
  [key: string]: any;
}

// =============================================================================
// FILTER CONFIGURATION MANAGEMENT
// =============================================================================

/**
 * Filter configuration manager class
 * Handles loading, caching, and providing filter configurations
 */
export class FilterConfigurationManager {
  private static instance: FilterConfigurationManager;
  private configuration: FilterConfiguration | null = null;
  private loading: Promise<FilterConfiguration> | null = null;
  private endpoints: FilterEndpoints;
  private cache: Map<string, any> = new Map();

  private constructor(endpoints: FilterEndpoints) {
    this.endpoints = endpoints;
  }

  /**
   * Get singleton instance of the configuration manager
   */
  public static getInstance(endpoints?: FilterEndpoints): FilterConfigurationManager {
    if (!FilterConfigurationManager.instance) {
      if (!endpoints) {
        throw new Error("FilterEndpoints must be provided for first instantiation");
      }
      FilterConfigurationManager.instance = new FilterConfigurationManager(endpoints);
    }
    return FilterConfigurationManager.instance;
  }

  /**
   * Load filter configuration from backend
   */
  public async loadConfiguration(force = false): Promise<FilterConfiguration> {
    // Return cached configuration if available and not forcing reload
    if (this.configuration && !force) {
      return this.configuration;
    }

    // Return existing loading promise if one is in progress
    if (this.loading && !force) {
      return this.loading;
    }

    // Start new loading process
    this.loading = this.fetchConfiguration();

    try {
      this.configuration = await this.loading;
      return this.configuration;
    } finally {
      this.loading = null;
    }
  }

  /**
   * Get current configuration (returns null if not loaded)
   */
  public getConfiguration(): FilterConfiguration | null {
    return this.configuration;
  }

  /**
   * Get filter field by ID
   */
  public getField(fieldId: string): FilterFieldDefinition | null {
    return this.configuration?.fields.find(field => field.id === fieldId) || null;
  }

  /**
   * Get fields by category
   */
  public getFieldsByCategory(categoryId: string): FilterFieldDefinition[] {
    if (!this.configuration) return [];
    return this.configuration.fields.filter(field => field.category === categoryId);
  }

  /**
   * Get enabled fields only
   */
  public getEnabledFields(): FilterFieldDefinition[] {
    if (!this.configuration) return [];
    return this.configuration.fields.filter(field => field.enabled);
  }

  /**
   * Get fields sorted by order
   */
  public getSortedFields(): FilterFieldDefinition[] {
    if (!this.configuration) return [];
    return [...this.configuration.fields].sort((a, b) => a.order - b.order);
  }

  /**
   * Validate filter values against configuration
   */
  public validateFilters(filters: FilterValues): ValidationResult {
    const errors: ValidationError[] = [];
    const warnings: ValidationWarning[] = [];

    if (!this.configuration) {
      errors.push({
        fieldId: '_system',
        message: 'Filter configuration not loaded',
        type: 'system'
      });
      return { isValid: false, errors, warnings };
    }

    // Validate each filter value
    Object.entries(filters).forEach(([fieldId, value]) => {
      const field = this.getField(fieldId);

      if (!field) {
        warnings.push({
          fieldId,
          message: `Unknown filter field: ${fieldId}`,
          type: 'unknown_field'
        });
        return;
      }

      if (!field.enabled) {
        warnings.push({
          fieldId,
          message: `Filter field is disabled: ${fieldId}`,
          type: 'disabled_field'
        });
        return;
      }

      // Validate based on field configuration
      const fieldErrors = this.validateFieldValue(field, value);
      errors.push(...fieldErrors);
    });

    return {
      isValid: errors.length === 0,
      errors,
      warnings
    };
  }

  /**
   * Create empty filter values based on current configuration
   */
  public createEmptyFilters(): FilterValues {
    if (!this.configuration) return {};

    const emptyFilters: FilterValues = {};
    this.configuration.fields.forEach(field => {
      emptyFilters[field.id] = field.defaultValue || "";
    });
    return emptyFilters;
  }

  /**
   * Convert frontend filters to backend format
   */
  public async convertToBackendFilters(filters: FilterValues): Promise<BackendFilterObject> {
    if (!this.endpoints.convertFilters) {
      // Fallback: simple key mapping
      return this.convertFiltersLocally(filters);
    }

    try {
      const response = await fetch(this.endpoints.convertFilters, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ filters })
      });

      if (!response.ok) {
        throw new Error(`Failed to convert filters: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      console.warn('Backend filter conversion failed, using local conversion:', error);
      return this.convertFiltersLocally(filters);
    }
  }

  // =============================================================================
  // PRIVATE METHODS
  // =============================================================================

  private async fetchConfiguration(): Promise<FilterConfiguration> {
    try {
      const response = await fetch(this.endpoints.configuration);

      if (!response.ok) {
        throw new Error(`Failed to load filter configuration: ${response.statusText}`);
      }

      const config = await response.json();

      // Validate configuration structure
      this.validateConfiguration(config);

      return config;
    } catch (error) {
      console.error('Failed to load filter configuration:', error);
      throw new Error(`Filter configuration loading failed: ${error.message}`);
    }
  }

  private validateConfiguration(config: any): void {
    if (!config.fields || !Array.isArray(config.fields)) {
      throw new Error('Invalid configuration: fields array is required');
    }

    if (!config.categories || !Array.isArray(config.categories)) {
      throw new Error('Invalid configuration: categories array is required');
    }

    // Validate each field has required properties
    config.fields.forEach((field: any, index: number) => {
      if (!field.id || !field.backendKey || !field.displayName) {
        throw new Error(`Invalid field at index ${index}: id, backendKey, and displayName are required`);
      }
    });
  }

  private validateFieldValue(field: FilterFieldDefinition, value: string): ValidationError[] {
    const errors: ValidationError[] = [];

    // Required field validation
    if (field.required && (!value || value.trim() === '')) {
      errors.push({
        fieldId: field.id,
        message: `${field.displayName} is required`,
        type: 'required'
      });
      return errors;
    }

    // Skip validation for empty optional fields
    if (!value || value.trim() === '') {
      return errors;
    }

    // Validation based on field type and rules
    if (field.validation) {
      const validation = field.validation;

      // Length validation
      if (validation.minLength && value.length < validation.minLength) {
        errors.push({
          fieldId: field.id,
          message: `${field.displayName} must be at least ${validation.minLength} characters`,
          type: 'minLength'
        });
      }

      if (validation.maxLength && value.length > validation.maxLength) {
        errors.push({
          fieldId: field.id,
          message: `${field.displayName} must be no more than ${validation.maxLength} characters`,
          type: 'maxLength'
        });
      }

      // Pattern validation
      if (validation.pattern) {
        const regex = new RegExp(validation.pattern);
        if (!regex.test(value)) {
          errors.push({
            fieldId: field.id,
            message: `${field.displayName} format is invalid`,
            type: 'pattern'
          });
        }
      }

      // Number validation
      if (field.inputType === FilterInputType.Number) {
        const numValue = parseFloat(value);
        if (isNaN(numValue)) {
          errors.push({
            fieldId: field.id,
            message: `${field.displayName} must be a valid number`,
            type: 'type'
          });
        } else {
          if (validation.min !== undefined && numValue < validation.min) {
            errors.push({
              fieldId: field.id,
              message: `${field.displayName} must be at least ${validation.min}`,
              type: 'min'
            });
          }

          if (validation.max !== undefined && numValue > validation.max) {
            errors.push({
              fieldId: field.id,
              message: `${field.displayName} must be no more than ${validation.max}`,
              type: 'max'
            });
          }
        }
      }
    }

    return errors;
  }

  private convertFiltersLocally(filters: FilterValues): BackendFilterObject {
    if (!this.configuration) {
      return {};
    }

    const backendFilters: BackendFilterObject = {};

    Object.entries(filters).forEach(([fieldId, value]) => {
      const field = this.getField(fieldId);
      if (field && field.enabled && value !== '') {
        backendFilters[field.backendKey] = value;
      }
    });

    return backendFilters;
  }
}

// =============================================================================
// TYPES FOR VALIDATION
// =============================================================================

export interface ValidationError {
  fieldId: string;
  message: string;
  type: string;
}

export interface ValidationWarning {
  fieldId: string;
  message: string;
  type: string;
}

export interface ValidationResult {
  isValid: boolean;
  errors: ValidationError[];
  warnings: ValidationWarning[];
}

// =============================================================================
// ENDPOINT CONFIGURATION
// =============================================================================

/**
 * Backend endpoints configuration
 */
export interface FilterEndpoints {
  /** Endpoint to load filter configuration */
  configuration: string;
  /** Optional endpoint to convert frontend filters to backend format */
  convertFilters?: string;
  /** Optional endpoint to validate filters */
  validateFilters?: string;
  /** Optional endpoint to get filter options dynamically */
  getOptions?: string;
}

// =============================================================================
// UTILITY FUNCTIONS
// =============================================================================

/**
 * Create filter configuration manager instance
 */
export function createFilterManager(endpoints: FilterEndpoints): FilterConfigurationManager {
  return FilterConfigurationManager.getInstance(endpoints);
}


/**
 * Extract field IDs from inherited filters for disabling in UI
 */
export function getDisabledFieldsFromInherited(
  inheritedFilters: InheritedFilterValues
): string[] {
  return inheritedFilters.map((inherited) => inherited.fieldId);
}

/**
 * Create initial filter values from inherited filters
 */
export function createInitialFiltersFromInherited(
  inheritedFilters: InheritedFilterValues,
  configManager: FilterConfigurationManager
): FilterValues {
  const initialFilters = configManager.createEmptyFilters();

  inheritedFilters.forEach((inherited) => {
    if (configManager.getField(inherited.fieldId)) {
      initialFilters[inherited.fieldId] = inherited.value;
    }
  });

  return initialFilters;
}

/**
 * Convert filter values to inherited filter format
 */
export function convertFiltersToInherited(
  filters: FilterValues | null
): InheritedFilterValues {
  if (!filters) {
    return [];
  }

  return Object.entries(filters)
    .filter(([, value]) => value !== "")
    .map(([fieldId, value]) => ({
      fieldId,
      value,
    }));
}

/**
 * Check if any filters have non-empty values
 */
export function hasActiveFilters(filters: FilterValues): boolean {
  return Object.values(filters).some(value => value !== "");
}

/**
 * Clear all filter values
 */
export function clearAllFilters(configManager: FilterConfigurationManager): FilterValues {
  return configManager.createEmptyFilters();
}

/**
 * Clear specific filter fields
 */
export function clearSpecificFilters(
  filters: FilterValues,
  fieldsToClear: string[]
): FilterValues {
  const newFilters = { ...filters };
  fieldsToClear.forEach(fieldId => {
    newFilters[fieldId] = "";
  });
  return newFilters;
}



export const makeFilterEndpoints = (): FilterEndpoints => {
  const contextual_url = contextualApiUrl(getEnvConfig())
  return {
    configuration: `${contextual_url}/api/filters/configuration`,
    convertFilters: `${contextual_url}/api/filters/convert`,
    validateFilters: `${contextual_url}/api/filters/validate`,
    getOptions: `${contextual_url}/api/filters/options`
  }
};

// =============================================================================
// LEGACY COMPATIBILITY
// =============================================================================

// Backward compatibility exports
export const disableListFromInherited = getDisabledFieldsFromInherited;
export const initialFiltersFromInherited = createInitialFiltersFromInherited;
export const inheritedFiltersFromValues = convertFiltersToInherited;
export type FilterField = string;
