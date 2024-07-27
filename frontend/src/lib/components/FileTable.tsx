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

interface RowData {
  selected: boolean;
  data: FileType;
}

interface Layout {
  columns: {
    key: string;
    label: string;
    width: string;
    truncate: boolean;
  }[];
  showExtraFeatures: boolean;
}

interface FileTableProps {
  files: FileType[];
  layout: Layout;
}

const FileTable: React.FC<FileTableProps> = ({ files, layout }) => {
  const [fileState, setFileState] = useState<RowData[]>(
    files.map((file) => ({ selected: false, data: file })),
  );
  const [loading, setLoading] = useState(false);

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

  return (
    <TableContainer>
      <Table>
        <Thead>
          <Tr>
            <Th width="2%">Select</Th>
            {layout.columns.map((col) => (
              <Th key={col.key} width={col.width}>
                {col.label}
              </Th>
            ))}
            {layout.showExtraFeatures && <Th width="6%">View</Th>}
            <Th width="2%">Status</Th>
          </Tr>
        </Thead>
        <Tbody>
          {loading && (
            <Tr>
              <Td></Td>
              <Td>
                <Center>
                  <LoadingSpinner />
                </Center>
              </Td>
            </Tr>
          )}
          {!loading &&
            fileState.map((file) => (
              <Tr key={file.data.id}>
                <Td>
                  <Box>
                    <Checkbox
                      isChecked={file.selected}
                      onChange={() => updateSelected(file.data.id)}
                    />
                  </Box>
                </Td>
                {layout.columns.map((col) =>
                  col.key in file.data ? (
                    <Td key={col.key}>
                      {col.truncate
                        ? truncateString(String(file.data[col.key]))
                        : String(file.data[col.key])}
                    </Td>
                  ) : (
                    <Td key={col.key}>
                      {col.truncate ? "Unknown" : "Unknown"}
                    </Td>
                  ),
                )}
                {layout.showExtraFeatures && (
                  <Td>
                    <DocumentViewer document_object={file.data} />
                  </Td>
                )}
                <Td>
                  {file.data.stage === "completed" ? (
                    <IoMdCheckmarkCircleOutline />
                  ) : (
                    <ImCross />
                  )}
                </Td>
              </Tr>
            ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default FileTable;
