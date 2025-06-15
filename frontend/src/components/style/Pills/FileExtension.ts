export enum FileExtension {
  PDF,
  XLSX,
  DOCX,
  HTML,
  UNKNOWN,
}

export const fileExtensionFromText = (text?: string): FileExtension => {
  if (!text) return FileExtension.UNKNOWN;
  const lowerText = text.toLowerCase();

  if (lowerText.endsWith("pdf")) return FileExtension.PDF;
  if (lowerText.endsWith("xlsx")) return FileExtension.XLSX;
  if (lowerText.endsWith("docx")) return FileExtension.DOCX;
  if (lowerText.endsWith("doc")) return FileExtension.DOCX;
  if (lowerText.endsWith("html")) return FileExtension.HTML;

  return FileExtension.UNKNOWN;
};
