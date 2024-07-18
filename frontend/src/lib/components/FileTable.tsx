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
import { IoMdCheckmarkCircleOutline } from "react-icons/io";
import { ImCross } from "react-icons/im";
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
            <Th width="70%">Filename</Th>
            <Th width="20%">Source</Th>
            <Th width="6%">View</Th>
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
                <Td>{truncateString(file.data.name)}</Td>
                <Td>{truncateString(file.data.source)}</Td>
                <Td>
                  <DocumentViewer document_uuid={file.data.id} />
                </Td>
                <Td>
                  {file.data.stage == "completed" ? (
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
