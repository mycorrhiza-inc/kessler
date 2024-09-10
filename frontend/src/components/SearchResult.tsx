import { Card, Modal, ModalClose, ModalDialog } from "@mui/joy";
import { useState } from "react";
import DocumentModalBody from "./DocumentModalBody";
import zIndex from "@mui/material/styles/zIndex";
type SearchFields = {
  id: string;
  name: string;
  text: string;
  docketID: string;
};

type SearchResultProps = {
  data: SearchFields;
};

const SearchResult = ({ data }: SearchResultProps) => {
  const [open, setOpen] = useState(false);
  // Huge fan of dasiui for refactoring the card here, easy extensionality
  return (
    <>
      <div
        className="card w-[90%] shadow-xl dark:card-bordered"
        onClick={() => setOpen(true)}
      >
        <div className="card-body">
          <h2 className="card-title">
            <h1>{data.name}</h1>
          </h2>
          <span />
          <div dangerouslySetInnerHTML={{ __html: data.text }} />
          <span />
          <p>{data.docketID}</p>
        </div>
      </div>
      <Modal
        aria-labelledby="modal-title"
        aria-describedby="modal-desc"
        open={open}
        onClose={() => setOpen(false)}
        sx={{ display: "flex", justifyContent: "center", alignItems: "center" }}
        style={{ zIndex: 99 }}
      >
        <ModalDialog className="standard-box">
          <ModalClose />
          <DocumentModalBody open={open} objectId={data.id} />
        </ModalDialog>
      </Modal>
    </>
  );
};

export default SearchResult;
