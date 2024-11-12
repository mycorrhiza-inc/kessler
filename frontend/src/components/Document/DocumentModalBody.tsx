import React, { Suspense, memo } from "react";
import * as Tabs from "@radix-ui/react-tabs";

import PDFViewer from "./PDFViewer";
import MarkdownRenderer from "@/components/MarkdownRenderer";
import LoadingSpinner from "@/components/styled-components/LoadingSpinner";
import { apiURL } from "@/lib/env_variables";
import { fetchObjectDataFromURL, fetchTextDataFromURL } from "./documentLoader";
import useSWR from "swr";

// import { ErrorBoundary } from "react-error-boundary";

type ModalProps = {
  open: boolean;
  objectId: string;
  overridePDFUrl?: string;
  children?: React.ReactNode;
  title?: string;
};

const MarkdownContentRaw = memo(({ docUUID }: { docUUID: string }) => {
  const markdown_url = `${apiURL}/v2/public/files/${docUUID}/markdown`;
  // axios.get(`https://api.kessler.xyz/v2/public/files/${objectid}/markdown`),
  const { data, error, isLoading } = useSWR(markdown_url, fetchTextDataFromURL);
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

const MarkdownContent = (props: { docUUID: string }) => {
  return <MarkdownContentRaw {...props} />;
};

const MetadataContentRaw = memo(({ docUUID }: { docUUID: string }) => {
  const object_url = `${apiURL}/v2/public/files/${docUUID}/metadata`;
  // axios.get(`https://api.kessler.xyz/v2/public/files/${objectid}/markdown`),
  const { data, error, isLoading } = useSWR(object_url, fetchObjectDataFromURL);
  if (isLoading) {
    return <LoadingSpinner loadingText="Loading Document Metadata" />;
  }
  if (error) {
    return <p>Encountered an error getting text from the server.</p>;
  }
  if (typeof data !== "object") {
    console.log("mdata type: ", typeof data);
    console.log("mdata: ", data);
    return <p>Expected an object for metadata, got something else.</p>;
  }
  const mdata_str = atob(data.Mdata);
  console.log("mdata str: ", mdata_str);
  const mdata = JSON.parse(mdata_str);
  console.log("metadata: ", mdata);
  // mdata.Mdata = metadata;

  return (
    <div className="overflow-x-auto">
      <table className="table table-zebra">
        {/* head */}
        <thead>
          <tr>
            <th>Field</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          {Object.entries(mdata).map(([key, value]) => (
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

const MetadataContent = (props: { docUUID: string }) => {
  return <MetadataContentRaw {...props} />;
};
const PDFContent = ({
  docUUID,
  overridePDFUrl,
}: {
  docUUID: string;
  overridePDFUrl?: string;
}) => {
  const pdfUrl = overridePDFUrl || `${apiURL}/v2/public/files/${docUUID}/raw`;
  // const pdfUrl = `${apiURL}/v2/public/files/${docUUID}/raw`;
  // return (
  //   <a className="btn btn-primary" href={pdfUrl} target="_blank">
  //     PDF Viewer coming soon
  //   </a>
  // );
  // return <PDFViewer file={pdfUrl}></PDFViewer>;
  return (
    <>
      <LoadingSpinner loadingText="PDF Viewer coming soon" />
      {/* This apparently gets an undefined network error when trying to fetch the pdf from their website not exactly sure why, we need to get the s3 fetch working in golang */}
      <PDFViewer file={pdfUrl}></PDFViewer>
    </>
  );
};

const DocumentModalBody = ({
  open,
  objectId,
  children,
  title,
  overridePDFUrl,
}: ModalProps) => {
  const pdfUrl = overridePDFUrl || `${apiURL}/v2/public/files/${objectId}/raw`;
  return (
    <div className="modal-content standard-box ">
      {/* children are components passed to result modals from special searches */}
      <div className="card-title">
        <h1>{title}</h1>
      </div>
      {/* Deleted all the MUI stuff, this should absolutely be refactored into its own styled component soonish*/}

      <a className="btn btn-primary" href={pdfUrl} target="_blank">
        Download PDF
      </a>

      <Tabs.Root
        className="TabsRoot"
        role="tablist tabs tabs-bordered tabs-lg"
        defaultValue="tab2"
      >
        <Tabs.List className="TabsList" aria-label="What Documents ">
          <Tabs.Trigger className="TabsTrigger tab" value="tab2">
            Document Text
          </Tabs.Trigger>
          <Tabs.Trigger className="TabsTrigger tab" value="tab3">
            Metadata
          </Tabs.Trigger>
          <Tabs.Trigger className="TabsTrigger tab" value="tab1">
            Document
          </Tabs.Trigger>
        </Tabs.List>
        <Tabs.Content className="TabsContent" value="tab1">
          <PDFContent docUUID={objectId} overridePDFUrl={overridePDFUrl} />
        </Tabs.Content>
        <Tabs.Content className="TabsContent" value="tab2">
          <MarkdownContent docUUID={objectId} />
        </Tabs.Content>
        <Tabs.Content className="TabsContent" value="tab3">
          <MetadataContent docUUID={objectId} />
        </Tabs.Content>
      </Tabs.Root>
    </div>
  );
};

export default DocumentModalBody;
