// Proxy module for Document Filters

// Re-exporting from DynamicFilters
import * as DynamicFilters from "./DynamicFilters";

export const DynamicDocumentFilters = DynamicFilters.default;
export type DynamicDocumentFiltersProps = DynamicFilters.KesslerDynamicFiltersProps;
export const DynamicDocumentFiltersList = DynamicFilters.KesslerDocumentFiltersList;
export const DynamicDocumentFiltersGrid = DynamicFilters.KesslerDocumentFiltersGrid;
export const ResponsiveDynamicDocumentFilters = DynamicFilters.KesslerResponsiveDynamicDocumentFilters;
export const InlineDynamicDocumentFilters = DynamicFilters.KesslerInlineDynamicDocumentFilters;

export default DynamicDocumentFilters;
