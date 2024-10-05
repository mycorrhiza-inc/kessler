import { useState } from "react";
import DocumentModalBody from "./DocumentModalBody";
import Modal from "./styled-components/Modal";
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
  return (
    <>
      <div
        className="card w-[90%] shadow-xl dark:card-bordered"
        onClick={() => {
          setOpen(true);
        }}
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
      <Modal open={open} setOpen={setOpen} uuid={data.sourceID}>
        <DocumentModalBody open={open} objectId={data.sourceID} />
      </Modal>
    </>
  );
};

export default SearchResult;
