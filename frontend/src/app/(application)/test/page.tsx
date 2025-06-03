"use client"

import { useState, useEffect } from "react";
import axios from "axios";

/**
 * Test configuration for legal document filters
 * This provides a realistic example of filter configuration for testing
 */
import {
  FilterConfiguration,
  FilterFieldDefinition,
  FilterInputType,
  FilterCategory,
  FilterEndpoints
} from "@/lib/filters";

// =============================================================================
// TEST FILTER CONFIGURATION
// =============================================================================

/**
 * Complete test filter configuration for legal document system
 */
const testFilterConfiguration: FilterConfiguration = {
  fields: [
    {
      id: "case_number",
      backendKey: "case_number",
      displayName: "Case Number",
      description: "Search by case number or docket number",
      inputType: FilterInputType.Text,
      required: false,
      placeholder: "Enter case number (e.g., 2024-CV-001234)",
      order: 1,
      category: "case_info",
      validation: {
        minLength: 3,
        maxLength: 50,
        pattern: "^[A-Za-z0-9\\-_]+$"
      },
      defaultValue: "",
      enabled: true
    },
    {
      id: "created_at",
      backendKey: "document_created_date",
      displayName: "Document Created Date",
      description: "Select the date when the document was created in the system",
      inputType: FilterInputType.Date,
      required: false,
      placeholder: "",
      order: 2,
      category: "dates",
      validation: {},
      defaultValue: "",
      enabled: true
    },
    {
      id: "filed_date",
      backendKey: "court_filing_date",
      displayName: "Filed Date",
      description: "Date when the document was filed with the court",
      inputType: FilterInputType.Date,
      required: false,
      placeholder: "",
      order: 3,
      category: "dates",
      validation: {},
      defaultValue: "",
      enabled: true
    },
    {
      id: "filing_type",
      backendKey: "filing_type_codes",
      displayName: "Filing Type",
      description: "Type of court filing or document category",
      inputType: FilterInputType.MultiSelect,
      required: false,
      placeholder: "Select filing types...",
      order: 4,
      category: "document_types",
      validation: {},
      options: [
        { value: "motion", label: "Motion", disabled: false },
        { value: "pleading", label: "Pleading", disabled: false },
        { value: "brief", label: "Brief", disabled: false },
        { value: "order", label: "Court Order", disabled: false },
        { value: "judgment", label: "Judgment", disabled: false },
        { value: "discovery", label: "Discovery Request", disabled: false },
        { value: "response", label: "Discovery Response", disabled: false },
        { value: "deposition", label: "Deposition", disabled: false },
        { value: "expert_report", label: "Expert Report", disabled: false },
        { value: "settlement", label: "Settlement Agreement", disabled: false },
        { value: "appeal", label: "Appeal Document", disabled: false },
        { value: "administrative", label: "Administrative Filing", disabled: false }
      ],
      defaultValue: "",
      enabled: true
    },
    {
      id: "matter_subtype",
      backendKey: "matter_subtype_ids",
      displayName: "Matter Subtype",
      description: "Specific legal matter subtypes (searchable)",
      inputType: FilterInputType.MultiSelect,
      required: false,
      placeholder: "Search and select matter subtypes...",
      order: 5,
      category: "matter_classification",
      validation: {},
      options: [
        // Civil Litigation Subtypes
        { value: "contract_dispute", label: "Contract Dispute", disabled: false },
        { value: "employment_litigation", label: "Employment Litigation", disabled: false },
        { value: "personal_injury", label: "Personal Injury", disabled: false },
        { value: "property_dispute", label: "Property Dispute", disabled: false },
        { value: "intellectual_property", label: "Intellectual Property", disabled: false },
        { value: "securities_litigation", label: "Securities Litigation", disabled: false },

        // Corporate/Business Subtypes
        { value: "merger_acquisition", label: "Merger & Acquisition", disabled: false },
        { value: "corporate_governance", label: "Corporate Governance", disabled: false },
        { value: "regulatory_compliance", label: "Regulatory Compliance", disabled: false },
        { value: "tax_matters", label: "Tax Matters", disabled: false },

        // Criminal Subtypes
        { value: "white_collar", label: "White Collar Crime", disabled: false },
        { value: "fraud_investigation", label: "Fraud Investigation", disabled: false },
        { value: "regulatory_enforcement", label: "Regulatory Enforcement", disabled: false },

        // Family Law Subtypes
        { value: "divorce", label: "Divorce Proceedings", disabled: false },
        { value: "child_custody", label: "Child Custody", disabled: false },
        { value: "adoption", label: "Adoption", disabled: false },

        // Real Estate Subtypes
        { value: "commercial_real_estate", label: "Commercial Real Estate", disabled: false },
        { value: "residential_real_estate", label: "Residential Real Estate", disabled: false },
        { value: "zoning_land_use", label: "Zoning & Land Use", disabled: false }
      ],
      defaultValue: "",
      enabled: true
    },
    {
      id: "matter_type",
      backendKey: "matter_type_categories",
      displayName: "Matter Type",
      description: "Primary legal matter categories (searchable)",
      inputType: FilterInputType.MultiSelect,
      required: false,
      placeholder: "Search and select matter types...",
      order: 6,
      category: "matter_classification",
      validation: {},
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
        { value: "environmental", label: "Environmental Law", disabled: false },
        { value: "healthcare", label: "Healthcare Law", disabled: false },
        { value: "securities", label: "Securities Law", disabled: false },
        { value: "antitrust", label: "Antitrust Law", disabled: false },
        { value: "administrative", label: "Administrative Law", disabled: false }
      ],
      defaultValue: "",
      enabled: true
    },
    {
      id: "party_name",
      backendKey: "party_names",
      displayName: "Party Name",
      description: "Names of parties involved in the case (searchable)",
      inputType: FilterInputType.MultiSelect,
      required: false,
      placeholder: "Search and select party names...",
      order: 7,
      category: "parties",
      validation: {
        minLength: 2
      },
      options: [
        // Example parties - in a real system, these would be loaded dynamically
        { value: "acme_corp", label: "ACME Corporation", disabled: false },
        { value: "smith_john", label: "John Smith", disabled: false },
        { value: "doe_jane", label: "Jane Doe", disabled: false },
        { value: "global_tech_inc", label: "Global Tech Inc.", disabled: false },
        { value: "city_springfield", label: "City of Springfield", disabled: false },
        { value: "johnson_mary", label: "Mary Johnson", disabled: false },
        { value: "brown_robert", label: "Robert Brown", disabled: false },
        { value: "mega_corp_ltd", label: "MegaCorp Ltd.", disabled: false },
        { value: "williams_sarah", label: "Sarah Williams", disabled: false },
        { value: "tech_innovations", label: "Tech Innovations LLC", disabled: false },
        { value: "state_california", label: "State of California", disabled: false },
        { value: "jones_michael", label: "Michael Jones", disabled: false },
        { value: "davis_lisa", label: "Lisa Davis", disabled: false },
        { value: "enterprise_solutions", label: "Enterprise Solutions Inc.", disabled: false },
        { value: "wilson_david", label: "David Wilson", disabled: false }
      ],
      defaultValue: "",
      enabled: true
    }
  ],
  categories: [
    {
      id: "case_info",
      name: "Case Information",
      description: "Basic case identification and reference information",
      order: 1,
      collapsible: false
    },
    {
      id: "dates",
      name: "Important Dates",
      description: "Date-related filters for document timing",
      order: 2,
      collapsible: true
    },
    {
      id: "document_types",
      name: "Document Types",
      description: "Classification by document and filing types",
      order: 3,
      collapsible: true
    },
    {
      id: "matter_classification",
      name: "Matter Classification",
      description: "Legal matter types and subtypes",
      order: 4,
      collapsible: true
    },
    {
      id: "parties",
      name: "Parties & Participants",
      description: "People and entities involved in the case",
      order: 5,
      collapsible: true
    }
  ],
  config: {
    version: "1.0.0",
    lastUpdated: "2025-06-03T12:00:00Z",
    defaultCategory: "case_info"
  }
};

