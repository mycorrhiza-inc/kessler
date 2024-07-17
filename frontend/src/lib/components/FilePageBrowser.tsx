import {
  Box,
  Button,
  Center,
  Select,
  Spinner,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import FileTable from "./FileTable";
import { FileType } from "../interfaces/file";

interface FilePageBrowserProps {
  fileUrl: string;
  data: any;
}

const FilePageBrowser: React.FC<FilePageBrowserProps> = ({ fileUrl, data }) => {
  const [files, setFiles] = useState<FileType[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [maxPage, setMaxPage] = useState(1);
  const [numResults, setNumResults] = useState(10);

  const fetchFiles = async () => {
    setLoading(true);
    const response = await fetch(
      `${fileUrl}?num_results=${numResults}&page=${page}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      },
    );
    const result = await response.json();
    console.log(result);
    setFiles(result[0]);
    setMaxPage(result[1]);
    setLoading(false);
  };

  useEffect(() => {
    fetchFiles();
  }, [page, numResults]);

  return (
    <Box>
      {loading ? (
        <Center>
          <Spinner />
        </Center>
      ) : (
        <FileTable files={files} />
      )}
      <Box mt={4} display="flex" justifyContent="space-between">
        <Button
          onClick={() => setPage((prev) => Math.max(prev - 1, 1))}
          isDisabled={page === 1}
        >
          Previous
        </Button>
        <Select
          value={numResults}
          onChange={(e) => {
            setNumResults(Number(e.target.value));
            setPage(1);
          }}
        >
          <option value={10}>10</option>
          <option value={20}>20</option>
          <option value={50}>50</option>
        </Select>
        <Button onClick={() => setPage((prev) => prev + 1)}>Next</Button>
      </Box>
    </Box>
  );
};

export default FilePageBrowser;
