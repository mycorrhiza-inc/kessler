import React, { Suspense } from "react";
import * as Tabs from "@radix-ui/react-tabs";

import PDFViewer from "./PDFViewer";
import MarkdownRenderer from "@/components/MarkdownRenderer";
import LoadingSpinner from "@/components/styled-components/LoadingSpinner";
import { apiURL } from "@/lib/env_variables";
import { fetchObjectDataFromURL, fetchTextDataFromURL } from "./documentLoader";

type ModalProps = {
  open: boolean;
  objectId: string;
  overridePDFUrl?: string;
  children?: React.ReactNode;
  title?: string;
};

const MarkdownContent = async ({ docUUID }: { docUUID: string }) => {
  const markdown_url = `${apiURL}/v2/public/files/${docUUID}/markdown`;
  // axios.get(`https://api.kessler.xyz/v2/public/files/${objectid}/markdown`),
  const text = await fetchTextDataFromURL(markdown_url);
  return <MarkdownRenderer>{text}</MarkdownRenderer>;
};

const MetadataContent = async ({ docUUID }: { docUUID: string }) => {
  const object_url = `${apiURL}/v2/public/files/${docUUID}/markdown`;
  // axios.get(`https://api.kessler.xyz/v2/public/files/${objectid}/markdown`),
  const mdata = await fetchObjectDataFromURL(object_url);

  return (
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
  );
};
const PDFContent = ({
  docUUID,
  overridePDFUrl,
}: {
  docUUID: string;
  overridePDFUrl?: string;
}) => {
  const pdfUrl = overridePDFUrl || `${apiURL}/v2/public/files/${docUUID}/raw`;
  return <PDFViewer file={pdfUrl}></PDFViewer>;
};

const DocumentModalBody = ({
  open,
  objectId,
  children,
  title,
  overridePDFUrl,
}: ModalProps) => {
  return (
    <div className="modal-content standard-box ">
      {/* children are components passed to result modals from special searches */}
      <div className="card-title">
        <h1>{title}</h1>
      </div>
      {/* Deleted all the MUI stuff, this should absolutely be refactored into its own styled component soonish*/}

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
          <PDFContent docUUID={objectId} overridePDFUrl={overridePDFUrl} />
        </Tabs.Content>
        <Tabs.Content className="TabsContent" value="tab2">
          <Suspense
            fallback={<LoadingSpinner loadingText="Loading Document Text" />}
          >
            <MarkdownContent docUUID={objectId} />
          </Suspense>
        </Tabs.Content>
        <Tabs.Content className="TabsContent" value="tab3">
          <div className="overflow-x-auto">
            <Suspense
              fallback={
                <LoadingSpinner loadingText="Loading Document Metadata" />
              }
            >
              <MetadataContent docUUID={objectId} />
            </Suspense>
          </div>
        </Tabs.Content>
      </Tabs.Root>
    </div>
  );
};

export default DocumentModalBody;
