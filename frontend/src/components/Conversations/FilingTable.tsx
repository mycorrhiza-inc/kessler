import { useState } from "react";
import Modal from "../styled-components/Modal";
import DocumentModalBody from "../Document/DocumentModalBody";
import { Filing } from "../../lib/types/FilingTypes";
import { QueryDataFile, QueryFilterFields } from "@/lib/filters";
import { getSearchResults } from "@/lib/requests/search";
import { memo } from "react";
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
          objectId={filing.id}
          overridePDFUrl={filing.url}
        />
      </Modal>
    </>
  );
};
export const FilingTable = ({ filings, scroll }: { filings: Filing[], scroll?: boolean }) => {
  return (
    <div className={"min-h-[500px] overflow-y-auto" + (scroll ? "max-h-[500px] overflow-x-scroll" : "")}>
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

const FilingTableQuery = memo(
  async ({ queryData }: { queryData: QueryDataFile }) => {
    try {
      const filings: Filing[] = await getSearchResults(queryData);
      if (filings == undefined) {
        return <p>Filings returned from server is undefined.</p>;
      }
      return <FilingTable filings={filings} />;
    } catch (error) {
      return (
        <p>
          Encountered an Error getting files from the server. <br />
          {String(error)}
        </p>
      );
    }
  },
);

export default FilingTableQuery;
