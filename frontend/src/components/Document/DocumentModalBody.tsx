"use client";
// TODO: change the network fetch stuff so that this can be SSR'd
import React, { Suspense, memo, useEffect } from "react";
import * as Tabs from "@radix-ui/react-tabs";

import PDFViewer from "./PDFViewer";
import MarkdownRenderer from "@/components/MarkdownRenderer";
import LoadingSpinner from "@/components/styled-components/LoadingSpinner";
import { apiURL } from "@/lib/env_variables";
import { fetchObjectDataFromURL, fetchTextDataFromURL } from "./documentLoader";
import useSWRImmutable from "swr";
import Link from "next/link";
import { completeFileSchemaGet } from "@/lib/requests/search";
import { AuthorInfoPill } from "../Conversations/TextPills";

// import { ErrorBoundary } from "react-error-boundary";

type ModalProps = {
  open: boolean;
  objectId: string;
  children?: React.ReactNode;
  title?: string;
  isPage: boolean;
  setTitle?: (val: string) => void;
};

const MarkdownContent = memo(({ docUUID }: { docUUID: string }) => {
  const markdown_url = `${apiURL}/v2/public/files/${docUUID}/markdown`;
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

const MetadataContent = memo(
  ({
    metadata,
    isLoading,
    error,
  }: {
    metadata: any;
    isLoading: boolean;
    error: any;
  }) => {
    if (isLoading) {
      return <LoadingSpinner loadingText="Loading Document Metadata" />;
    }
    if (error) {
      return <p>Encountered an error getting text from the server.</p>;
    }

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
  },
);

const PDFContent = ({ docUUID }: { docUUID: string }) => {
  const [loading, setLoading] = React.useState(true);
  const pdfUrl = `${apiURL}/v2/public/files/${docUUID}/raw`;
  return (
    <>
      {/* This apparently gets an undefined network error when trying to fetch the pdf from their website not exactly sure why, we need to get the s3 fetch working in golang */}
      <PDFViewer file={pdfUrl} setLoading={setLoading}></PDFViewer>

      {loading && <LoadingSpinner loadingText="PDF Viewer Loading" />}
    </>
  );
};

const DocumentModalBody = ({
  open,
  objectId,
  title,
  setTitle,
  isPage,
}: ModalProps) => {
  const semiCompleteFileUrl = `${apiURL}/v2/public/files/${objectId}`;
  const { data, error, isLoading } = useSWRImmutable(
    semiCompleteFileUrl,
    completeFileSchemaGet,
  );
  const actualTitle: string = (data?.name || title) as string;
  const underscoredTitle = title ? title.replace(/ /g, "_") : "Unkown_Document";
  const fileUrlNamedDownload = `${apiURL}/v2/public/files/${objectId}/raw/${underscoredTitle}.pdf`;
  const kesslerFileUrl = `/files/${objectId}`;
  const metadata = data?.mdata;
  useEffect(() => {
    if (setTitle) {
      setTitle(actualTitle);
    }
  }, [actualTitle]);
  return (
    <div className="modal-content standard-box ">
      <div className="card-title flex justify-between items-center">
        <h1>{actualTitle}</h1>
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

      {isLoading ? (
        <LoadingSpinner loadingText="Loading Document Summary" />
      ) : (
        <>
          <h3 className="text-lg font-bold">Summary:</h3>
          <br />
          <MarkdownRenderer>{data?.extra.summary as string}</MarkdownRenderer>
          <br />
          <h3 className="text-lg font-bold">Authors:</h3>
          <br />
          {data?.authors?.map((auth_info) => (
            <AuthorInfoPill author_info={auth_info} />
          ))}
        </>
      )}

      <Tabs.Root
        className="TabsRoot"
        role="tablist tabs tabs-bordered tabs-lg"
        defaultValue="tab1"
      >
        <Tabs.List className="TabsList" aria-label="What Documents ">
          <Tabs.Trigger className="TabsTrigger tab" value="tab1">
            Document
          </Tabs.Trigger>
          <Tabs.Trigger className="TabsTrigger tab" value="tab2">
            Document Text
          </Tabs.Trigger>
          <Tabs.Trigger className="TabsTrigger tab" value="tab3">
            Metadata
          </Tabs.Trigger>
        </Tabs.List>
        <Tabs.Content className="TabsContent" value="tab1">
          <PDFContent docUUID={objectId} />
        </Tabs.Content>
        <Tabs.Content className="TabsContent" value="tab2">
          <MarkdownContent docUUID={objectId} />
        </Tabs.Content>
        <Tabs.Content className="TabsContent" value="tab3">
          <MetadataContent
            isLoading={isLoading}
            metadata={metadata}
            error={error}
          />
        </Tabs.Content>
      </Tabs.Root>
    </div>
  );
};

export default DocumentModalBody;
