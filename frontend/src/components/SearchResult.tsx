import { useState } from "react";
import DocumentModal from "./Document/DocumentModal";
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
          setOpen((prev) => !prev);
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
      <DocumentModal open={open} setOpen={setOpen} objectId={data.sourceID} />
    </>
  );
};

export default SearchResult;