// =============================================================================
// TEST ENDPOINTS CONFIGURATION
// =============================================================================

/**
 * Mock endpoints for testing (you can replace these with real backend URLs)
 */
const testFilterEndpoints: FilterEndpoints = {
  configuration: "/api/filters/configuration",
  convertFilters: "/api/filters/convert",
  validateFilters: "/api/filters/validate",
  getOptions: "/api/filters/options"
};

// =============================================================================
// MOCK API FUNCTIONS FOR TESTING
// =============================================================================

/**
 * Mock API function that simulates loading filter configuration from backend
 * In a real application, this would be handled by your backend
 */
export async function mockLoadFilterConfiguration(): Promise<FilterConfiguration> {
  // Simulate network delay
  await new Promise(resolve => setTimeout(resolve, 500));
  return testFilterConfiguration;
}

/**
 * Mock API function for converting frontend filters to backend format
 */
export async function mockConvertFilters(filters: Record<string, string>): Promise<any> {
  // Simulate network delay
  await new Promise(resolve => setTimeout(resolve, 200));

  // Example conversion logic
  const backendFilters: any = {
    metadata_filters: {},
    search_filters: {},
    date_filters: {},
    multiselect_filters: {}
  };

  Object.entries(filters).forEach(([fieldId, value]) => {
    if (!value) return;

    const field = testFilterConfiguration.fields.find(f => f.id === fieldId);
    if (!field) return;

    switch (field.inputType) {
      case FilterInputType.Text:
        backendFilters.metadata_filters[field.backendKey] = value;
        break;
      case FilterInputType.Date:
        backendFilters.date_filters[field.backendKey] = value;
        break;
      case FilterInputType.MultiSelect:
        backendFilters.multiselect_filters[field.backendKey] = value.split(',').filter(Boolean);
        break;
    }
  });

  return backendFilters;
}

