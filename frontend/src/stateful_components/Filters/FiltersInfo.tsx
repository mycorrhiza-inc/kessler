import { FilterField } from "@/lib/filters";

export enum InputType {
  Text = "text",
  Select = "select",
  Datalist = "datalist",
  Date = "date",
  NotShown = "not_shown",
  OrgMultiselect = "org_multiselect",
  ConvoMultiselect = "convo_multiselect",
  Organization = "organization",
  Docket = "docket",
  FileClass = "file_class",
  FileExtension = "extension",
  NYDocket = "nypuc_docket_industry",
}
export type PropertyInformation = {
  type: InputType;
  index: number;
  displayName: string;
  description: string;
  details: string;
  options?: { label: string; value: string }[];
  isDate?: boolean;
};
export type QueryFiltersInformation = {
  [key in FilterField]: PropertyInformation;
};
export const queryFiltersInformation: QueryFiltersInformation = {
  match_name: {
    type: InputType.Text,
    index: 1,
    displayName: "Title",
    description: "Filter on the title of the document.",
    details: "Searches for items approximately matching the title",
  },
  match_docket_id: {
    type: InputType.ConvoMultiselect,
    index: 2,
    displayName: "Docket ID",
    description: "The unique identifier for the docket.",
    details: "Filters search results based on the docket ID.",
  },
  match_author: {
    type: InputType.OrgMultiselect,
    index: 3,
    displayName: "Author",
    description: "The author of the document.",
    details: "Searches for items created or written by the specified author.",
  },
  match_source: {
    type: InputType.Select,
    index: 4,
    displayName: "State",
    description: "What state do you want to limit your search to?",
    details: "",
    options: [
      { label: "All", value: "" },
      { label: "New York", value: "usa-ny" },
      { label: "Colorado", value: "usa-co" },
      { label: "California", value: "usa-ca" },
      { label: "National", value: "usa-national" },
    ],
  },
  match_extension: {
    type: InputType.Select,
    index: 5,
    displayName: "Document File Type",
    description: "The type or category of the document.",
    details: "Searches for items that match the specified document type.",
    options: [
      { label: "All", value: "" },
      { label: "PDF", value: "pdf" },
      { label: "Docx", value: "docx" },
      { label: "Xls", value: "xls" },
    ],
  },
  match_file_class: {
    type: InputType.Select,
    index: 6,
    displayName: "Document Class",
    description: "The classification or category of the document.",
    details: "Searches for documents that fall under the specified class.",
    options: [
      { label: "All", value: "" },
      { label: "Correspondence", value: "Correspondence" },
      { label: "Comments", value: "Comments" },
      { label: "Reports", value: "Reports" },
      { label: "Plans and Proposals", value: "Plans and Proposals" },
      { label: "Motions", value: "Motions" },
      { label: "Letters", value: "Letters" },
      { label: "Orders", value: "Orders" },
      { label: "Notices", value: "Notices" },
    ],
  },
  match_before_date: {
    type: InputType.Date,
    index: 7,
    displayName: "Before Date",
    description: "The date related to the document.",
    details: "Filters results by the specified date.",
  },
  match_after_date: {
    type: InputType.Date,
    index: 8,
    displayName: "After Date",
    description: "The date related to the document.",
    details: "Filters results by the specified date.",
  },
  match_date_range: {
    type: InputType.Date,
    index: 8,
    displayName: "Within Date Range",
    description: "The range of dates related to the document.",
    details: "Filters results by the specified date range.",
  },
  match_file_uuid: {
    type: InputType.NotShown,
    index: -1,
    displayName: "",
    description: "",
    details: "",
  },
  match_author_uuids: {
    type: InputType.NotShown,
    index: -1,
    displayName: "",
    description: "",
    details: "",
  },
  match_conversation_uuid: {
    type: InputType.NotShown,
    index: -1,
    displayName: "",
    description: "",
    details: "",
  },
};
