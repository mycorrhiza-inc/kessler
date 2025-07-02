"use client";
import React, { memo, useState } from "react";

import PDFViewer from "./PDFViewer";
import XlsxViewer from "@/components/style/messages/XlsxCannotBeViewedMessage";
import ErrorMessage from "@/components/style/messages/ErrorMessage";
import { FileExtension, fileExtensionFromText } from "@/components/style/Pills/FileExtension";
import { CLIENT_API_URL } from "@/lib/env_variables";
import { FilePageInfo, FillingAttachmentInfo } from '../RenderedObjectCards/RednderedObjectCard';
import clsx from "clsx";
import { color, subdividedHueFromSeed, subdividedHueRaw } from "@/components/style/Pills/TextPills";

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
  attachmentInfos,
  isPage,
}: {
  attachmentInfos: Array<FillingAttachmentInfo>;
  isPage: boolean;
}) => {
  // Precompute the attachments

  const [activeAttachmentIndex, setActiveAttachmentIndex] = useState<number>(0);
  type TabKey = 'raw' | 'text';
  const [activeTabPerAttachment, setActiveTabPerAttachment] = useState<{ [key: number]: TabKey }>(() => {
    let result: { [key: number]: TabKey } = {};
    attachmentInfos.forEach((attachment, i) => {
      result[i] = 'raw'; // default tab per attachment
    });
    return result;
  });

  const activeAttachment = attachmentInfos[activeAttachmentIndex];

  // For each attachment, we check showRaw and showText
  const showRaw = true; // Always show raw tab
  const showText = true && activeAttachment.attachment_extension !== 'xlsx';

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
  const attachmentsHue = attachmentInfos.map((val, index) => subdividedHueRaw(index * 4))
  const attachmentHue = attachmentsHue[activeAttachmentIndex]

  return (
    <div className="modal-content standard-box flex flex-col items-center w-screen"

      style={
        {
          backgroundColor: color({
            lightness: 97,
            chroma: 0.04,
            hue: attachmentHue
          })
        }
      }
    >
      {/* Top-level tabs for each attachment */}
      <div role="tablist" aria-label="Attachment Sections" className="flex gap-2 border-b border-base-300 overflow-x-auto">
        {attachmentInfos.map((attachment, i) => (
          <button
            key={`${attachment.attachment_uuid}_top`}
            role="tab"
            aria-selected={activeAttachmentIndex === i}
            onClick={() => changeActiveAttachment(i)}
            style={
              {
                backgroundColor: color({
                  lightness: 80 + 10 * Number(activeAttachmentIndex === i),
                  chroma: 0.07,
                  hue: attachmentsHue[i]
                })
              }
            }
            className={clsx(`px-6 py-3 font-bold flex-shrink-0`, activeAttachmentIndex === i ? 'border-b-2 border-primary text-primary' : 'hover:bg-base-200')}
          // className={clsx(`px-6 py-3 font-bold flex-shrink-0 bg-[oklch(78.2% 0.066 180)]`, activeAttachmentIndex === i ? 'border-b-2 border-primary text-primary' : 'hover:bg-base-200')}
          // className={clsx(`px-6 py-3 font-bold flex-shrink-0 bg-[oklch(80% 0.12 ${attachmentsHue[i]})]`, activeAttachmentIndex === i ? 'border-b-2 border-primary text-primary' : 'hover:bg-base-200')}
          >
            {attachment.attachment_name}
          </button>
        ))}
      </div>

      {/* Sub-tabs for the active attachment */}
      {attachmentInfos.length > 0 && (
        <>
          <div role="tablist"
            aria-label="Document Sections" className="flex gap-2 border-b border-base-300">
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
                docUUID={activeAttachment.attachment_uuid}
                extension={fileExtensionFromText(activeAttachment.attachment_extension)}
              />
            }
            {showText && activeTabPerAttachment[activeAttachmentIndex] === 'text' &&
              <MarkdownContent docUUID={activeAttachment.attachment_uuid} />
            }
          </div>
        </>
      )}
    </div>
  );
};
