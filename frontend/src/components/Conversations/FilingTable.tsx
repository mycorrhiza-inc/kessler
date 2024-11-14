import { useState } from "react";
import Modal from "../styled-components/Modal";
import DocumentModalBody from "../Document/DocumentModalBody";
import { Filing } from "../../lib/types/FilingTypes";
const TableRow = ({ filing }: { filing: Filing }) => {
  const [open, setOpen] = useState(false);
  return (
    <>
      <tr
        className="border-b border-gray-200"
        onClick={() => {
          setOpen((previous) => !previous);
        }}
      >
        <td>{filing.date}</td>
        <td>{filing.title}</td>
        <td>{filing.author}</td>
        <td>{filing.file_class}</td>
        <td>{filing.item_number}</td>
      </tr>
      <Modal open={open} setOpen={setOpen}>
        <DocumentModalBody
          open={open}
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
          ? "max-h-[500px] overflow-y-auto overflow-x-scroll"
          : "max-h-[500px] overflow-y-auto"
      }
    >
      <table className="w-full divide-y divide-gray-200 table">
        <tbody>
          <tr className="border-b border-gray-200">
            <th className="text-left p-2 sticky top-0 bg-white">Date Filed</th>
            <th className="text-left p-2 sticky top-0 bg-white">Title</th>
            <th className="text-left p-2 sticky top-0 bg-white">Author</th>
            <th className="text-left p-2 sticky top-0 bg-white">
              Document Class
            </th>
            <th className="text-left p-2 sticky top-0 bg-white">Item Number</th>
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