/**
 * Mock API function for validating filters
 */
export async function mockValidateFilters(filters: Record<string, string>): Promise<any> {
  // Simulate network delay
  await new Promise(resolve => setTimeout(resolve, 100));

  const errors: any[] = [];
  const warnings: any[] = [];

  // Example validation logic
  Object.entries(filters).forEach(([fieldId, value]) => {
    const field = testFilterConfiguration.fields.find(f => f.id === fieldId);
    if (!field || !value) return;

    // Check required fields
    if (field.required && !value.trim()) {
      errors.push({
        fieldId,
        message: `${field.displayName} is required`,
        type: 'required'
      });
    }

    // Check pattern validation
    if (field.validation?.pattern && value) {
      const regex = new RegExp(field.validation.pattern);
      if (!regex.test(value)) {
        errors.push({
          fieldId,
          message: `${field.displayName} format is invalid`,
          type: 'pattern'
        });
      }
    }
  });

  return {
    isValid: errors.length === 0,
    errors,
    warnings
  };
}

// =============================================================================
// HELPER FUNCTIONS FOR TESTING
// =============================================================================

/**
 * Get all field IDs for testing
 */
export function getAllTestFieldIds(): string[] {
  return testFilterConfiguration.fields.map(field => field.id);
}

/**
 * Get field IDs by category for testing
 */
export function getTestFieldIdsByCategory(categoryId: string): string[] {
  return testFilterConfiguration.fields
    .filter(field => field.category === categoryId)
    .map(field => field.id);
}

/**
 * Create sample filter values for testing
 */
export function createSampleFilterValues(): Record<string, string> {
  return {
    case_number: "2024-CV-001234",
    created_at: "2024-01-15",
    filed_date: "2024-01-10",
    filing_type: "motion,brief",
    matter_subtype: "contract_dispute,employment_litigation",
    matter_type: "litigation,employment",
    party_name: "acme_corp,smith_john"
  };
}

/**
 * Setup mock fetch for testing - Browser compatible version
 */
export function setupMockFetch(): void {
  // Store original fetch
  const originalFetch = window.fetch;

  // Create mock implementation
  window.fetch = async (input: RequestInfo | URL, init?: RequestInit): Promise<Response> => {
    const url = typeof input === 'string' ? input : input.toString();

    if (url.includes('/api/filters/configuration')) {
      return new Response(JSON.stringify(testFilterConfiguration), {
        status: 200,
        headers: { 'Content-Type': 'application/json' }
      });
    }

    if (url.includes('/api/filters/convert')) {
      const body = init?.body ? JSON.parse(init.body as string) : {};
      const result = await mockConvertFilters(body.filters || {});
      return new Response(JSON.stringify(result), {
        status: 200,
        headers: { 'Content-Type': 'application/json' }
      });
    }

    if (url.includes('/api/filters/validate')) {
      const body = init?.body ? JSON.parse(init.body as string) : {};
      const result = await mockValidateFilters(body.filters || {});
      return new Response(JSON.stringify(result), {
        status: 200,
        headers: { 'Content-Type': 'application/json' }
      });
    }

    // Fallback to original fetch for other URLs
    return originalFetch(input, init);
  };
}

