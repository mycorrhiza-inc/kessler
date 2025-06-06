import React, { useState } from "react";
import { Filing } from "@/lib/types/FilingTypes";
import { AuthorInformation } from "@/lib/types/backend_schemas";
import {
  AuthorInfoPill,
  DocketPill,
  FilePill,
  TextPill,
} from "./TextPills";
import DocumentModal from '../Document/DocumentModal';
import { TableStyle } from "../styles/Table";
import { fileExtensionFromText } from "./FileExtension";

const NoclickSpan = ({ children }: { children: React.ReactNode }) => {
  return <span className="noclick">{children}</span>;
};

const AuthorColumn = ({ filing }: { filing: Filing }) => {
  return (
    <>
      {filing.authors_information
        ? filing.authors_information.map((auth_info: AuthorInformation) => (
          <AuthorInfoPill author_info={auth_info} />
        ))
        : filing.author + " Something isnt working"}
    </>
  );
};

const TableRow = ({
  filing,
  DocketColumn,
  OpenRow
}: {
  filing: Filing;
  DocketColumn?: boolean;
  OpenRow: (id: string) => void
}) => {
  const [open, setOpen] = useState(false);
  const handleRowClick = (event: React.MouseEvent<HTMLTableRowElement>) => {
    // Check if the clicked element is the text inside the row
    const element = event.target as HTMLElement;
    const tagName = element.tagName; // Ensure type safety
    const className = element.className;
    console.log("element", element);
    console.log("tagName", tagName);
    console.log("className", className);

    const includesNoClick = className.includes("noclick");

    if (tagName === "SPAN" || tagName === "BUTTON" || includesNoClick) {
      console.log("Text was clicked");
      // Prevent the event from propagating further if needed
      event.stopPropagation();
      return;
    }

    // If it's not the specific text, allow the row click to proceed
    console.log("Row was clicked");
    OpenRow(filing.id);
  };

  return (
    <>
      <tr
        className="border-b border-base-300 hover:bg-base-200 transition duration-500 ease-out"
        onClick={handleRowClick}
      >
        <td>
          <NoclickSpan>{filing.date}</NoclickSpan>
        </td>
        <td>
          <NoclickSpan>
            <TextPill text={filing.file_class} />
          </NoclickSpan>
        </td>
        <td>
          {/* Removing the noclick around this, since I think clicking on the tile should actually open the modal */}
          {filing.title}
        </td>
        <td>
          <NoclickSpan>
            <AuthorColumn filing={filing} />
          </NoclickSpan>
        </td>
        {DocketColumn && (
          <td>
            {filing.docket_id && (
              <DocketPill docketId={filing.docket_id} />
            )}
          </td>
        )}

        <td>{filing.item_number}</td>
        <td>
          <FilePill extension={fileExtensionFromText(filing.extension)} />
        </td>
      </tr>
    </>
  );
};
export const FilingTable = ({
  filings,
  scroll,
  DocketColumn,
  PinTableHeader,
}: {
  filings: Filing[];
  scroll?: boolean;
  DocketColumn?: boolean;
  PinTableHeader?: boolean;
}) => {
  // const pinClassName = PinTableHeader ? "table-pin-rows" : "";
  const [modalIsOpen, setModalIsOpen] = useState(false);
  const [modalObjectId, setModalObjectId] = useState("");
  const OpenRowInModal = (id: string) => {
    setModalObjectId(id);
    setModalIsOpen(true);
  };

  return (
    <div>
      <table className={TableStyle}>
        <colgroup>
          <col width="40px" />
          <col width="80px" />
          <col width="200px" />
          {/* authors column */}
          {DocketColumn ? <col width="200px" /> : <col width="300px" />}
          {/* docket id column */}
          {DocketColumn && <col width="50px" />}
          <col width="50px" />
          <col width="50px" />
        </colgroup>
        <thead>
          <tr>
            <th className="text-left sticky top-0 filing-table-date">
              Date Filed
            </th>
            <th className="text-left sticky top-0 filing-table-doc-class">
              Document Class
            </th>
            <th className="text-left sticky top-0">Title</th>
            <th className="text-left sticky top-0">Author</th>
            {DocketColumn && (
              <th className="text-left sticky top-0">Docket ID</th>
            )}
            <th className="text-left sticky top-0">Item Number</th>
            <th className="text-left sticky top-0">Extension</th>
          </tr>
        </thead>
        <tbody>
          {filings.map((filing) => (
            <TableRow
              OpenRow={OpenRowInModal}
              filing={filing}
              {...(DocketColumn !== undefined ? { DocketColumn } : {})}
              key={filing.id}
            />
          ))}
        </tbody>
      </table>
      <DocumentModal objectId={modalObjectId} open={modalIsOpen} setOpen={setModalIsOpen} />
    </div>
  );
};
export default FilingTable;
