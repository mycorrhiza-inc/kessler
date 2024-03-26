import React, { useState } from "react";
import Card from "@mui/joy/Card";

import Typography from "@mui/joy/Typography";

import Button from "@mui/joy/Button";

import Modal from "@mui/joy/Modal";
import Sheet from "@mui/joy/Sheet";
import ModalClose from "@mui/joy/ModalClose";

type DocumentNodeProps = { 
  data: any; 
  isConnectable: boolean 
};

function DocumentNode({ data, isConnectable }: DocumentNodeProps) {
  const [expanded, setExpanded] = useState(false);
  const [viewRawdoc, setViewRawdoc] = useState(false);
  return (
    <>
      <Card variant="plain">
        <Typography level="title-md">{data.docid.metadata.title}</Typography>
        <Typography level="body-sm">
          {data.docid.extras.short_summary}
        </Typography>
        <div>
          <Button
            color="primary"
            onClick={function () {
              setExpanded(expanded == false);
            }}
          >
            {expanded ? "Collapse" : "Expand"}
          </Button>

          <Button
            color="success"
            onClick={() => {
              setViewRawdoc(true);
            }}
          >
            View Raw Document
          </Button>
        </div>
      </Card>
      <Modal
        aria-labelledby="modal-title"
        aria-describedby="modal-desc"
        open={viewRawdoc}
        onClose={() => setViewRawdoc(false)}
        sx={{ display: "flex", justifyContent: "center", alignItems: "center" }}
      >
        <Sheet
          variant="outlined"
          sx={{
            maxWidth: 500,
            borderRadius: "md",
            p: 3,
            boxShadow: "lg",
          }}
        >
          <ModalClose variant="plain" sx={{ m: 1 }} />
          <Typography
            component="h2"
            id="modal-title"
            level="h4"
            textColor="inherit"
            fontWeight="lg"
            mb={1}
          >
            {data.docid.metadata.title}
          </Typography>
          <Typography id="modal-desc" textColor="text.tertiary">
            {data.document_text}
          </Typography>
        </Sheet>
      </Modal>
    </>
  );
}

export default DocumentNode;
