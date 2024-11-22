import React, { useState } from "react";
import Modal from "../styled-components/Modal";
import DocumentModalBody from "../Document/DocumentModalBody";
import { Filing } from "../../lib/types/FilingTypes";
import { is } from "date-fns/locale";
import Link from "next/link";
import { AuthorInformation } from "@/lib/types/backend_schemas";

const oklchSubdivide = (colorNum: number, divisions?: number) => {
  const defaultDivisions = divisions || 15;
  const hue = (colorNum % defaultDivisions) * (360 / defaultDivisions);
  return `oklch(73% 0.123 ${hue})`;
};

const subdivide15 = [
  "oklch(73% 0.123 0)",
  "oklch(73% 0.123 30)",
  "oklch(73% 0.123 60)",
  "oklch(73% 0.123 90)",
  "oklch(73% 0.123 120)",
  "oklch(73% 0.123 150)",
  "oklch(73% 0.123 180)",
  "oklch(73% 0.123 210)",
  "oklch(73% 0.123 240)",
  "oklch(73% 0.123 270)",
  "oklch(73% 0.123 300)",
  "oklch(73% 0.123 330)",
];

type FileColor = {
  pdf: string;
  doc: string;
  xlsx: string;
};

const fileTypeColor = {
  pdf: "oklch(65.55% 0.133 0)",
  doc: "oklch(60.55% 0.13 240)",
  xlsx: "oklch(75.55% 0.133 140)",
};

const IsFiletypeColor = (key: string): key is keyof FileColor => {
  return key in fileTypeColor;
};

const TextPill = ({
  text,
  href,
  seed,
}: {
  text?: string;
  href?: string;
  seed?: string;
}) => {
  const textDefined: string = text || "Unknown";
  const actualSeed = seed || textDefined;
  var pillColor = "";
  if (IsFiletypeColor(textDefined)) {
    pillColor = fileTypeColor[textDefined];
  } else {
    const pillInteger =
      Math.abs(
        actualSeed
          .split("")
          .reduce((acc, char) => (acc * 3 + 2 * char.charCodeAt(0)) % 27, 0),
      ) % 9;
    pillColor = oklchSubdivide(pillInteger, 15);
  }
  // btn-[${pillColor}]
  if (href) {
    return (
      <Link
        style={{ backgroundColor: pillColor }}
        className={`btn btn-xs text-black`}
        href={href}
      >
        {text}
      </Link>
    );
  }
  return (
    <button
      style={{ backgroundColor: pillColor }}
      className={`btn btn-xs no-animation text-black mt-2 mb-2`}
    >
      {text}
    </button>
  );
};

const NoclickSpan = ({ children }: { children: React.ReactNode }) => {
  return <span className="noclick p-5">{children}</span>;
}

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
    const element = (event.target as HTMLElement);
    const tagName = element.tagName; // Ensure type safety
    const className = element.className;
    console.log("element", element);
    console.log('tagName', tagName);
    console.log('className', className);

    const includesNoClick = className.includes('noclick');

    if (tagName === 'SPAN' || tagName === 'BUTTON' || includesNoClick) {
      console.log('Text was clicked');
      // Prevent the event from propagating further if needed
      event.stopPropagation();
      return;
    }

    // If it's not the specific text, allow the row click to proceed
    console.log('Row was clicked');
    setOpen((previous) => !previous)
  };

  return (
    <>
      <tr
        className="border-b border-base-300 hover:bg-base-200 transition duration-500 ease-out"
        onClick={handleRowClick}
      >
        <td><NoclickSpan>{filing.date}</NoclickSpan></td>
        <td>
          <NoclickSpan><TextPill text={filing.file_class} /></NoclickSpan>
        </td>
        <td><NoclickSpan>{filing.title}</NoclickSpan></td>
        <td>
          <NoclickSpan>
          {filing.authors_information
            ? filing.authors_information.map((auth_info: AuthorInformation) => (
              <TextPill
                text={auth_info.author_name}
                seed={auth_info.author_id}
                href={`/orgs/${auth_info.author_id}`}
              />
            ))
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
          overridePDFUrl={filing.url}
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
