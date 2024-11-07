import React, { Suspense } from "react";
import * as Tabs from "@radix-ui/react-tabs";
import { fetchDocumentData } from "@/utils/documentLoader";

import PDFViewer from "./PDFViewer";
import MarkdownRenderer from "@/bin/MarkdownRenderer";
import { testMarkdownContent } from "./MarkdownRenderer";
import { LoadingSpinner } from "./styled-components/LoadingSpinner";

type ModalProps = {
  open: boolean;
  objectId?: string;
  overridePDFUrl?: string;
  children?: React.ReactNode;
  title?: string;
};

const DocumentModalBody = ({
  open,
  objectId,
  children,
  title,
  overridePDFUrl,
}: ModalProps) => {
  const resource = React.useMemo(
    () => fetchDocumentData(objectId!, overridePDFUrl),
    [objectId, overridePDFUrl],
  );

  const { text: docText, metadata: docMetadata, pdfUrl } = resource.read();

  return (
    <div className="modal-content standard-box ">
      <></>
      {/* children are components passed to result modals from special searches */}
      <div className="card-title">
        <h1>{title}</h1>
      </div>
      {/* REFACTOR MODALS to fix stupid BS with MUI       {children ? <div>{children}</div> : null}*/}
      {/* Deleted all the MUI stuff, this should absolutely be refactored into its own styled component soonish*/}

      <Tabs.Root
        className="TabsRoot"
        role="tablist tabs tabs-bordered tabs-lg"
        defaultValue="tab1"
      >
        <Tabs.List className="TabsList" aria-label="Manage your account">
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
          <Suspense fallback={<LoadingSpinner text="Loading PDF" />}>
            <PDFViewer file={pdfUrl}></PDFViewer>
          </Suspense>
        </Tabs.Content>
        <Tabs.Content className="TabsContent" value="tab2">
          <Suspense fallback={<LoadingSpinner text="Loading Document Text" />}>
            <MarkdownRenderer>{docText}</MarkdownRenderer>
          </Suspense>
        </Tabs.Content>
        <Tabs.Content className="TabsContent" value="tab3">
          <div className="overflow-x-auto">
            <Suspense
              fallback={<LoadingSpinner text="Loading Document Metadata" />}
            >
              <table className="table table-zebra">
                {/* head */}
                <thead>
                  <tr>
                    <th>Field</th>
                    <th>Value</th>
                  </tr>
                </thead>
                <tbody>
                  {Object.entries(docMetadata).map(([key, value]) => (
                    <tr key={key}>
                      <td>{key}</td>
                      <td>{String(value)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </Suspense>
          </div>
        </Tabs.Content>
      </Tabs.Root>
    </div>
  );
};

// <Tabs aria-label="Basic tabs" defaultValue={0}>
//   <TabList>
//     <Tab>Document Text</Tab>
//     <Tab>Document</Tab>
//     <Tab>Metadata</Tab>
//   </TabList>
//   <TabPanel value={0}>
//     <MarkdownRenderer>{docText}</MarkdownRenderer>
//   </TabPanel>
//   <TabPanel value={1}>
//     <b>Second</b> tab panel
//   </TabPanel>
//   <TabPanel value={2}>
//     <b>Third</b> tab panel
//   </TabPanel>
// </Tabs>
export default DocumentModalBody;
