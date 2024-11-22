import React, { useState } from "react";
import Modal from "../styled-components/Modal";
import DocumentModalBody from "../Document/DocumentModalBody";
import { Filing } from "../../lib/types/FilingTypes";
import { is } from "date-fns/locale";
import Link from "next/link";
import { AuthorInformation } from "@/lib/types/backend_schemas";
import { AuthorInfoPill, TextPill } from "./TextPills";

const NoclickSpan = ({ children }: { children: React.ReactNode }) => {
  return <span className="noclick p-5">{children}</span>;
};

const TableRow = ({
  filing,
  DocketColumn,
}: {
  filing: Filing;
  DocketColumn?: boolean;
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
          <NoclickSpan>{filing.title}</NoclickSpan>
        </td>
        <td>
          <NoclickSpan>
            {filing.authors_information
              ? filing.authors_information.map(
                  (auth_info: AuthorInformation) => (
                    <AuthorInfoPill author_info={auth_info} />
                  ),
                )
              : filing.author + " Something isnt working"}
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
          ? "max-h-[1000px] overflow-y-auto overflow-x-scroll"
          : "overflow-y-auto"
      }
    >
      <table className="w-full divide-y divide-gray-200 table table-pin-rows">
        <tbody>
          <tr className="border-b border-gray-200">
            <th className="text-left p-2 sticky top-0">Date Filed</th>
            <th className="text-left p-2 sticky top-0">Document Class</th>
            <th className="text-left p-2 sticky top-0">Title</th>
            <th className="text-left p-2 sticky top-0">Author</th>
            {DocketColumn && (
              <th className="text-left p-2 sticky top-0">Proceeding ID</th>
            )}
            <th className="text-left p-2 sticky top-0">Item Number</th>
          </tr>
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
