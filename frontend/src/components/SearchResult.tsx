import { Card, Modal, ModalClose, ModalDialog } from "@mui/joy";
import { useState } from "react";
import DocumentModalBody from "./DocumentModalBody";
import zIndex from "@mui/material/styles/zIndex";
type SearchFields = {
  sourceID: string;
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
  const docid: string = data.sourceID;
  const onClick = () => {
    setOpen(true);
    // Oh come on, the docid is always not null, its defined right below you
    // @ts-ignore
    document.getElementById(`doc_modal_${docid}`).showModal();
  };
  return (
    <>
      <div
        className="card w-[90%] shadow-xl dark:card-bordered"
        onClick={onClick}
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
      <dialog id={`doc_modal_${docid}`} className="modal ">
        <div
          className="modal-box bg-white dark:bg-black "
          style={{
            minHeight: "80vh",
            minWidth: "60vw",
          }}
          // This should just work and not require a background override, its an inidication something is deeply wrong
        >
          <form method="dialog">
            {/* if there is a button in form, it will close the modal */}
            <button className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">
              ✕
            </button>
          </form>
          <DocumentModalBody open={open} objectId={data.sourceID} />
        </div>
      </dialog>
    </>
  );
};
// <Modal
//   aria-labelledby="modal-title"
//   aria-describedby="modal-desc"
//   open={open}
//   onClose={() => setOpen(false)}
//   sx={{ display: "flex", justifyContent: "center", alignItems: "center" }}
//   style={{ zIndex: 99 }}
// >
//   <ModalDialog className="standard-box">
//     <ModalClose />
//   </ModalDialog>
// </Modal>

export default SearchResult;
