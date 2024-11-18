export type PGStage = "pending" | "processing" | "completed" | "errored";

export type DocProcStatus =
  | "unprocessed"
  | "completed"
  | "encounters_analyzed"
  | "organization_assigned"
  | "summarization_completed"
  | "embeddings_completed"
  | "upload_document_to_db"
  | "stage3"
  | "stage2"
  | "stage1";

export interface FileChildTextSource {
  is_original_text: boolean;
  text: string;
  language: string;
}

export interface DocProcStage {
  pg_stage: PGStage;
  docproc_stage: DocProcStatus;
  is_errored: boolean;
  is_completed: boolean;
  processing_error_msg: string;
  database_error_msg: string;
}

export interface FileGeneratedExtras {
  summary: string;
  short_summary: string;
  purpose: string;
  impressiveness: number;
}

export interface AuthorInformation {
  author_name: string;
  is_person: boolean;
  is_primary_author: boolean;
  author_id: string;
}

export interface JuristictionInformation {
  country: string;
  state: string;
  municipality: string;
  agency: string;
  proceeding_name: string;
  extra_object: Record<string, any>;
}

export interface CompleteFileSchema {
  id: string;
  verified: boolean;
  extension: string;
  lang: string;
  name: string;
  hash: string;
  is_private: boolean;
  mdata: FileMetadataSchema;
  stage: DocProcStage;
  extra: FileGeneratedExtras;
  authors: AuthorInformation[];
  juristiction: JuristictionInformation;
  doc_texts: FileChildTextSource[];
}

export type FileMetadataSchema = Record<string, any>;
