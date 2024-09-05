// Modal.js
import React, { useEffect } from "react";
import { Box, Tab, Tabs, TabList } from "@mui/joy";
import axios from "axios";

type ModalProps = {
  open: boolean;
  objectId?: string;
  children?: React.ReactNode;
};

const ResultModal = ({ open, objectId, children }: ModalProps) => {
  const [loading, setLoading] = React.useState(false);

  const [title, setTitle] = React.useState("Demo Title");

  const peekDocument = async () => {
    try {
      const response = await axios.post(
        "http://localhost:4041/documents/peek",
        {
          id: objectId,
        }
      );
      if (response.data.length === 0) {
        return;
      }
    } catch (error) {
      console.log(error);
    } finally {
    }
  };

  useEffect(() => {
    if (open) {
      peekDocument();
    }
  });

  return (
    <div
      className="modal-content standard-box"
      style={{
        minHeight: "80vh",
        minWidth: "60vw",
      }}
    >
      {/* box containing the  */}
      {/* children are components passed to result modals from special searches */}

      <div className="card-title">
        {!loading ? <h1>{title}</h1> : <h1>{title}</h1>}
      </div>
      {/* passed content for filter */}
      {children ? <div>{children}</div> : null}
      <div className="modal-body">
        <Tabs>
          <TabList>
            <Tab variant="plain" color="neutral">
              ...
            </Tab>
          </TabList>
        </Tabs>
      </div>
    </div>
  );
};

export default ResultModal;
