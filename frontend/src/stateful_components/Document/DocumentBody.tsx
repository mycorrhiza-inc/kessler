"use client";
// TODO: change the network fetch stuff so that this can be SSR'd
import React, { memo } from "react";
import * as Tabs from "@radix-ui/react-tabs";

import PDFViewer from "./PDFViewer";
import { fetchTextDataFromURL } from "./documentLoader";
import useSWRImmutable from "swr";
import {
  AuthorInformation,
  CompleteFileSchema,
} from "@/lib/types/backend_schemas";
import { CLIENT_API_URL } from "@/lib/env_variables";
import { FileExtension } from "@/style_components/Pills/FileExtension";
import XlsxViewer from "@/style_components/messages/XlsxCannotBeViewedMessage";
import ErrorMessage from "@/style_components/messages/ErrorMessage";
import { DocketPill } from "@/style_components/Pills/TextPills";
import LoadingSpinner from "@/style_components/misc/LoadingSpinner";
import MarkdownRenderer from "@/style_components/misc/MarkdownRenderer";


// import { ErrorBoundary } from "react-error-boundary";
//

const MarkdownContent = memo(({ docUUID }: { docUUID: string }) => {
  const markdown_url = `${CLIENT_API_URL}/v2/public/files/${docUUID}/markdown`;
  // axios.get(`https://api.kessler.xyz/v2/public/files/${objectid}/markdown`),
  const { data, error, isLoading } = useSWRImmutable(
    markdown_url,
    fetchTextDataFromURL,
  );
  if (isLoading) {
    return <LoadingSpinner loadingText="Loading Document" />;
  }
  const text = data;
  if (error) {
    return (
      <p>
        Encountered an error getting text from the server.
        <br /> {String(error)}
      </p>
    );
  }
  if (text == undefined) {
    return <p>Document text is undefined.</p>;
  }
  return <MarkdownRenderer>{text}</MarkdownRenderer>;
});

const MetadataContent = memo(({ metadata }: { metadata: any }) => {
  // if (isLoading) {
  //   return <LoadingSpinner loadingText="Loading Document Metadata" />;
  // }
  // if (error) {
  //   return <p>Encountered an error getting text from the server.</p>;
  // }

  return (
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
  );
});

const DocumentContent = ({
  docUUID,
  extension,
}: {
  docUUID: string;
  extension: FileExtension;
}) => {
  const documentUrl = `${CLIENT_API_URL}/v2/public/files/${docUUID}/raw`;
  if (extension == FileExtension.PDF) {
    return <PDFViewer file={documentUrl} />;
  }
  if (extension == FileExtension.XLSX) {
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
  const fileUrlNamedDownload = `${CLIENT_API_URL}/v2/public/files/${objectId}/raw/${underscoredTitle}.${extension}`;
  const kesslerFileUrl = `/files/${objectId}`;
  const authors_unpluralized =
    documentObject.authors?.length == 1 ? "Author" : "Authors";
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
          {/* {verified && ( */}
          {/* <ExperimentalChatModalClickDiv */}
          {/*   className="btn btn-accent" */}
          {/*   inheritedFilters={[ */}
          {/*     { filter: FilterField.MatchFileUUID, value: objectId }, */}
          {/*   ]} */}
          {/* > */}
          {/*   Chat with Document */}
          {/* </ExperimentalChatModalClickDiv> */}
          {/* )} */}
        </div>
      </div>
      <p>
        <b>Case Number:</b> {"   "}
        <DocketPill
          docketId={documentObject.mdata.docket_id as string}
        />
      </p>
      {documentObject.authors && (
        <p>
          <b>{authors_unpluralized}:</b>{" "}
          {documentObject.authors.map((a: AuthorInformation) => (
            <AuthorInfoPill author_info={a} />
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
  var metadata = documentObject.mdata;
  metadata.hash = documentObject.hash;
  // TODO: Make this into a library function or something.
  const showText =
    documentObject?.verified && documentObject?.extension != "xlsx";
  // Temporary fix while we sort out the bad documents
  const showRawDocument = documentObject?.verified;
  const extension = fileExtensionFromText(documentObject?.extension);
  const getDefaultTab = (
    showRawDocument: boolean,
    showText: boolean,
  ): string => {
    if (!showRawDocument) {
      if (!showText) {
        return "tab3";
      }
      return "tab2";
    }
    return "tab1";
  };
  const defaultTab = getDefaultTab(showRawDocument, showText);
  return (
    <div className="modal-content standard-box ">
      <DocumentHeader documentObject={documentObject} isPage={isPage} />
      <Tabs.Root
        className="TabsRoot"
        role="tablist tabs tabs-bordered tabs-lg"
        defaultValue={defaultTab}
      >
        <Tabs.List
          className="TabsList flex gap-2 border-b border-base-300"
          aria-label="Document Sections"
        >
          {showRawDocument && (
            <Tabs.Trigger
              className="TabsTrigger px-6 py-3 font-bold hover:bg-base-200 transition-colors duration-200 data-[state=active]:border-b-2 data-[state=active]:border-primary data-[state=active]:text-primary"
              value="tab1"
            >
              📄 Document
            </Tabs.Trigger>
          )}

          {showText && (
            <Tabs.Trigger
              className="TabsTrigger px-6 py-3 font-bold hover:bg-base-200 transition-colors duration-200 data-[state=active]:border-b-2 data-[state=active]:border-primary data-[state=active]:text-primary"
              value="tab2"
            >
              📝 Document Text
            </Tabs.Trigger>
          )}
          <Tabs.Trigger
            className="TabsTrigger px-6 py-3 font-bold hover:bg-base-200 transition-colors duration-200 data-[state=active]:border-b-2 data-[state=active]:border-primary data-[state=active]:text-primary"
            value="tab3"
          >
            ℹ️ Metadata
          </Tabs.Trigger>
        </Tabs.List>
        {showRawDocument && (
          <Tabs.Content className="TabsContent" value="tab1">
            <DocumentContent docUUID={objectId} extension={extension} />
          </Tabs.Content>
        )}
        {showText && (
          <Tabs.Content className="TabsContent" value="tab2">
            <MarkdownContent docUUID={objectId} />
          </Tabs.Content>
        )}
        <Tabs.Content className="TabsContent" value="tab3">
          <MetadataContent metadata={metadata} />
        </Tabs.Content>
      </Tabs.Root>
    </div>
  );
};
