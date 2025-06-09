"use client"

import { useState, useEffect } from "react";

// Import types and utilities
import {
  FilterConfiguration,
  FilterFieldDefinition,
  FilterInputType,
  FilterEndpoints,
  createFilterManager,
} from "@/lib/filters";

// Import Kessler-integrated components
import {
  KesslerDocumentFiltersList,
  KesslerDocumentFiltersGrid,
  KesslerResponsiveDynamicDocumentFilters,
  KesslerInlineDynamicDocumentFilters,
  KesslerFilterControls,
  KesslerFilterStatus,
  useKesslerFilters,
  useKesslerFilterField,

} from '@/components/Filters/DynamicFilters';

import { FilterErrorBoundary } from "@/components/Filters/Errors";
// Import Kessler store hooks
import {
  useFilters,
  useFilterSystemLifecycle,
  useUrlSync,
  useFilterPersistence,
  getFilterSystemStatus,
  useFilterPerformance,
  useMemoryMonitor,
  safeBulkUpdateFilters,
} from "@/lib/store";

// Import the direct multiselect for comparison
import { DynamicMultiSelect } from '@/components/Filters/FilterMultiSelect';

// =============================================================================
// ENHANCED TEST CONFIGURATION WITH VALIDATION
// =============================================================================

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
      validation: { minLength: 3, maxLength: 50 },
      defaultValue: "",
      enabled: true
    },
    {
      id: "created_at",
      backendKey: "document_created_date",
      displayName: "Document Created Date",
      description: "Select the date when the document was created",
      inputType: FilterInputType.Date,
      required: false,
      order: 2,
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
      order: 3,
      category: "document_types",
      options: [
        { value: "motion", label: "Motion", disabled: false },
        { value: "pleading", label: "Pleading", disabled: false },
        { value: "brief", label: "Brief", disabled: false },
        { value: "order", label: "Court Order", disabled: false },
        { value: "judgment", label: "Judgment", disabled: false },
        { value: "discovery", label: "Discovery Request", disabled: false },
        { value: "response", label: "Discovery Response", disabled: false },
        { value: "deposition", label: "Deposition", disabled: false },
      ],
      defaultValue: "",
      enabled: true
    },
    {
      id: "matter_type",
      backendKey: "matter_type_categories",
      displayName: "Matter Type",
      description: "Primary legal matter categories",
      inputType: FilterInputType.MultiSelect,
      required: false,
      placeholder: "Search and select matter types...",
      order: 4,
      category: "matter_classification",
      options: [
        { value: "litigation", label: "Litigation", disabled: false },
        { value: "corporate", label: "Corporate Law", disabled: false },
        { value: "criminal", label: "Criminal Law", disabled: false },
        { value: "family", label: "Family Law", disabled: false },
        { value: "real_estate", label: "Real Estate Law", disabled: false },
        { value: "employment", label: "Employment Law", disabled: false },
        { value: "intellectual_property", label: "Intellectual Property", disabled: false },
        { value: "tax", label: "Tax Law", disabled: false },
      ],
      defaultValue: "",
      enabled: true
    },
    {
      id: "party_name",
      backendKey: "party_names",
      displayName: "Party Name",
      description: "Names of parties involved in the case",
      inputType: FilterInputType.MultiSelect,
      required: false,
      placeholder: "Search and select party names...",
      order: 5,
      category: "parties",
      options: [
        { value: "acme_corp", label: "ACME Corporation", disabled: false },
        { value: "smith_john", label: "John Smith", disabled: false },
        { value: "doe_jane", label: "Jane Doe", disabled: false },
        { value: "global_tech_inc", label: "Global Tech Inc.", disabled: false },
        { value: "city_springfield", label: "City of Springfield", disabled: false },
        { value: "johnson_mary", label: "Mary Johnson", disabled: false },
      ],
      defaultValue: "",
      enabled: true
    }
  ],
  categories: [
    { id: "case_info", name: "Case Information", order: 1 },
    { id: "dates", name: "Important Dates", order: 2 },
    { id: "document_types", name: "Document Types", order: 3 },
    { id: "matter_classification", name: "Matter Classification", order: 4 },
    { id: "parties", name: "Parties & Participants", order: 5 }
  ],
  config: {
    version: "1.0.0",
    lastUpdated: "2025-06-03T12:00:00Z",
    defaultCategory: "case_info"
  }
};

