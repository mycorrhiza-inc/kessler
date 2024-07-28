import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  TableContainer,
  Checkbox,
  Box,
  Center,
} from "@chakra-ui/react";
import { useState, useEffect } from "react";
import DocumentViewer from "./DocumentViewer";
import { IoMdCheckmarkCircleOutline } from "react-icons/io";
import { ImCross } from "react-icons/im";
import { FileType } from "../interfaces/file";

interface RowData {
  selected: boolean;
  data: FileType;
}

interface Layout {
  columns: {
    key: string;
    label: string;
    width: string;
    enabled: boolean;
  }[];
  showExtraFeatures: boolean;
}

interface FileTableProps {
  files: FileType[];
  layout: Layout;
}

export const defaultLayout: Layout = {
  columns: [
    { key: "name", label: "Filename", width: "80%", enabled: true },
    { key: "source", label: "Source", width: "20%", enabled: true },
    { key: "author", label: "Author", width: "20%", enabled: true },
  ],
  showExtraFeatures: true,
};

const FileTable: React.FC<FileTableProps> = ({ files, layout }) => {
  const [fileState, setFileState] = useState<RowData[]>(
    files.map((file) => ({ selected: false, data: file })),
  );

  const updateSelected = (id: string) => {
    setFileState((prevState) =>
      prevState.map((file) =>
        file.data.id === id ? { ...file, selected: !file.selected } : file,
      ),
    );
  };

  function truncateString(str: string) {
    const length = 60;
    return str.length < length ? str : str.slice(0, length - 3) + "...";
  }
  function getFieldFromFile(key: string, file: FileType): string {
    // Please shut up, I know what Im doing
    // @ts-ignore
    var result = file[key];
    if (result == undefined) {
      // @ts-ignore
      result = file.mdata[key];
    }
    if (result != undefined) {
      return String(result);
    }
    return "Unknown";
  }
  const layoutFiltered: Layout = {
    ...layout,
    columns: layout.columns.filter((column) => column.enabled),
  };

  return (
    <TableContainer>
      <Table>
        <Thead>
          <Tr>
            <Th width="2%">Select</Th>
            {layoutFiltered.columns.map((col) => (
              <Th key={col.key} width={col.width}>
                {col.label}
              </Th>
            ))}
            {layoutFiltered.showExtraFeatures && (
              <>
                <Th width="6%">View</Th>
                <Th width="2%">Status</Th>{" "}
              </>
            )}
          </Tr>
        </Thead>
        <Tbody>
          {fileState.map((file) => (
            <Tr key={file.data.id}>
              <Td>
                <Box>
                  <Checkbox
                    isChecked={file.selected}
                    onChange={() => updateSelected(file.data.id)}
                  />
                </Box>
              </Td>
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
                  <Td>
                    {file.data.stage === "completed" ? (
                      <IoMdCheckmarkCircleOutline />
                    ) : (
                      <ImCross />
                    )}
                  </Td>
                </>
              )}
            </Tr>
          ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default FileTable;
