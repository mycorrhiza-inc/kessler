import { z } from "zod";

export const PGStageValidator = z.enum([
  "pending",
  "processing",
  "completed",
  "errored",
]);

export type PGStage = z.infer<typeof PGStageValidator>;

export const DocProcStatusValidator = z.enum([
  "unprocessed",
  "completed",
  "encounters_analyzed",
  "organization_assigned",
  "summarization_completed",
  "embeddings_completed",
  "upload_document_to_db",
  "stage3",
  "stage2",
  "stage1",
]);

export type DocProcStage = z.infer<typeof DocProcStatusValidator>;

export const FileChildTextSourceValidator = z.object({
  is_original_text: z.boolean(),
  text: z.string(),
  language: z.string(),
});

export const DocProcStageValidator = z.object({
  pg_stage: PGStageValidator,
  docproc_stage: DocProcStatusValidator,
  is_errored: z.boolean(),
  is_completed: z.boolean(),
  processing_error_msg: z.string(),
  database_error_msg: z.string(),
});

export const FileGeneratedExtrasValidator = z.object({
  summary: z.string(),
  short_summary: z.string(),
  purpose: z.string(),
  impressiveness: z.number(),
});

export const AuthorInformationValidator = z.object({
  author_name: z.string(),
  is_person: z.boolean(),
  is_primary_author: z.boolean(),
  author_id: z.string(),
});

export const JuristictionInformationValidator = z.object({
  country: z.string(),
  state: z.string(),
  municipality: z.string(),
  agency: z.string(),
  proceeding_name: z.string(),
  extra_object: z.record(z.any()),
});

export const CompleteFileSchemaValidator = z.object({
  id: z.string(),
  verified: z.boolean(),
  extension: z.string(),
  lang: z.string(),
  name: z.string(),
  hash: z.string(),
  is_private: z.boolean(),
  mdata: z.record(z.any()),
  stage: DocProcStageValidator,
  extra: FileGeneratedExtrasValidator,
  authors: z.array(AuthorInformationValidator),
  juristiction: JuristictionInformationValidator,
  doc_texts: z.array(FileChildTextSourceValidator),
});

export type CompleteFileSchema = z.infer<typeof CompleteFileSchemaValidator>;
