// Modal.js
import React, { useEffect } from "react";
import { Box, Tab, Tabs, TabList, TabPanel } from "@mui/joy";
import axios from "axios";

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
      peekDocument();
      getDocumentText();
    }
  }, [open]);

  return (
    <div
      className="modal-content standard-box "
      style={{
        minHeight: "80vh",
        minWidth: "60vw",
      }}
    >
      {/* children are components passed to result modals from special searches */}
      <div className="card-title">
        {!loading ? <h1>{title}</h1> : <h1>{title}</h1>}
      </div>
      {/* REFACTOR MODALS to fix stupid BS with MUI */}
      {children ? <div>{children}</div> : null}
      <div className="modal-body">
        <Tabs aria-label="Basic tabs" defaultValue={0}>
          <TabList>
            <Tab>Document Text</Tab>
            <Tab>Document</Tab>
            <Tab>Metadata</Tab>
          </TabList>
          <TabPanel value={0}>
            <MarkdownRenderer>{docText}</MarkdownRenderer>
          </TabPanel>
          <TabPanel value={1}>
            <b>Second</b> tab panel
          </TabPanel>
          <TabPanel value={2}>
            <b>Third</b> tab panel
          </TabPanel>
        </Tabs>
      </div>
    </div>
  );
};

export default DocumentModalBody;
