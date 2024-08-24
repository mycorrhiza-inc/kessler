import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
  Box,
  Text,
} from "@chakra-ui/react";
import { useState } from "react";
import DocumentViewer from "./DocumentViewer";
import { FileType } from "../interfaces";

import { TableLayout } from "../interfaces";
// Define color variables

interface RowData {
  selected: boolean;
  data: FileType;
}

interface FileTableProps {
  files: FileType[];
  layout: TableLayout;
}

export const defaultLayout: TableLayout = {
  columns: [
    { key: "name", label: "Filename", width: "60%", enabled: true },
    { key: "source", label: "Source", width: "20%", enabled: true },
    { key: "author", label: "Author", width: "20%", enabled: false },
    { key: "docket_id", label: "Docket ID", width: "20%", enabled: false },
  ],
  showExtraFeatures: true,
  showDisplayText: true,
};

const FileTable: React.FC<FileTableProps> = ({ files, layout }) => {
  const [tableState, setTableState] = useState<RowData[]>(
    files.map((file) => ({ selected: false, data: file })),
  );

  const toggleSelected = (id: string) => {
    setTableState((prevState) =>
      prevState.map((file) =>
        file.data.id === id ? { ...file, selected: !file.selected } : file,
      ),
    );
  };

  const handleKeydown = (e: React.KeyboardEvent, id: string) => {
    if (e.shiftKey && e.type === "click") {
      toggleSelected(id);
    }
  };

  const updateSelected = (id: string) => {
    setTableState((prevState) =>
      prevState.map((file) =>
        file.data.id === id ? { ...file, selected: !file.selected } : file,
      ),
    );
  };
  function truncateString(str: string, length: number = 60) {
    return str.length < length ? str : str.slice(0, length - 3) + "...";
  }
  function getFieldFromFile(key: string, file: FileType): string {
    // Please shut up, I know what I'm doing
    // @ts-ignore
    let result = file[key];
    if (result === undefined) {
      // @ts-ignore
      result = file.mdata[key];
    }
    return result !== undefined ? String(result) : "Unknown";
  }

  const layoutFiltered: TableLayout = {
    ...layout,
    columns: layout.columns.filter((column) => column.enabled),
  };

  return (
    <TableContainer>
      <Table>
        <Thead>
          <Tr>
            {layoutFiltered.columns.map((col) => (
              <Th key={col.key} width={col.width}>
                {col.label}
              </Th>
            ))}
            {layoutFiltered.showExtraFeatures && (
              <>
                <Th width="6%">View</Th>
              </>
            )}
          </Tr>
        </Thead>
        <Tbody>
          {tableState.map((file) => (
            <>
              <Tr key={file.data.id}>
                {layoutFiltered.columns.map((col) => (
                  <Td key={col.key}>
                    {truncateString(getFieldFromFile(col.key, file.data))}
                  </Td>
                ))}
                {layoutFiltered.showExtraFeatures && (
                  <>
                    <Td>
                      <DocumentViewer document_object={file.data} />
                    </Td>
                  </>
                )}
              </Tr>
              {layoutFiltered.showDisplayText && (
                <Tr>
                  <Td colSpan={layoutFiltered.columns.length + 2}>
                    <Box w="100%">
                      {truncateString(file.data.display_text, 150)}
                    </Box>
                  </Td>
                </Tr>
              )}

              <Tr>
                <Td colSpan={layoutFiltered.columns.length + 2}>
                  <hr
                    style={{
                      background: "black",
                      // color: "lime",
                      // borderColor: "lime",
                      height: "3px",
                    }}
                  />
                </Td>
              </Tr>
            </>
          ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default FileTable;

// Figure out way to get backgrounds working with a div
// <div
//   onMouseOver={(e) => {
//     e.currentTarget.style.backgroundColor = highlightColor;
//   }}
//   onMouseOut={(e) => {
//     e.currentTarget.style.backgroundColor = file.selected
//       ? selectedColor
//       : "";
//   }}
//   onClick={(e) => handleKeydown(e as any, file.data.id)}
//   backgroundColor={file.selected ? selectedColor : ""}
// >
// <Td>
//   {file.data.stage === "completed" ? (
//     <IoMdCheckmarkCircleOutline />
//   ) : (
//     <ImCross />
//   )}
// </Td>
