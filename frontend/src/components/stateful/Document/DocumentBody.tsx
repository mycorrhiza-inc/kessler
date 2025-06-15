"use client";

import Link from 'next/link';
import React, { memo, useState } from "react";

import PDFViewer from "./PDFViewer";
import { fetchTextDataFromURL } from "./documentLoader";
import {
  AuthorInformation,
  CompleteFileSchema,
} from "@/lib/types/backend_schemas";
import { CLIENT_API_URL } from "@/lib/env_variables";
import { FileExtension, fileExtensionFromText } from "@/componenets/style/Pills/FileExtension";
import XlsxViewer from "@/componenets/style/messages/XlsxCannotBeViewedMessage";
import ErrorMessage from "@/componenets/style/messages/ErrorMessage";
import { AuthorInfoPill, DocketPill } from "@/componenets/style/Pills/TextPills";
import LoadingSpinner from "@/componenets/style/misc/LoadingSpinner";
import MarkdownRenderer from "@/componenets/style/misc/MarkdownRenderer";

const MarkdownContent = memo(({ docUUID }: { docUUID: string }) => {
  // TODO: Replace this placeholder with real fetch + SSR support
  return <p>IMPLEMENT SOME KIND OF NOT SHITTY FETCHING FRAMEWORK FOR THIS, AND ALSO LET IT BE SSRABLE</p>;
});

const MetadataContent = memo(({ metadata }: { metadata: any }) => (
  <div className="overflow-x-auto">
    <table className="table table-zebra">
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
  const documentUrl = `${CLIENT_API_URL}/v2/public/files/${docUUID}/raw`;
  if (extension === FileExtension.PDF) {
    return <PDFViewer file={documentUrl} />;
  }
  if (extension === FileExtension.XLSX) {
    return <XlsxViewer file={documentUrl} />;
  }
  return <ErrorMessage error={`Cannot Display Document Type`} />;
};

const DocumentHeader = ({
  documentObject,
  isPage,
}: {
  documentObject: CompleteFileSchema;
  isPage: boolean;
}) => {
  const title: string = documentObject.name;
  const objectId: string = documentObject.id;
  const extension = documentObject.extension || "pdf";
  const verified = (documentObject.verified || false) as boolean;
  const summary = documentObject.extra.summary;
  const underscoredTitle = title ? title.replace(/ /g, "_") : "Unkown_Document";
  const fileUrlNamedDownload =
    `${CLIENT_API_URL}/v2/public/files/${objectId}/raw/${underscoredTitle}.${extension}`;
  const kesslerFileUrl = `/files/${objectId}`;
  const authors_unpluralized =
    documentObject.authors?.length === 1 ? "Author" : "Authors";

  return (
    <>
      <div className="card-title flex justify-between items-start">
        <h1 className="text-3xl break-words max-w-[70%]">{title}</h1>
        <div className="flex gap-2">
          <a
            className="btn btn-primary"
            href={fileUrlNamedDownload}
            target="_blank"
            download={title}
          >
            Download File
          </a>
          {!isPage && (
            <Link
              className="btn btn-secondary"
              href={kesslerFileUrl}
              target="_blank"
            >
              Open in New Tab
            </Link>
          )}
        </div>
      </div>
      <p>
        <b>Case Number:</b> {'   '}
        <DocketPill docketId={documentObject.mdata.docket_id as string} />
      </p>
      {documentObject.authors && (
        <p>
          <b>{authors_unpluralized}:</b>{' '}
          {documentObject.authors.map((a: AuthorInformation) => (
            <AuthorInfoPill author_info={a} key={a.author_id} />
          ))}
        </p>
      )}
      <div className="p-4" />
      <h2 className="text-xl">
        <b>LLM Summary:</b>
      </h2>
      <MarkdownRenderer>
        {verified
          ? summary
          : "The document hasnt finished processing yet. Come back later for a completed summary and document chat!"}
      </MarkdownRenderer>
      <div className="p-12" />
    </>
  );
};

export const DocumentMainTabs = ({
  documentObject,
  isPage,
}: {
  documentObject: CompleteFileSchema;
  isPage: boolean;
}) => {
  const objectId = documentObject.id;
  const metadata = { ...documentObject.mdata, hash: documentObject.hash };

  const showText =
    documentObject.verified && documentObject.extension !== "xlsx";
  const showRawDocument = documentObject.verified;
  const extension = fileExtensionFromText(documentObject.extension);

  type TabValue = "tab1" | "tab2" | "tab3";
  const getDefaultTab = (): TabValue => {
    if (showRawDocument) return "tab1";
    if (showText) return "tab2";
    return "tab3";
  };

  const [activeTab, setActiveTab] = useState<TabValue>(getDefaultTab());

  return (
    <div className="modal-content standard-box">
      <DocumentHeader documentObject={documentObject} isPage={isPage} />

      {/* Tabs List */}
      <div
        role="tablist"
        aria-label="Document Sections"
        className="TabsList flex gap-2 border-b border-base-300"
      >
        {showRawDocument && (
          <button
            role="tab"
            aria-selected={activeTab === "tab1"}
            onClick={() => setActiveTab("tab1")}
            className={`px-6 py-3 font-bold hover:bg-base-200 transition-colors duration-200 ${
              activeTab === "tab1"
                ? "border-b-2 border-primary text-primary"
                : ""
            }`}
          >
            📄 Document
          </button>
        )}

        {showText && (
          <button
            role="tab"
            aria-selected={activeTab === "tab2"}
            onClick={() => setActiveTab("tab2")}
            className={`px-6 py-3 font-bold hover:bg-base-200 transition-colors duration-200 ${
              activeTab === "tab2"
                ? "border-b-2 border-primary text-primary"
                : ""
            }`}
          >
            📝 Document Text
          </button>
        )}

        <button
          role="tab"
          aria-selected={activeTab === "tab3"}
          onClick={() => setActiveTab("tab3")}
          className={`px-6 py-3 font-bold hover:bg-base-200 transition-colors duration-200 ${
            activeTab === "tab3"
              ? "border-b-2 border-primary text-primary"
              : ""
          }`}
        >
          ℹ️ Metadata
        </button>
      </div>

      {/* Tabs Content */}
      <div className="mt-4">
        {showRawDocument && activeTab === "tab1" && (
          <div className="TabsContent">
            <DocumentContent docUUID={objectId} extension={extension} />
          </div>
        )}

        {showText && activeTab === "tab2" && (
          <div className="TabsContent">
            <MarkdownContent docUUID={objectId} />
          </div>
        )}

        {activeTab === "tab3" && (
          <div className="TabsContent">
            <MetadataContent metadata={metadata} />
          </div>
        )}
      </div>
    </div>
  );
};