const testFilterEndpoints: FilterEndpoints = {
  configuration: "/api/filters/configuration",
  convertFilters: "/api/filters/convert",
  validateFilters: "/api/filters/validate",
  getOptions: "/api/filters/options"
};

// =============================================================================
// ENHANCED MOCK FETCH WITH BETTER ERROR HANDLING
// =============================================================================

interface MockFetchState {
  isSetup: boolean;
  requestCount: number;
  errors: string[];
}

const mockFetchState: MockFetchState = {
  isSetup: false,
  requestCount: 0,
  errors: []
};

function setupMockFetch(): void {
  if (mockFetchState.isSetup) {
    console.log('Mock fetch already set up');
    return;
  }

  // Store original fetch for fallback
  const originalFetch = window.fetch;

  window.fetch = async (input: RequestInfo | URL, init?: RequestInit): Promise<Response> => {
    const url = typeof input === 'string' ? input : input.toString();
    mockFetchState.requestCount++;

    console.log(`Mock fetch request #${mockFetchState.requestCount}:`, url);

    try {
      if (url.includes('/api/filters/configuration')) {
        console.log('Returning mock filter configuration');
        // Simulate realistic network delay
        await new Promise(resolve => setTimeout(resolve, 300 + Math.random() * 200));

        return new Response(JSON.stringify(testFilterConfiguration), {
          status: 200,
          headers: {
            'Content-Type': 'application/json',
            'X-Mock-Response': 'true'
          }
        });
      }

      if (url.includes('/api/filters/convert')) {
        console.log('Returning mock filter conversion');
        const body = init?.body ? JSON.parse(init.body as string) : {};
        await new Promise(resolve => setTimeout(resolve, 100));

        return new Response(JSON.stringify({
          converted: body.filters || {},
          timestamp: new Date().toISOString()
        }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' }
        });
      }

      if (url.includes('/api/filters/validate')) {
        console.log('Returning mock filter validation');
        await new Promise(resolve => setTimeout(resolve, 50));

        return new Response(JSON.stringify({
          isValid: true,
          errors: [],
          warnings: []
        }), {
          status: 200,
          headers: { 'Content-Type': 'application/json' }
        });
      }

      // Fall back to original fetch for other URLs
      console.log('Falling back to original fetch for:', url);
      return originalFetch(input, init);

    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Unknown fetch error';
      mockFetchState.errors.push(errorMessage);
      console.error('Mock fetch error:', error);

      return new Response(JSON.stringify({
        error: errorMessage,
        timestamp: new Date().toISOString()
      }), {
        status: 500,
        headers: { 'Content-Type': 'application/json' }
      });
    }
  };

  mockFetchState.isSetup = true;
  console.log('Mock fetch setup completed');
}

function getMockFetchStats() {
  return {
    ...mockFetchState,
    isWorking: mockFetchState.isSetup && mockFetchState.requestCount > 0
  };
}

function getAllTestFieldIds(): string[] {
  return testFilterConfiguration.fields.map(field => field.id);
}

// =============================================================================
// SYSTEM DIAGNOSTICS COMPONENT
// =============================================================================

