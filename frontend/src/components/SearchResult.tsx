import { Card, Modal, ModalClose, ModalDialog } from "@mui/joy";
import { useState } from "react";
import ResultModal from "./ResultModal";
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

  return (
    <>
      <Card
        style={{
          padding: "15px",
          backgroundColor: "white",
          borderRadius: "10px",
          border: "2px solid grey",
          width: "90%",
          maxHeight: "15em",
        }}
        onClick={() => setOpen(true)}
      >
        <h1>{data.name}</h1>
        <span />
        <div dangerouslySetInnerHTML={{ __html: data.text }} />
        <span />
        <p>{data.docketID}</p>
      </Card>
      <Modal
        aria-labelledby="modal-title"
        aria-describedby="modal-desc"
        open={open}
        onClose={() => setOpen(false)}
        sx={{ display: "flex", justifyContent: "center", alignItems: "center" }}
      >
        <ModalDialog className="standard-box">
          <ModalClose />
          <ResultModal open={open} />
        </ModalDialog>
      </Modal>
    </>
  );
};

export default SearchResult;
