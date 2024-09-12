// Modal.js
import React, { useEffect } from "react";
import * as Tabs from "@radix-ui/react-tabs";
import axios from "axios";

import PDFViewer from "./PDFViewer";
import MarkdownRenderer from "./MarkdownRenderer";
import { testMarkdownContent } from "./MarkdownRenderer";

type ModalProps = {
  open: boolean;
  objectId?: string;
  children?: React.ReactNode;
  title?: string;
};

const DocumentModalBody = ({ open, objectId, children, title }: ModalProps) => {
  const [loading, setLoading] = React.useState(false);

  const [docText, setDocText] = React.useState("Loading Document Text");
  const [pdfUrl, setPdfUrl] = React.useState("");

  const peekDocument = async () => {
    try {
      const response = await axios.post(
        "http://localhost:4041/documents/peek",
        {
          id: objectId,
        },
      );
      if (response.data.length === 0) {
        return;
      }
    } catch (error) {
      console.log(error);
    } finally {
    }
  };
  const getDocumentText = async () => {
    const response = await axios.get(`/api/v1/files/markdown/${objectId}`);
    setDocText(response.data);
  };

  useEffect(() => {
    if (open) {
      setPdfUrl(`/api/v1/files/raw/${objectId}`);
      getDocumentText();
    }
  }, [open]);

  return (
    <div className="modal-content standard-box ">
      {/* children are components passed to result modals from special searches */}
      <div className="card-title">
        {!loading ? <h1>{title}</h1> : <h1>{title}</h1>}
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
          <PDFViewer file={pdfUrl}></PDFViewer>
        </Tabs.Content>
        <Tabs.Content className="TabsContent" value="tab2">
          <MarkdownRenderer>{docText}</MarkdownRenderer>
        </Tabs.Content>
        <Tabs.Content className="TabsContent" value="tab3">
          <div className="overflow-x-auto">
            <table className="table table-zebra">
              {/* head */}
              <thead>
                <tr>
                  <th></th>
                  <th>Name</th>
                  <th>Job</th>
                  <th>Favorite Color</th>
                </tr>
              </thead>
              <tbody>
                {/* row 1 */}
                <tr>
                  <th>1</th>
                  <td>Cy Ganderton</td>
                  <td>Quality Control Specialist</td>
                  <td>Blue</td>
                </tr>
                {/* row 2 */}
                <tr>
                  <th>2</th>
                  <td>Hart Hagerty</td>
                  <td>Desktop Support Technician</td>
                  <td>Purple</td>
                </tr>
                {/* row 3 */}
                <tr>
                  <th>3</th>
                  <td>Brice Swyre</td>
                  <td>Tax Accountant</td>
                  <td>Red</td>
                </tr>
              </tbody>
            </table>
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
