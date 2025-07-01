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
import { ConversationPill, AuthorInfoPill } from "@/components/style/Pills/TextPills";

// Minimal data shape required by DocumentMainTabs

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

// Legacy Data
export interface DocumentMainTabsData {
  id: string;
  mdata: any;
  hash: string;
  verified: boolean;
  extension: string;
  name: string; // Added for attachment names
}

// New Data
export interface DocumentTabsData {
  attachments: {
    attachment_uuid: string;
    attachment_hash: string;
    attachment_name: string;
    attachment_extension: string;
    // Note: we are adding verified to Data, because we want to validate text view per atttachment
    attachment_verified: boolean; 
  }[]
}

export type AttachmentItem = {
  id: string;
  name: string;
  extension: string;
  verified: boolean;
}

// All of these tabs should be instantiated for every single attachment, it should end up requiring 2 layers of tabs 
// - Attachment 1 
//   - Attachment 1 Preview 
//   - Attachment 1 Text
// - Attachment 2 
//   - Attachment 2 Preview 
//   - Attachment 2 Text
export const DocumentMainTabsClient = ({
  documentObject,
  attachments,
  isPage,
}: {
  documentObject: DocumentMainTabsData;
  attachments: DocumentTabsData;
  isPage: boolean;
}) => {
  // Precompute the attachments
  const parsedAttachments: Array<AttachmentItem> = [{
    id: documentObject.id,
    name: documentObject.name,
    extension: documentObject.extension,
    verified: documentObject.verified
  }, ...attachments.attachments.map(a => ({
    id: a.attachment_uuid,
    name: a.attachment_name,
    extension: a.attachment_extension,
    verified: a.attachment_verified
  }))];

  const [activeAttachmentIndex, setActiveAttachmentIndex] = useState<number>(0);
  type TabKey = 'raw' | 'text';
  const [activeTabPerAttachment, setActiveTabPerAttachment] = useState<{[key: number]: TabKey}>(() => {
    let result: {[key: number]: TabKey} = {};
    parsedAttachments.forEach((attachment, i) => {
      result[i] = 'raw'; // default tab per attachment
    });
    return result;
  });

  const activeAttachment = parsedAttachments[activeAttachmentIndex];
  // For each attachment, we check showRaw and showText
  const showRaw = true; // Always show raw tab
  const showText = activeAttachment.verified && activeAttachment.extension !== 'xlsx';

  // Called when changing attachment
  const changeActiveAttachment = (index: number) => {
    setActiveAttachmentIndex(index);
    // Keep the current active tab for this attachment, if not set then default to 'raw'
    if (activeTabPerAttachment[index] === undefined) {
      setActiveTabPerAttachment(prev => ({
        ...prev,
        [index]: 'raw'
      }));
    }
  };

  // Called when changing tab for a given attachment
  const changeActiveTab = (tab: TabKey, index: number) => {
    setActiveTabPerAttachment(prev => ({
      ...prev,
      [index]: tab
    }));
  }

  return (
    <div className="modal-content standard-box">
      {/* Top-level tabs for each attachment */}
      <div role="tablist" aria-label="Attachment Sections" className="flex gap-2 border-b border-base-300 overflow-x-auto">
        {parsedAttachments.map((attachment, i) => (
          <button
            key={`${attachment.id}_top`}
            role="tab"
            aria-selected={activeAttachmentIndex === i}
            onClick={() => changeActiveAttachment(i)}
            className={`px-6 py-3 font-bold flex-shrink-0 ${activeAttachmentIndex === i ? 'border-b-2 border-primary text-primary' : 'hover:bg-base-200'}`}
          >
            {attachment.name}
          </button>
        ))}
      </div>

      {/* Sub-tabs for the active attachment */}
      {parsedAttachments.length > 0 && (
        <>
          <div role="tablist" aria-label="Document Sections" className="flex gap-2 border-b border-base-300">
            {showRaw && (
              <button
                role="tab"
                aria-selected={activeTabPerAttachment[activeAttachmentIndex] === 'raw'}
                onClick={() => changeActiveTab('raw', activeAttachmentIndex)}
                className={`px-6 py-3 font-bold ${activeTabPerAttachment[activeAttachmentIndex] === 'raw' ? 'border-b-2 border-primary text-primary' : 'hover:bg-base-200'}`}
              >
                📄 Document
              </button>
            )}
            {showText && (
              <button
                role="tab"
                aria-selected={activeTabPerAttachment[activeAttachmentIndex] === 'text'}
                onClick={() => changeActiveTab('text', activeAttachmentIndex)}
                className={`px-6 py-3 font-bold ${activeTabPerAttachment[activeAttachmentIndex] === 'text' ? 'border-b-2 border-primary text-primary' : 'hover:bg-base-200'}`}
              >
                📝 Text View
              </button>
            )}
          </div>

          {/* Content for the active attachment */}
          <div className="mt-4">
            {showRaw && activeTabPerAttachment[activeAttachmentIndex] === 'raw' && 
              <DocumentContent 
                docUUID={activeAttachment.id} 
                extension={fileExtensionFromText(activeAttachment.extension)} 
              />
            }
            {showText && activeTabPerAttachment[activeAttachmentIndex] === 'text' && 
              <MarkdownContent docUUID={activeAttachment.id} />
            }
          </div>
        </>
      )}
    </div>
  );
};