function SystemDiagnostics() {
  const [diagnostics, setDiagnostics] = useState<any>({});
  const [isRunning, setIsRunning] = useState(false);

  const runDiagnostics = async () => {
    setIsRunning(true);
    try {
      const mockStats = getMockFetchStats();
      const systemStatus = getFilterSystemStatus();

      // Test filter manager creation
      let managerTest = { success: false, error: null };
      try {
        const testManager = createFilterManager(testFilterEndpoints);
        managerTest.success = true;
      } catch (error) {
        managerTest.error = error instanceof Error ? error.message : 'Unknown error';
      }

      // Test mock fetch
      let fetchTest = { success: false, error: null, response: null };
      try {
        const response = await fetch('/api/filters/configuration');
        fetchTest.success = response.ok;
        fetchTest.response = await response.json();
      } catch (error) {
        fetchTest.error = error instanceof Error ? error.message : null;
      }

      setDiagnostics({
        timestamp: new Date().toISOString(),
        mockFetch: mockStats,
        systemStatus,
        managerTest,
        fetchTest,
        environment: {
          isSSR: typeof window === 'undefined',
          hasLocalStorage: typeof window !== 'undefined' && !!window.localStorage,
          userAgent: typeof window !== 'undefined' ? navigator.userAgent : 'SSR',
        }
      });
    } catch (error) {
      console.error('Diagnostics error:', error);
    } finally {
      setIsRunning(false);
    }
  };

  useEffect(() => {
    runDiagnostics();
  }, []);

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h2 className="card-title">🔍 System Diagnostics</h2>

        <div className="flex gap-2 mb-4">
          <button
            onClick={runDiagnostics}
            className={`btn btn-primary btn-sm ${isRunning ? 'loading' : ''}`}
            disabled={isRunning}
          >
            {isRunning ? 'Running...' : '🔄 Run Diagnostics'}
          </button>
        </div>

        {diagnostics.timestamp && (
          <div className="space-y-4">
            {/* Mock Fetch Status */}
            <div className="alert alert-info">
              <div className="w-full">
                <h3 className="font-bold">Mock Fetch Status</h3>
                <div className="grid grid-cols-2 gap-2 text-sm mt-2">
                  <div>Setup: {diagnostics.mockFetch.isSetup ? '✅' : '❌'}</div>
                  <div>Working: {diagnostics.mockFetch.isWorking ? '✅' : '❌'}</div>
                  <div>Requests: {diagnostics.mockFetch.requestCount}</div>
                  <div>Errors: {diagnostics.mockFetch.errors.length}</div>
                </div>
              </div>
            </div>

            {/* Filter Manager Test */}
            <div className={`alert ${diagnostics.managerTest.success ? 'alert-success' : 'alert-error'}`}>
              <div className="w-full">
                <h3 className="font-bold">Filter Manager Creation</h3>
                <div className="text-sm mt-1">
                  {diagnostics.managerTest.success ?
                    '✅ Filter manager created successfully' :
                    `❌ Error: ${diagnostics.managerTest.error}`
                  }
                </div>
              </div>
            </div>

            {/* Fetch Test */}
            <div className={`alert ${diagnostics.fetchTest.success ? 'alert-success' : 'alert-error'}`}>
              <div className="w-full">
                <h3 className="font-bold">Configuration Fetch Test</h3>
                <div className="text-sm mt-1">
                  {diagnostics.fetchTest.success ?
                    '✅ Successfully fetched configuration' :
                    `❌ Fetch Error: ${diagnostics.fetchTest.error}`
                  }
                </div>
                {diagnostics.fetchTest.response && (
                  <details className="mt-2">
                    <summary className="cursor-pointer text-xs">View Response</summary>
                    <pre className="text-xs bg-base-200 p-2 rounded mt-1 overflow-auto">
                      {JSON.stringify(diagnostics.fetchTest.response, null, 2)}
                    </pre>
                  </details>
                )}
              </div>
            </div>

            {/* System Status */}
            <div className="alert alert-info">
              <div className="w-full">
                <h3 className="font-bold">System Status</h3>
                <div className="grid grid-cols-2 gap-2 text-sm mt-2">
                  <div>Ready: {diagnostics.systemStatus.isReady ? '✅' : '❌'}</div>
                  <div>Initialized: {diagnostics.systemStatus.isInitialized ? '✅' : '❌'}</div>
                  <div>Error: {diagnostics.systemStatus.hasError ? '❌' : '✅'}</div>
                  <div>Fields: {diagnostics.systemStatus.totalFields}</div>
                </div>
                {diagnostics.systemStatus.errorMessage && (
                  <div className="text-error text-sm mt-2">
                    Error: {diagnostics.systemStatus.errorMessage}
                  </div>
                )}
              </div>
            </div>

            {/* Environment */}
            <div className="collapse collapse-arrow bg-base-200">
              <input type="checkbox" />
              <div className="collapse-title text-sm font-medium">
                🌐 Environment Details
              </div>
              <div className="collapse-content text-xs">
                <pre>{JSON.stringify(diagnostics.environment, null, 2)}</pre>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

// =============================================================================
// WORKING EXAMPLE COMPONENT
// =============================================================================

function WorkingExample() {
  const [testResults, setTestResults] = useState<any>({});

  const testKesslerIntegration = async () => {
    const results: any = {
      timestamp: new Date().toISOString(),
      tests: {}
    };

    try {
      // Test 1: Basic filter system status
      results.tests.systemStatus = getFilterSystemStatus();

      // Test 2: Mock fetch test
      try {
        const response = await fetch('/api/filters/configuration');
        const data = await response.json();
        results.tests.fetchTest = {
          success: true,
          fieldsCount: data.fields?.length || 0,
          categoriesCount: data.categories?.length || 0
        };
      } catch (error) {
        results.tests.fetchTest = {
          success: false,
          error: error instanceof Error ? error.message : 'Unknown error'
        };
      }

      // Test 3: Filter manager creation
      try {
        const manager = createFilterManager(testFilterEndpoints);
        results.tests.managerCreation = {
          success: true,
          hasLoadConfiguration: typeof manager.loadConfiguration === 'function',
          hasGetField: typeof manager.getField === 'function'
        };
      } catch (error) {
        results.tests.managerCreation = {
          success: false,
          error: error instanceof Error ? error.message : 'Unknown error'
        };
      }

      setTestResults(results);
    } catch (error) {
      console.error('Test error:', error);
      setTestResults({
        error: error instanceof Error ? error.message : 'Unknown test error'
      });
    }
  };

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h2 className="card-title">✅ Working Example</h2>

        <div className="alert alert-success mb-4">
          <span><strong>This example should work:</strong> Basic integration test without complex components</span>
        </div>

        <button
          onClick={testKesslerIntegration}
          className="btn btn-primary mb-4"
        >
          🧪 Test Kessler Integration
        </button>

        {testResults.timestamp && (
          <div className="space-y-3">
            <h3 className="font-semibold">Test Results:</h3>

            {testResults.error ? (
              <div className="alert alert-error">
                <span>Test failed: {testResults.error}</span>
              </div>
            ) : (
              <div className="space-y-2">
                {/* System Status Test */}
                <div className={`alert ${testResults.tests.systemStatus.isReady ? 'alert-success' : 'alert-warning'}`}>
                  <div>
                    <div className="font-medium">System Status Test</div>
                    <div className="text-sm">
                      Ready: {testResults.tests.systemStatus.isReady ? '✅' : '❌'} •
                      Initialized: {testResults.tests.systemStatus.isInitialized ? '✅' : '❌'} •
                      Fields: {testResults.tests.systemStatus.totalFields}
                    </div>
                  </div>
                </div>

                {/* Fetch Test */}
                <div className={`alert ${testResults.tests.fetchTest.success ? 'alert-success' : 'alert-error'}`}>
                  <div>
                    <div className="font-medium">Fetch Test</div>
                    <div className="text-sm">
                      {testResults.tests.fetchTest.success ?
                        `✅ Fetched ${testResults.tests.fetchTest.fieldsCount} fields, ${testResults.tests.fetchTest.categoriesCount} categories` :
                        `❌ Error: ${testResults.tests.fetchTest.error}`
                      }
                    </div>
                  </div>
                </div>

                {/* Manager Creation Test */}
                <div className={`alert ${testResults.tests.managerCreation.success ? 'alert-success' : 'alert-error'}`}>
                  <div>
                    <div className="font-medium">Manager Creation Test</div>
                    <div className="text-sm">
                      {testResults.tests.managerCreation.success ?
                        `✅ Manager created with required methods` :
                        `❌ Error: ${testResults.tests.managerCreation.error}`
                      }
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Simple Filter Test */}
        <div className="mt-6 p-4 border rounded-lg">
          <h4 className="font-semibold mb-2">Simple Filter System Test</h4>
          <FilterErrorBoundary
            fallback={
              <div className="alert alert-error">
                <span>Filter component failed to render</span>
              </div>
            }
          >
            <div className="space-y-4">
              <KesslerFilterStatus />
              <KesslerDocumentFiltersList
                showFields={["case_number", "created_at"]} // Only show simple fields
                endpoints={testFilterEndpoints}
                enableUrlSync={false} // Disable to avoid URL issues
                enablePersistence={false} // Disable to avoid storage issues
                onFilterChange={(fieldId, value) => {
                  console.log(`✅ Filter changed: ${fieldId} = ${value}`);
                }}
              />
            </div>
          </FilterErrorBoundary>
        </div>
      </div>
    </div>
  );
}

// =============================================================================
// STEP BY STEP TUTORIAL
// =============================================================================

function StepByStepTutorial() {
  const [currentStep, setCurrentStep] = useState(0);
  const [stepResults, setStepResults] = useState<any[]>([]);

  const steps = [
    {
      title: "Setup Mock Fetch",
      description: "Initialize mock API responses",
      test: async () => {
        setupMockFetch();
        return { success: true, message: "Mock fetch initialized" };
      }
    },
    {
      title: "Test Configuration Fetch",
      description: "Verify we can fetch filter configuration",
      test: async () => {
        try {
          const response = await fetch('/api/filters/configuration');
          const data = await response.json();
          return {
            success: true,
            message: `Fetched configuration with ${data.fields.length} fields`
          };
        } catch (error) {
          return {
            success: false,
            message: error instanceof Error ? error.message : 'Fetch failed'
          };
        }
      }
    },
    {
      title: "Create Filter Manager",
      description: "Test filter manager creation",
      test: async () => {
        try {
          const manager = createFilterManager(testFilterEndpoints);
          return {
            success: true,
            message: "Filter manager created successfully"
          };
        } catch (error) {
          return {
            success: false,
            message: error instanceof Error ? error.message : 'Manager creation failed'
          };
        }
      }
    },
    {
      title: "Load Configuration",
      description: "Test loading configuration through manager",
      test: async () => {
        try {
          const manager = createFilterManager(testFilterEndpoints);
          const config = await manager.loadConfiguration();
          return {
            success: true,
            message: `Configuration loaded with ${config.fields.length} fields`
          };
        } catch (error) {
          return {
            success: false,
            message: error instanceof Error ? error.message : 'Configuration loading failed'
          };
        }
      }
    }
  ];

  const runStep = async (stepIndex: number) => {
    setCurrentStep(stepIndex);
    try {
      const result = await steps[stepIndex].test();
      setStepResults(prev => {
        const newResults = [...prev];
        newResults[stepIndex] = result;
        return newResults;
      });
    } catch (error) {
      setStepResults(prev => {
        const newResults = [...prev];
        newResults[stepIndex] = {
          success: false,
          message: error instanceof Error ? error.message : 'Step failed'
        };
        return newResults;
      });
    }
  };

  const runAllSteps = async () => {
    for (let i = 0; i < steps.length; i++) {
      await runStep(i);
    }
  };

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h2 className="card-title">📚 Step-by-Step Tutorial</h2>

        <div className="flex gap-2 mb-4">
          <button
            onClick={runAllSteps}
            className="btn btn-primary btn-sm"
          >
            ▶️ Run All Steps
          </button>
          <button
            onClick={() => setStepResults([])}
            className="btn btn-outline btn-sm"
          >
            🔄 Reset
          </button>
        </div>

        <div className="space-y-3">
          {steps.map((step, index) => (
            <div key={index} className="border rounded-lg p-3">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-semibold">
                    Step {index + 1}: {step.title}
                  </h3>
                  <p className="text-sm text-gray-600">{step.description}</p>
                </div>
                <div className="flex items-center gap-2">
                  {stepResults[index] && (
                    <span className={`badge ${stepResults[index].success ? 'badge-success' : 'badge-error'}`}>
                      {stepResults[index].success ? '✅' : '❌'}
                    </span>
                  )}
                  <button
                    onClick={() => runStep(index)}
                    className="btn btn-xs btn-outline"
                    disabled={currentStep === index}
                  >
                    {currentStep === index ? '...' : 'Run'}
                  </button>
                </div>
              </div>

              {stepResults[index] && (
                <div className={`mt-2 p-2 rounded text-sm ${stepResults[index].success ? 'bg-green-50 text-green-800' : 'bg-red-50 text-red-800'
                  }`}>
                  {stepResults[index].message}
                </div>
              )}
            </div>
          ))}
        </div>

        {stepResults.length === steps.length && (
          <div className="mt-4 alert alert-info">
            <span>
              Tutorial complete! {stepResults.filter(r => r.success).length} of {steps.length} steps passed.
            </span>
          </div>
        )}
      </div>
    </div>
  );
}

// =============================================================================
// MAIN TEST PAGE COMPONENT
// =============================================================================

type TestMode = 'diagnostics' | 'working' | 'tutorial' | 'layouts';

interface TestOption {
  id: TestMode;
  name: string;
  icon: string;
  description: string;
  component: React.ComponentType;
}

export default function KesslerTestPageFixed() {
  const [activeMode, setActiveMode] = useState<TestMode>('diagnostics');

  // Setup mock fetch on component mount
  useEffect(() => {
    console.log('Setting up mock fetch for test page...');
    setupMockFetch();
  }, []);

  const testOptions: TestOption[] = [
    {
      id: 'diagnostics',
      name: 'System Diagnostics',
      icon: '🔍',
      description: 'Debug system issues and check configuration',
      component: SystemDiagnostics
    },
    {
      id: 'working',
      name: 'Working Example',
      icon: '✅',
      description: 'Simple integration test that should work',
      component: WorkingExample
    },
    {
      id: 'tutorial',
      name: 'Step-by-Step Tutorial',
      icon: '📚',
      description: 'Guided setup and testing process',
      component: StepByStepTutorial
    },
    {
      id: 'layouts',
      name: 'Layout Tests',
      icon: '📋',
      description: 'Test different layout configurations (advanced)',
      component: () => (
        <div className="space-y-6">
          <div className="alert alert-warning">
            <span><strong>Advanced:</strong> Run diagnostics first to ensure basic functionality works</span>
          </div>

          <FilterErrorBoundary>
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title">List Layout Test</h3>
                <KesslerDocumentFiltersList
                  showFields={getAllTestFieldIds()}
                  endpoints={testFilterEndpoints}
                  enableUrlSync={false}
                  enablePersistence={false}
                  onFilterChange={(fieldId, value) => {
                    console.log(`📋 List Layout - ${fieldId}:`, value);
                  }}
                />
              </div>
            </div>
          </FilterErrorBoundary>

          <FilterErrorBoundary>
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title">Responsive Layout Test</h3>
                <KesslerResponsiveDynamicDocumentFilters
                  showFields={getAllTestFieldIds()}
                  endpoints={testFilterEndpoints}
                  enableUrlSync={false}
                  enablePersistence={false}
                  onFilterChange={(fieldId, value) => {
                    console.log(`📱 Responsive Layout - ${fieldId}:`, value);
                  }}
                />
              </div>
            </div>
          </FilterErrorBoundary>

          <FilterErrorBoundary>
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title">Inline Layout Test</h3>
                <KesslerInlineDynamicDocumentFilters
                  showFields={["case_number", "created_at", "filing_type"]}
                  endpoints={testFilterEndpoints}
                  enableUrlSync={false}
                  enablePersistence={false}
                  onFilterChange={(fieldId, value) => {
                    console.log(`➡️ Inline Layout - ${fieldId}:`, value);
                  }}
                />
              </div>
            </div>
          </FilterErrorBoundary>
        </div>
      )
    }
  ];

  const activeOption = testOptions.find(option => option.id === activeMode);
  const ActiveComponent = activeOption?.component || SystemDiagnostics;

  return (
    <div className="container mx-auto p-6 space-y-6 max-w-7xl">
      {/* Header */}
      <div className="text-center">
        <h1 className="text-4xl font-bold mb-2">🔧 Kessler Integration Debugger</h1>
        <p className="text-gray-600">Diagnose and fix filter system initialization issues</p>
      </div>

      {/* Quick Status */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title">🚨 Quick Status Check</h2>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="stat bg-base-200 rounded-lg">
              <div className="stat-title">Mock Fetch</div>
              <div className={`stat-value text-lg ${getMockFetchStats().isSetup ? 'text-success' : 'text-error'}`}>
                {getMockFetchStats().isSetup ? '✅ Ready' : '❌ Not Setup'}
              </div>
              <div className="stat-desc">
                Requests: {getMockFetchStats().requestCount}
              </div>
            </div>

            <div className="stat bg-base-200 rounded-lg">
              <div className="stat-title">Filter System</div>
              <div className={`stat-value text-lg ${getFilterSystemStatus().isReady ? 'text-success' : 'text-warning'}`}>
                {getFilterSystemStatus().isReady ? '✅ Ready' : '⏳ Not Ready'}
              </div>
              <div className="stat-desc">
                {getFilterSystemStatus().hasError ? 'Has Errors' : 'No Errors'}
              </div>
            </div>

            <div className="stat bg-base-200 rounded-lg">
              <div className="stat-title">Configuration</div>
              <div className={`stat-value text-lg ${testFilterConfiguration.fields.length > 0 ? 'text-success' : 'text-error'}`}>
                {testFilterConfiguration.fields.length} Fields
              </div>
              <div className="stat-desc">
                {testFilterConfiguration.categories.length} Categories
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Mode Selector */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title mb-4">Select Debug Mode</h2>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-4">
            {testOptions.map(option => (
              <button
                key={option.id}
                className={`btn btn-outline h-auto flex-col p-4 ${activeMode === option.id ? 'btn-active' : ''
                  }`}
                onClick={() => setActiveMode(option.id)}
              >
                <span className="text-2xl mb-1">{option.icon}</span>
                <span className="text-sm font-semibold">{option.name}</span>
              </button>
            ))}
          </div>

          <div className="p-4 bg-base-200 rounded-lg">
            <div className="flex items-center gap-3">
              <span className="text-3xl">{activeOption?.icon}</span>
              <div>
                <h3 className="font-bold text-lg">{activeOption?.name}</h3>
                <p className="text-sm text-gray-600">{activeOption?.description}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Active Debug Component */}
      <ActiveComponent />

      {/* Common Issues & Solutions */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title mb-4">🛠️ Common Issues & Solutions</h2>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h3 className="font-semibold text-red-700 mb-2">❌ Common Errors</h3>
              <ul className="text-sm space-y-2">
                <li className="flex items-start gap-2">
                  <span className="text-red-500">•</span>
                  <div>
                    <strong>"Failed to initialize filter system"</strong>
                    <br />
                    <span className="text-gray-600">Usually caused by missing mock fetch or invalid endpoints</span>
                  </div>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-red-500">•</span>
                  <div>
                    <strong>"syncWithUrl is not a function"</strong>
                    <br />
                    <span className="text-gray-600">Store hooks not properly exported or initialized</span>
                  </div>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-red-500">•</span>
                  <div>
                    <strong>"FilterManager not initialized"</strong>
                    <br />
                    <span className="text-gray-600">Filter manager creation failed or async timing issue</span>
                  </div>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-red-500">•</span>
                  <div>
                    <strong>Components not rendering</strong>
                    <br />
                    <span className="text-gray-600">Error boundaries catching initialization failures</span>
                  </div>
                </li>
              </ul>
            </div>

            <div>
              <h3 className="font-semibold text-green-700 mb-2">✅ Solutions</h3>
              <ul className="text-sm space-y-2">
                <li className="flex items-start gap-2">
                  <span className="text-green-500">•</span>
                  <div>
                    <strong>Run System Diagnostics first</strong>
                    <br />
                    <span className="text-gray-600">Check if mock fetch is working and configuration loads</span>
                  </div>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-green-500">•</span>
                  <div>
                    <strong>Use Working Example tab</strong>
                    <br />
                    <span className="text-gray-600">Test basic integration without advanced features</span>
                  </div>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-green-500">•</span>
                  <div>
                    <strong>Disable URL sync and persistence initially</strong>
                    <br />
                    <span className="text-gray-600">Reduce complexity while debugging core issues</span>
                  </div>
                </li>
                <li className="flex items-start gap-2">
                  <span className="text-green-500">•</span>
                  <div>
                    <strong>Check browser console</strong>
                    <br />
                    <span className="text-gray-600">Look for detailed error messages and network requests</span>
                  </div>
                </li>
              </ul>
            </div>
          </div>

          <div className="mt-6 p-4 bg-amber-50 border border-amber-200 rounded-lg">
            <h4 className="font-semibold text-amber-800 mb-2">💡 Debug Process</h4>
            <ol className="text-sm text-amber-700 list-decimal list-inside space-y-1">
              <li>Start with <strong>System Diagnostics</strong> to identify core issues</li>
              <li>Use <strong>Working Example</strong> to test basic functionality</li>
              <li>Follow <strong>Step-by-Step Tutorial</strong> for guided troubleshooting</li>
              <li>Only test <strong>Layout Tests</strong> after core functionality works</li>
              <li>Check browser DevTools console for detailed error messages</li>
              <li>Verify mock fetch is working by checking Network tab</li>
            </ol>
          </div>
        </div>
      </div>

      {/* Debug Information */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title mb-4">🔍 Current Debug Information</h2>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <h3 className="font-semibold mb-2">Environment</h3>
              <div className="text-sm space-y-1">
                {/* <div>localStorage: {typeof window !== 'undefined' && window.localStorage ? 'Available' : 'Not available'}</div> */}
                {/* <div>URL: {typeof window !== 'undefined' ? window.location.pathname : 'SSR'}</div> */}
              </div>
            </div>

            <div>
              <h3 className="font-semibold mb-2">Mock Fetch</h3>
              <div className="text-sm space-y-1">
                <div>Setup: {getMockFetchStats().isSetup ? 'Yes' : 'No'}</div>
                <div>Requests: {getMockFetchStats().requestCount}</div>
                <div>Errors: {getMockFetchStats().errors.length}</div>
              </div>
            </div>

            <div>
              <h3 className="font-semibold mb-2">Configuration</h3>
              <div className="text-sm space-y-1">
                <div>Fields: {testFilterConfiguration.fields.length}</div>
                <div>Categories: {testFilterConfiguration.categories.length}</div>
                <div>Version: {testFilterConfiguration.config.version}</div>
              </div>
            </div>
          </div>

          <div className="mt-4">
            <details className="collapse collapse-arrow bg-base-200">
              <summary className="collapse-title text-sm font-medium">
                📊 Detailed Filter Configuration
              </summary>
              <div className="collapse-content text-xs">
                <pre className="overflow-auto max-h-64">
                  {JSON.stringify(testFilterConfiguration, null, 2)}
                </pre>
              </div>
            </details>
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className="text-center p-4 text-gray-500">
        <p>🔧 Kessler Integration Debugger • Built for troubleshooting filter system issues</p>
        <p className="text-sm mt-1">
          Use this tool to identify and fix initialization problems before using the main filter system
        </p>
      </div>
    </div>
  );
}