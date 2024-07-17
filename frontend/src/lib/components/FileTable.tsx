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
import { useState } from "react";
import { FileType } from "../interfaces/file";
import { LoadingSpinner } from "@saas-ui/react";
import DocumentViewer from "./DocumentViewer";

interface RowData {
  selected: boolean;
  data: FileType;
}

interface FileTableProps {
  files: FileType[];
}

const FileTable: React.FC<FileTableProps> = ({ files }) => {
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

  return (
    <TableContainer>
      <Table>
        <Thead>
          <Tr>
            <Th width="2%">Select</Th>
            <Th width="96%">Filename</Th>
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
                <Td>{file.data.name}</Td>
                <Td>
                  <DocumentViewer document_uuid={file.data.id} />
                </Td>
              </Tr>
            ))}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default FileTable;
