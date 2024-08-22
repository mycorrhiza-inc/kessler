import { Box, Button, Center, Select, Spinner, Text } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import FileTable from "./FileTable";
import CustomizeFileTableButton from "./CustomizeFileTableButton";
import { FileType } from "../interfaces";

import { defaultLayout } from "./FileTable";
import Paginator from "./Paginator";

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
  const [layout, setLayout] = useState(defaultLayout);
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
      <CustomizeFileTableButton layout={layout} setLayout={setLayout} />
      {loading ? (
        <Center>
          <Spinner />
        </Center>
      ) : (
        <FileTable files={files} layout={layout} />
      )}
      <Paginator
        page={page}
        setPage={setPage}
        maxPage={maxPage}
        numResults={numResults}
        setNumResults={setNumResults}
      />
    </Box>
  );
};

export default FilePageBrowser;
