import { string } from "zod";

export type Conversation = {
  id: string; // UUID
  docket_id: string;
  name: string;
  description: string;
  updated_at: string;
};

export type ExtendedNYConversation = {
  id: string;
  matter_number: string;
  docket_id: string;
  industry_affected: string;
  organization: string;
  matter_type: string;
  matter_subtype: string;
  description: string; // title of matter
  related_matters: string[]; // Related Matter/Case No. array of docket_ids
  assigned_judge: string;
  updated_at: string;
};

export type NYConversation = {
  docket_id: string;
  matter_type: string;
  matter_subtype: string;
  title: string;
  organization: string;
  date_filed: string;
};

