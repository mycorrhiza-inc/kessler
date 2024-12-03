import React, { useState } from "react";
import Modal from "../styled-components/Modal";
import DocumentModalBody from "../Document/DocumentModalBody";
import { Filing } from "../../lib/types/FilingTypes";
import { AuthorInformation } from "@/lib/types/backend_schemas";
import { AuthorInfoPill, TextPill } from "./TextPills";

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
  header,
}: {
  filing: Filing;
  DocketColumn?: boolean;
  header?: boolean;
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
    setOpen((previous) => !previous);
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
            <TextPill
              text={filing.docket_id}
              href={`/proceedings/${filing.docket_id}`}
            />
          </td>
        )}

        <td>{filing.item_number}</td>
      </tr>
      <Modal open={open} setOpen={setOpen}>
        <DocumentModalBody
          open={open}
          title={filing.title}
          objectId={filing.id}
          isPage={false}
        />
      </Modal>
    </>
  );
};
export const FilingTable = ({
  filings,
  scroll,
  DocketColumn,
}: {
  filings: Filing[];
  scroll?: boolean;
  DocketColumn?: boolean;
}) => {
  return (
    <div
      className={
        scroll
          ? "max-h-[1000px] overflow-y-scroll overflow-x-scroll"
          : "overflow-y-auto"
      }
    >
      <table className="w-full divide-y divide-gray-200  border-collaps table table-pin-rows lg:table-fixed md:table-auto sm:table-auto">
        <colgroup>
          <col width="40px" />
          <col width="80px" />
          <col width="200px" />
          {/* authors column */}
          {DocketColumn ? <col width="200px" /> : <col width="300px" />}
          {/* docket id column */}
          {DocketColumn && <col width="50px" />}
          <col width="50px" />
        </colgroup>
        <thead>
          <tr className="border-b border-gray-200">
            <th className="text-left sticky top-0 filing-table-date">
              Date Filed
            </th>
            <th className="text-left sticky top-0 filing-table-doc-class">
              Document Class
            </th>
            <th className="text-left sticky top-0">Title</th>
            <th className="text-left sticky top-0">Author</th>
            {DocketColumn && (
              <th className="text-left sticky top-0">Proceeding ID</th>
            )}
            <th className="text-left sticky top-0">Item Number</th>
          </tr>
        </thead>
        <tbody>
          {filings.map((filing) => (
            <TableRow
              filing={filing}
              {...(DocketColumn !== undefined ? { DocketColumn } : {})}
            />
          ))}
        </tbody>
      </table>
    </div>
  );
};
export default FilingTable;