/**
 * Cleanup mock fetch after testing
 */
export function cleanupMockFetch(): void {
  // This would restore original fetch if we stored it
  // For now, just reload the page or implement proper restoration
}

// =============================================================================
// REAL BACKEND INTEGRATION EXAMPLE
// =============================================================================

/**
 * Real backend integration for fetching filters
 */
function useRealBackendFilters() {
  const [filters, setFilters] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const getFilters = async (root?: string) => {
    console.log("Getting filters from real backend");
    setLoading(true);
    setError(null);

    try {
      const response = await axios.get(`http://localhost/fugu/filters/metadata`);
      console.log("Backend response:", response.data);
      setFilters(response.data.filters || []);
    } catch (error) {
      console.error('Error fetching filters:', error);
      setError(error instanceof Error ? error.message : 'Failed to fetch filters');
    } finally {
      setLoading(false);
    }
  };

  const filterTail = (facet: string): string => {
    const parts = facet.split("/");
    return parts[parts.length - 1];
  };

  return {
    filters,
    loading,
    error,
    getFilters,
    filterTail
  };
}

// =============================================================================
// IMPORTS - CANONICAL MULTISELECT COMPONENT
// =============================================================================

import {
  DynamicDocumentFiltersList,
  DynamicDocumentFiltersGrid,
  ResponsiveDynamicDocumentFilters,
  useDocumentFilters
} from '@/components/Filters/DocumentFilters';

// Import the canonical multiselect component
import { DynamicMultiSelect } from '@/components/Filters/FilterMultiSelect';

import {
  FilterValues,
  FilterConfigurationManager,
  createFilterManager,
  ValidationResult
} from '@/lib/filters';

// =============================================================================
// CANONICAL MULTISELECT TEST COMPONENT
// =============================================================================

/**
 * Canonical MultiSelect Test Component - Uses Direct Import
 */
