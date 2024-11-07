import { useState } from "react";
import Modal from "../styled-components/Modal";
import DocumentModalBody from "../Document/DocumentModalBody";
import { Filing } from "./conversationTypes";
import { QueryDataFile, QueryFilterFields } from "@/lib/filters";
import searchResultsGet from "./searchResultGet";

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
        <td>{filing.source}</td>
        <td>{filing.item_number}</td>
        <td>
          <a href={filing.url}>View</a>
        </td>
      </tr>
      <Modal open={open} setOpen={setOpen}>
        <DocumentModalBody
          open={open}
          objectId={filing.uuid}
          overridePDFUrl={filing.url}
        />
      </Modal>
    </>
  );
};
const FilingTable = ({ filings }: { filings: Filing[] }) => {
  return (
    <div className="overflow-x-scroll max-h-[500px] overflow-y-auto">
      <table className="w-full divide-y divide-gray-200 table">
        <tbody>
          <tr className="border-b border-gray-200">
            <th className="text-left p-2 sticky top-0 bg-white">Date Filed</th>
            <th className="text-left p-2 sticky top-0 bg-white">Title</th>
            <th className="text-left p-2 sticky top-0 bg-white">Author</th>
            <th className="text-left p-2 sticky top-0 bg-white">Source</th>
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

const FilingTableQuery = async ({
  queryData,
}: {
  queryData: QueryDataFile;
}) => {
  const filings = await searchResultsGet(queryData);
  return <FilingTable filings={filings} />;
};

export default FilingTableQuery;
