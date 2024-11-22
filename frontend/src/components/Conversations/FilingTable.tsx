import { useState } from "react";
import Modal from "../styled-components/Modal";
import DocumentModalBody from "../Document/DocumentModalBody";
import { Filing } from "../../lib/types/FilingTypes";

const pillColors = [
  "oklch(73% 0.123 0)",
  "oklch(73% 0.123 40)",
  "oklch(73% 0.123 80)",
  "oklch(73% 0.123 120)",
  "oklch(73% 0.123 160)",
  "oklch(73% 0.123 200)",
  "oklch(73% 0.123 240)",
  "oklch(73% 0.123 280)",
  "oklch(73% 0.123 320)",
];

const fileTypeColor = {
  pdf: "oklch(65.55% 0.133 0)",
  doc: "oklch(60.55% 0.13 240)",
  xlsx: "oklch(75.55% 0.133 140)",
};

const FileTypePill = ({ file_class }: { file_class?: string }) => {
  const fileClassDefined = file_class || "Unknown";
  const pillInteger =
    Math.abs(
      fileClassDefined
        .split("")
        .reduce((acc, char) => (acc * 3 + 2 * char.charCodeAt(0)) % 27, 0),
    ) % 9;
  const pillColor = pillColors[pillInteger];
  // btn-[${pillColor}]
  return (
    <button
      style={{ backgroundColor: pillColor }}
      className={`btn btn-xs no-animation text-black`}
    >
      {file_class}
    </button>
  );
};

const TableRow = ({ filing }: { filing: Filing }) => {
  const [open, setOpen] = useState(false);
  return (
    <>
      <tr
        className="border-b border-base-300 hover:bg-base-200 transition duration-500 ease-out"
        onDoubleClick={() => {
          setOpen((previous) => !previous);
        }}
      >
        <td>{filing.date}</td>
        <td>
          <FileTypePill file_class={filing.file_class} />
        </td>
        <td>{filing.title}</td>
        <td>{filing.author}</td>
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
}: {
  filings: Filing[];
  scroll?: boolean;
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
            <th className="text-left p-2 sticky top-0">Item Number</th>
          </tr>
          {filings.map((filing) => (
            <TableRow filing={filing} />
          ))}
        </tbody>
      </table>
    </div>
  );
};
export default FilingTable;