function CanonicalMultiSelectTest() {
  const [filingTypeValue, setFilingTypeValue] = useState("motion,brief");
  const [matterTypeValue, setMatterTypeValue] = useState("litigation,corporate");
  const [partyNameValue, setPartyNameValue] = useState("acme_corp,smith_john");

  // Get field definitions from test configuration
  const filingTypeField = testFilterConfiguration.fields.find(f => f.id === 'filing_type')!;
  const matterTypeField = testFilterConfiguration.fields.find(f => f.id === 'matter_type')!;
  const partyNameField = testFilterConfiguration.fields.find(f => f.id === 'party_name')!;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title">🎯 Canonical MultiSelect Component Test</h2>
          <div className="alert alert-info">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" className="stroke-current shrink-0 w-6 h-6">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <div>
              <p className="font-semibold">Direct import from @/components/Filters/FilterMultiSelect</p>
              <p className="text-sm mt-1">Testing the enhanced DynamicMultiSelect component with all features</p>
            </div>
          </div>
        </div>
      </div>

      {/* Individual MultiSelect Components */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Filing Type */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h3 className="card-title text-lg">Filing Type</h3>
            <DynamicMultiSelect
              fieldDefinition={filingTypeField}
              value={filingTypeValue}
              onChange={setFilingTypeValue}
              onFocus={() => console.log('🎯 Filing Type focused')}
              onBlur={() => console.log('👋 Filing Type blurred')}
            />
            <div className="mt-3 text-xs">
              <strong>Current value:</strong>
              <div className="font-mono bg-gray-100 p-2 rounded mt-1 break-all">
                "{filingTypeValue}"
              </div>
            </div>
          </div>
        </div>

        {/* Matter Type */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h3 className="card-title text-lg">Matter Type</h3>
            <DynamicMultiSelect
              fieldDefinition={matterTypeField}
              value={matterTypeValue}
              onChange={setMatterTypeValue}
              onFocus={() => console.log('🎯 Matter Type focused')}
              onBlur={() => console.log('👋 Matter Type blurred')}
            />
            <div className="mt-3 text-xs">
              <strong>Current value:</strong>
              <div className="font-mono bg-gray-100 p-2 rounded mt-1 break-all">
                "{matterTypeValue}"
              </div>
            </div>
          </div>
        </div>

        {/* Party Name */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h3 className="card-title text-lg">Party Name</h3>
            <DynamicMultiSelect
              fieldDefinition={partyNameField}
              value={partyNameValue}
              onChange={setPartyNameValue}
              onFocus={() => console.log('🎯 Party Name focused')}
              onBlur={() => console.log('👋 Party Name blurred')}
            />
            <div className="mt-3 text-xs">
              <strong>Current value:</strong>
              <div className="font-mono bg-gray-100 p-2 rounded mt-1 break-all">
                "{partyNameValue}"
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Test Controls */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h3 className="card-title">🧪 Test Controls</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
            <button
              onClick={() => setFilingTypeValue("motion,brief,order,judgment")}
              className="btn btn-sm btn-primary"
            >
              Set Filing Types
            </button>
            <button
              onClick={() => setMatterTypeValue("litigation,corporate,criminal,family")}
              className="btn btn-sm btn-secondary"
            >
              Set Matter Types
            </button>
            <button
              onClick={() => setPartyNameValue("acme_corp,smith_john,global_tech_inc,city_springfield")}
              className="btn btn-sm btn-accent"
            >
              Set Party Names
            </button>
            <button
              onClick={() => {
                setFilingTypeValue("");
                setMatterTypeValue("");
                setPartyNameValue("");
              }}
              className="btn btn-sm btn-outline"
            >
              Clear All
            </button>
          </div>

          {/* Combined Values Display */}
          <div className="mt-6 p-4 bg-gray-50 rounded-lg">
            <h4 className="font-semibold mb-3">📊 Combined State</h4>
            <div className="space-y-2 text-sm">
              <div>
                <span className="font-medium">Filing Types:</span>
                <div className="flex flex-wrap gap-1 mt-1">
                  {filingTypeValue.split(',').filter(Boolean).map((item, index) => (
                    <span key={index} className="badge badge-primary badge-sm">{item}</span>
                  ))}
                </div>
              </div>
              <div>
                <span className="font-medium">Matter Types:</span>
                <div className="flex flex-wrap gap-1 mt-1">
                  {matterTypeValue.split(',').filter(Boolean).map((item, index) => (
                    <span key={index} className="badge badge-secondary badge-sm">{item}</span>
                  ))}
                </div>
              </div>
              <div>
                <span className="font-medium">Party Names:</span>
                <div className="flex flex-wrap gap-1 mt-1">
                  {partyNameValue.split(',').filter(Boolean).map((item, index) => (
                    <span key={index} className="badge badge-accent badge-sm">{item}</span>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Feature Testing Guide */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h3 className="card-title">✨ Feature Testing Guide</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">

            {/* Button Features */}
            <div>
              <h4 className="font-semibold text-lg mb-3">🔘 Button Features</h4>
              <ul className="space-y-2 text-sm">
                <li className="flex items-start gap-2">
                  <span className="text-green-600">✅</span>
                  <span><strong>Filter name with count:</strong> Shows "Filing Type: 2 selected"</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-green-600">✅</span>
                  <span><strong>Pills in button:</strong> Each selection shown as removable pill</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-green-600">✅</span>
                  <span><strong>Remove buttons:</strong> X button on each pill for quick removal</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-green-600">✅</span>
                  <span><strong>Responsive wrapping:</strong> Pills wrap when many items selected</span>
                </li>
              </ul>
            </div>

            {/* Dropdown Features */}
            <div>
              <h4 className="font-semibold text-lg mb-3">📋 Dropdown Features</h4>
              <ul className="space-y-2 text-sm">
                <li className="flex items-start gap-2">
                  <span className="text-blue-600">🔍</span>
                  <span><strong>Search functionality:</strong> Type to filter options instantly</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-blue-600">🏷️</span>
                  <span><strong>Pills in rows:</strong> Selected items show pills on the right</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-blue-600">🗑️</span>
                  <span><strong>Remove from rows:</strong> Click X on pills to deselect</span>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-blue-600">📊</span>
                  <span><strong>Footer counter:</strong> Shows count and "Clear all" button</span>
                </li>
              </ul>
            </div>
          </div>

          {/* Test Instructions */}
          <div className="mt-6 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
            <h4 className="font-semibold text-yellow-800 mb-2">🧪 Test Instructions:</h4>
            <ol className="list-decimal list-inside space-y-1 text-sm text-yellow-700">
              <li>Click on any multiselect dropdown to open it</li>
              <li>Try typing in the search box to filter options</li>
              <li>Select/deselect items using checkboxes</li>
              <li>Use X buttons on pills to remove items (both in button and rows)</li>
              <li>Try the "Clear all" button in the dropdown footer</li>
              <li>Click outside the dropdown to close it</li>
              <li>Use test control buttons above to quickly set different states</li>
              <li>Check the console for focus/blur event logging</li>
            </ol>
          </div>
        </div>
      </div>
    </div>
  );
}

// =============================================================================
// OTHER EXAMPLE COMPONENTS (Simplified to focus on canonical component)
// =============================================================================

/**
 * Backend Integration Test Component
 */
function RealBackendTest() {
  const { filters, loading, error, getFilters, filterTail } = useRealBackendFilters();

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h2 className="card-title">Real Backend Integration Test</h2>

        <button
          className={`btn btn-accent ${loading ? 'loading' : ''}`}
          onClick={() => getFilters()}
          disabled={loading}
        >
          {loading ? 'Loading...' : 'Get Filters from Backend'}
        </button>

        {error && (
          <div className="alert alert-error mt-4">
            <span>Error: {error}</span>
          </div>
        )}

        {filters.length > 0 && (
          <div className="mt-4">
            <h3 className="text-lg font-semibold mb-2">Backend Filters ({filters.length}):</h3>
            <div className="max-h-60 overflow-y-auto">
              <ul className="list-disc list-inside space-y-1">
                {filters.map((filter, index) => (
                  <li key={index} className="text-sm">
                    <span className="font-medium">{filterTail(filter[0])}</span>
                    {filter[1] && <span className="text-gray-500 ml-2">({filter[1]})</span>}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

/**
 * Integrated Filter System Test - Uses Dynamic Document Filters
 */
function IntegratedFilterSystemTest() {
  const [filters, setFilters] = useState<FilterValues>({});

  useEffect(() => {
    setupMockFetch();
    return () => cleanupMockFetch();
  }, []);

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h2 className="card-title">Integrated Filter System Test</h2>
        <p className="text-gray-600 mb-4">
          Test the complete filter system using DynamicDocumentFiltersList with enhanced multiselect
        </p>

        {/* Filter Component - Uses the canonical DynamicMultiSelect internally */}
        <DynamicDocumentFiltersList
          queryOptions={filters}
          setQueryOptions={setFilters}
          showFields={getAllTestFieldIds()}
          endpoints={testFilterEndpoints}
          onFilterChange={(fieldId, value) => {
            console.log(`🔄 Integrated Filter ${fieldId} changed:`, value);
            if (fieldId.includes('type') || fieldId.includes('filing') || fieldId.includes('party')) {
              console.log(`📊 MultiSelect items: [${value.split(',').filter(Boolean).join(', ')}]`);
            }
          }}
          showValidationErrors={true}
        />

        {/* Current State Display */}
        <div className="mt-6 p-4 bg-gray-50 rounded-lg">
          <h3 className="font-semibold mb-3">📊 Current Filter State</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {Object.entries(filters)
              .filter(([, value]) => value !== '')
              .map(([key, value]) => (
                <div key={key} className="border rounded p-3 bg-white">
                  <div className="font-semibold text-sm text-gray-700 mb-1 capitalize">
                    {key.replace(/_/g, ' ')}:
                  </div>
                  {key.includes('type') || key.includes('filing') || key.includes('party') ? (
                    // MultiSelect fields - show as pills
                    <div className="flex flex-wrap gap-1">
                      {value.split(',').filter(Boolean).map((item, index) => (
                        <span key={index} className="badge badge-primary badge-sm">
                          {item}
                        </span>
                      ))}
                    </div>
                  ) : (
                    // Regular fields - show as text
                    <div className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">
                      {value}
                    </div>
                  )}
                </div>
              ))}
          </div>

          {Object.keys(filters).filter(key => filters[key] !== '').length === 0 && (
            <div className="text-gray-500 italic text-center py-4">
              No filters applied yet. Try using the filters above!
            </div>
          )}
        </div>

        {/* Quick Actions */}
        <div className="mt-4">
          <h4 className="font-semibold mb-2">🧪 Quick Test Actions:</h4>
          <div className="flex gap-2 flex-wrap">
            <button
              onClick={() => setFilters(createSampleFilterValues())}
              className="btn btn-sm btn-primary"
            >
              Load Sample Data
            </button>
            <button
              onClick={() => setFilters(prev => ({
                ...prev,
                filing_type: "motion,brief,order,judgment,discovery"
              }))}
              className="btn btn-sm btn-secondary"
            >
              Many Filing Types
            </button>
            <button
              onClick={() => setFilters({})}
              className="btn btn-sm btn-outline"
            >
              Clear All Filters
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

// =============================================================================
// MAIN TEST PAGE COMPONENT
// =============================================================================

export default function TestPage() {
  const [activeExample, setActiveExample] = useState<string>('canonical');

  const examples = [
    { id: 'canonical', name: 'Canonical MultiSelect', component: CanonicalMultiSelectTest },
    { id: 'integrated', name: 'Integrated System', component: IntegratedFilterSystemTest },
    { id: 'backend', name: 'Real Backend Test', component: RealBackendTest },
  ];

  const ActiveComponent = examples.find(ex => ex.id === activeExample)?.component || CanonicalMultiSelectTest;

  return (
    <div className="container mx-auto p-6 space-y-6 max-w-6xl">
      <div className="text-center">
        <h1 className="text-4xl font-bold mb-2">Enhanced MultiSelect Test Page</h1>
        <p className="text-gray-600">Test the canonical DynamicMultiSelect component with all enhanced features</p>
      </div>

      {/* Example Selector */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title mb-4">Select Test Example</h2>
          <div className="flex flex-wrap gap-2">
            {examples.map(example => (
              <button
                key={example.id}
                className={`btn ${activeExample === example.id ? 'btn-primary' : 'btn-outline'}`}
                onClick={() => setActiveExample(example.id)}
              >
                {example.name}
              </button>
            ))}
          </div>

          {/* Current Example Description */}
          <div className="mt-4 p-3 bg-blue-50 rounded-lg">
            {activeExample === 'canonical' && (
              <p className="text-blue-800 text-sm">
                <strong>Canonical MultiSelect:</strong> Direct testing of the DynamicMultiSelect component
                imported from @/components/Filters/FilterMultiSelect with all enhanced features.
              </p>
            )}
            {activeExample === 'integrated' && (
              <p className="text-blue-800 text-sm">
                <strong>Integrated System:</strong> Testing the complete filter system using
                DynamicDocumentFiltersList which internally uses the canonical DynamicMultiSelect.
              </p>
            )}
            {activeExample === 'backend' && (
              <p className="text-blue-800 text-sm">
                <strong>Real Backend Test:</strong> Integration test with your actual backend API
                at http://localhost/fugu/filters/metadata.
              </p>
            )}
          </div>
        </div>
      </div>

      {/* Active Example */}
      <ActiveComponent />

      {/* Technical Information */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title">🔧 Technical Information</h2>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Import Information */}
            <div>
              <h3 className="text-lg font-semibold mb-3">📦 Component Import</h3>
              <div className="bg-gray-100 p-3 rounded-lg font-mono text-sm">
                <div className="text-green-600">// Canonical import</div>
                <div>import &#123; DynamicMultiSelect &#125;</div>
                <div className="ml-4">from '@/components/Filters/FilterMultiSelect';</div>
              </div>
              <div className="mt-3 text-sm text-gray-600">
                This import ensures you're using the enhanced DynamicMultiSelect
                component with all the search and pill functionality.
              </div>
            </div>

            {/* Configuration Data */}
            <div>
              <h3 className="text-lg font-semibold mb-3">⚙️ Configuration</h3>
              <div className="stats stats-vertical shadow-sm">
                <div className="stat py-2">
                  <div className="stat-title text-xs">MultiSelect Fields</div>
                  <div className="stat-value text-lg">
                    {testFilterConfiguration.fields.filter(f => f.inputType === FilterInputType.MultiSelect).length}
                  </div>
                </div>
                <div className="stat py-2">
                  <div className="stat-title text-xs">Total Options</div>
                  <div className="stat-value text-lg">
                    {testFilterConfiguration.fields
                      .filter(f => f.inputType === FilterInputType.MultiSelect)
                      .reduce((sum, field) => sum + (field.options?.length || 0), 0)}
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Enhanced Features List */}
          <div className="mt-6">
            <h3 className="text-lg font-semibold mb-3">✨ Enhanced Features Verified</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">

              {/* Button Features */}
              <div className="p-4 border rounded-lg">
                <h4 className="font-semibold text-green-700 mb-2">🔘 Button Features</h4>
                <ul className="space-y-1 text-sm">
                  <li className="flex items-center gap-2">
                    <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                    Filter name with selection count display
                  </li>
                  <li className="flex items-center gap-2">
                    <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                    Selected items shown as removable pills
                  </li>
                  <li className="flex items-center gap-2">
                    <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                    Individual X buttons on each pill
                  </li>
                  <li className="flex items-center gap-2">
                    <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                    Responsive pill wrapping
                  </li>
                </ul>
              </div>

              {/* Dropdown Features */}
              <div className="p-4 border rounded-lg">
                <h4 className="font-semibold text-blue-700 mb-2">📋 Dropdown Features</h4>
                <ul className="space-y-1 text-sm">
                  <li className="flex items-center gap-2">
                    <span className="w-2 h-2 bg-blue-500 rounded-full"></span>
                    Real-time search with instant filtering
                  </li>
                  <li className="flex items-center gap-2">
                    <span className="w-2 h-2 bg-blue-500 rounded-full"></span>
                    Pills in option rows for selected items
                  </li>
                  <li className="flex items-center gap-2">
                    <span className="w-2 h-2 bg-blue-500 rounded-full"></span>
                    Remove buttons on row pills
                  </li>
                  <li className="flex items-center gap-2">
                    <span className="w-2 h-2 bg-blue-500 rounded-full"></span>
                    Footer with count and "Clear all" button
                  </li>
                </ul>
              </div>
            </div>
          </div>

          {/* Usage Tips */}
          <div className="mt-6 p-4 bg-amber-50 border border-amber-200 rounded-lg">
            <h4 className="font-semibold text-amber-800 mb-2">💡 Usage Tips</h4>
            <ul className="text-sm text-amber-700 space-y-1">
              <li>• <strong>Search:</strong> Type in the search box to filter options instantly</li>
              <li>• <strong>Select:</strong> Click checkboxes or anywhere on the row to select</li>
              <li>• <strong>Remove:</strong> Use X buttons on pills (in button or rows) to deselect</li>
              <li>• <strong>Clear All:</strong> Use the footer button to clear all selections</li>
              <li>• <strong>Close:</strong> Click outside the dropdown to close it</li>
              <li>• <strong>Console:</strong> Check browser console for detailed event logging</li>
            </ul>
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className="text-center p-4 text-gray-500">
        <p>Enhanced MultiSelect Test Page • Built with React + TypeScript + Tailwind CSS</p>
        <p className="text-sm mt-1">
          Using canonical DynamicMultiSelect from @/components/Filters/FilterMultiSelect
        </p>
      </div>
    </div>
  );
}
 
 
 
 
 
 
  
 
  
 
 
 
 
 
 
  
  
  
  
 
 
 
 
 
  
 
  
 
 
 
 
 
  
  
   
 
  
  
 
 
  
 
  
 
 
 
 
 
 
 
  
 
  
 
 
 
 
 
 
  
  
  
  
 
 
 
 
 
  
 
  
 
 
 
 
 
  
  
   
 
  
  
 
 
  
 
  
 
 
 
 
 
 
 
  
 
  
 
 
 
 
 
 
  
  
  
  
 
 
 
 
 
  
 
  
 
 
 
 
 
  
  
   
 
  
  
 
 
  
 
  
 
 
 
 
 
 
 
  
 
  
 
 
 
 
 
 
  
  
  
  
 
 
 
 
 
  
 
  
 
 
 
 
 
  
  
   
 
  
  
 
 
  
 
  
 
 
 
 
 
 
 
  
 
  
 
 
 
 
 
 
  
  
  
  
 
 
 
 
 
  
 
  
 
 
 
 
 
  
  
   
 
  
  
 
 
  
 
  
 
 
 
 
 
 
 
  
 
  
 
 
 
 
 
 
  
  
  
  
 
 
 
 
 
  
 
  
 
 
 
 
 
  
  
   
 
  
  
 
 
  
 
  
 
 
 
 
 
 
 
  
 
  
 
 
 
 
 
 
  
  
  
  
 
 
 
 
 
  
 
  
 
 
 
 
 
  
  
   
 
  
  
 
 
  
 
  
 
 
 
 
 
 
 
  
 
  
 
 
 
 
 
 
  
  
  
  
 
 
 
 
 
  
 
  
 
 
 
 
 
  
  
   
 
  
  
 
 
  
 
  
 
