import { Box, Button, Center, Select, Spinner } from "@chakra-ui/react";
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
  }, [fileUrl, data, page, numResults]);

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
        <Paginator page={page} setPage={setPage} maxPage={maxPage} />
        <br />
        <NumResultsSelector
          numResults={numResults}
          setNumResults={setNumResults}
          setPage={setPage}
        />
      </Box>
    </Box>
  );
};

interface NumResultsSelectorProps {
  numResults: number;
  setNumResults: (num: number) => void;
  setPage: (page: number) => void;
}

const NumResultsSelector: React.FC<NumResultsSelectorProps> = ({
  numResults,
  setNumResults,
  setPage,
}) => {
  return (
    <Select
      value={numResults}
      onChange={(e) => {
        setNumResults(Number(e.target.value));
        setPage(1);
      }}
    >
      <option value={10}>10</option>
      <option value={25}>25</option>
      <option value={50}>50</option>
      <option value={100}>100</option>
      <option value={500}>500</option>
      <option value={65535}>Unlimited</option>
    </Select>
  );
};

interface PaginatorProps {
  page: number;
  setPage: (page: number) => void;
  maxPage: number;
}

const Paginator: React.FC<PaginatorProps> = ({ page, setPage, maxPage }) => {
  const createPageButtons = () => {
    const buttons = [];
    if (maxPage <= 6) {
      for (let i = 1; i <= maxPage; i++) {
        buttons.push(
          <Button
            key={i}
            onClick={() => setPage(i)}
            bg={i === page ? "blue.500" : "gray.200"}
            color={i === page ? "white" : "black"}
          >
            {i}
          </Button>,
        );
      }
    } else {
      buttons.push(
        <Button
          key={1}
          onClick={() => setPage(1)}
          bg={1 === page ? "blue.500" : "gray.200"}
          color={1 === page ? "white" : "black"}
        >
          1
        </Button>,
      );
      if (page > 5) {
        buttons.push(<span key="front-ellipsis">...</span>);
      }
      for (
        let i = Math.max(2, page - 1);
        i <= Math.min(maxPage - 1, page + 1);
        i++
      ) {
        buttons.push(
          <Button
            key={i}
            onClick={() => setPage(i)}
            bg={i === page ? "blue.500" : "gray.200"}
            color={i === page ? "white" : "black"}
          >
            {i}
          </Button>,
        );
      }
      if (maxPage - page > 5) {
        buttons.push(<span key="back-ellipsis">...</span>);
      }
      buttons.push(
        <Button
          key={maxPage}
          onClick={() => setPage(maxPage)}
          bg={maxPage === page ? "blue.500" : "gray.200"}
          color={maxPage === page ? "white" : "black"}
        >
          {maxPage}
        </Button>,
      );
    }
    return buttons;
  };

  return (
    <Box display="flex" alignItems="center">
      <Button
        onClick={() => setPage(Math.max(page - 1, 1))}
        isDisabled={page === 1}
      >
        Previous
      </Button>
      {createPageButtons()}
      <Button
        onClick={() => setPage(Math.min(page + 1, maxPage))}
        isDisabled={page === maxPage}
      >
        Next
      </Button>
    </Box>
  );
};

export default FilePageBrowser;
