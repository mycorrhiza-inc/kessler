"use client";

import Link from 'next/link';
import React, { memo, useState } from "react";

import PDFViewer from "./PDFViewer";
import XlsxViewer from "@/components/style/messages/XlsxCannotBeViewedMessage";
import ErrorMessage from "@/components/style/messages/ErrorMessage";
import MarkdownRenderer from "@/components/style/misc/MarkdownRenderer";
import LoadingSpinner from "@/components/style/misc/LoadingSpinner";
import { FileExtension, fileExtensionFromText } from "@/components/style/Pills/FileExtension";
import { CLIENT_API_URL } from "@/lib/env_variables";
import { CompleteFileSchema } from "@/lib/types/backend_schemas";
import { AuthorInformation } from "@/lib/types/backend_schemas";
import { DocketPill, AuthorInfoPill } from "@/components/style/Pills/TextPills";

// Minimal data shape required by DocumentMainTabs
export interface DocumentMainTabsData {
  id: CompleteFileSchema['id'];
  mdata: CompleteFileSchema['mdata'];
  hash: CompleteFileSchema['hash'];
  verified: CompleteFileSchema['verified'];
  extension: CompleteFileSchema['extension'];
}

const MarkdownContent = memo(({ docUUID }: { docUUID: string }) => {
  // TODO: wire up an SSR-friendly fetch for text
  return <p>Text view coming soon…</p>;
});

const MetadataContent = memo(({ metadata }: { metadata: Record<string, any> }) => (
  <div className="overflow-x-auto">
    <table className="table table-zebra w-full">
      <thead>
        <tr>
          <th>Field</th>
          <th>Value</th>
        </tr>
      </thead>
      <tbody>
        {Object.entries(metadata).map(([key, value]) => (
          <tr key={key}>
            <td>{key}</td>
            <td>{String(value)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  </div>
));

const DocumentContent = ({
  docUUID,
  extension,
}: {
  docUUID: string;
  extension: FileExtension;
}) => {
  const documentUrl = `${CLIENT_API_URL}/public/files/${docUUID}/raw`;
  if (extension === FileExtension.PDF) return <PDFViewer file={documentUrl} />;
  if (extension === FileExtension.XLSX) return <XlsxViewer file={documentUrl} />;
  return <ErrorMessage error="Cannot display this document type" />;
};

const DocumentHeader = ({
  documentObject,
  isPage,
}: {
  documentObject: CompleteFileSchema;
  isPage: boolean;
}) => {
  const { name: title, id: objectId, extension = 'pdf', verified = false, extra, mdata, authors } = documentObject;
  const summary = extra.summary;
  const underscoredTitle = title.replace(/\s+/g, '_') || 'Unknown_Document';
  const downloadUrl = `${CLIENT_API_URL}/public/files/${objectId}/raw/${underscoredTitle}.${extension}`;
  const viewUrl = `/files/${objectId}`;
  return (
    <>
      <div className="card-title flex justify-between items-start">
        <h1 className="text-3xl break-words max-w-[70%]">{title}</h1>
        <div className="flex gap-2">
          <a className="btn btn-primary" href={downloadUrl} target="_blank" download={title}>
            Download File
          </a>
          {!isPage && (
            <Link className="btn btn-secondary" href={viewUrl} target="_blank">
              Open in New Tab
            </Link>
          )}
        </div>
      </div>
      <p><b>Case Number:</b> <DocketPill docketId={mdata.docket_id as string} /></p>
      {authors && (
        <p>
          <b>{authors.length === 1 ? 'Author' : 'Authors'}:</b>{' '}
          {authors.map((a: AuthorInformation) => (
            <AuthorInfoPill author_info={a} key={a.author_id} />
          ))}
        </p>
      )}
      <div className="p-4" />
      <h2 className="text-xl"><b>LLM Summary:</b></h2>
      <MarkdownRenderer>
        {verified
          ? summary
          : 'Document processing in progress. Please check back later.'}
      </MarkdownRenderer>
      <div className="p-12" />
    </>
  );
};

export const DocumentMainTabsClient = ({
  documentObject,
  isPage,
}: {
  documentObject: DocumentMainTabsData;
  isPage: boolean;
}) => {
  const { id: docUUID, mdata, hash, verified, extension: rawExt } = documentObject;
  const metadata = { ...mdata, hash };
  const fileExt = fileExtensionFromText(rawExt);
  const showRaw = verified;
  const showText = verified && rawExt !== 'xlsx';

  type TabKey = 'raw' | 'text' | 'meta';
  const defaultTab: TabKey = showRaw ? 'raw' : showText ? 'text' : 'meta';
  const [activeTab, setActiveTab] = useState<TabKey>(defaultTab);

  return (
    <div className="modal-content standard-box">

      <div role="tablist" aria-label="Document Sections" className="flex gap-2 border-b border-base-300">
        {showRaw && (
          <button
            role="tab"
            aria-selected={activeTab === 'raw'}
            onClick={() => setActiveTab('raw')}
            className={`px-6 py-3 font-bold ${activeTab === 'raw' ? 'border-b-2 border-primary text-primary' : 'hover:bg-base-200'}`}
          >
            📄 Document
          </button>
        )}
        {showText && (
          <button
            role="tab"
            aria-selected={activeTab === 'text'}
            onClick={() => setActiveTab('text')}
            className={`px-6 py-3 font-bold ${activeTab === 'text' ? 'border-b-2 border-primary text-primary' : 'hover:bg-base-200'}`}
          >
            📝 Text View
          </button>
        )}
        <button
          role="tab"
          aria-selected={activeTab === 'meta'}
          onClick={() => setActiveTab('meta')}
          className={`px-6 py-3 font-bold ${activeTab === 'meta' ? 'border-b-2 border-primary text-primary' : 'hover:bg-base-200'}`}
        >
          ℹ️ Metadata
        </button>
      </div>

      <div className="mt-4">
        {showRaw && activeTab === 'raw' && <DocumentContent docUUID={docUUID} extension={fileExt} />}
        {showText && activeTab === 'text' && <MarkdownContent docUUID={docUUID} />}
        {activeTab === 'meta' && <MetadataContent metadata={metadata} />}
      </div>
    </div>
  );
};
