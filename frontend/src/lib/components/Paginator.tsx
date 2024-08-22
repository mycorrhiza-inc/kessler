import { Box, Button, Center, Select, Spinner, Text } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import FileTable from "./FileTable";
import CustomizeFileTableButton from "./CustomizeFileTableButton";
import { FileType } from "../interfaces";

import { defaultLayout } from "./FileTable";

interface PaginatorProps {
  numResults: number;
  setNumResults: (num: number) => void;
  setPage: (page: number) => void;
  page: number;
  maxPage: number;
}

const Paginator: React.FC<PaginatorProps> = ({
  numResults,
  setNumResults,
  setPage,
  page,
  maxPage,
}) => {
  return (
    <Box>
      <Box mt={4} display="flex" justifyContent="space-between">
        <PaginatorButtons page={page} setPage={setPage} maxPage={maxPage} />
      </Box>
      <Box>
        <Text>
          <br />
          Select Number of Results
        </Text>
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
      <option value={1}>1</option>
      <option value={10}>10</option>
      <option value={25}>25</option>
      <option value={50}>50</option>
      <option value={100}>100</option>
      <option value={500}>500</option>
      <option value={65535}>Unlimited</option>
    </Select>
  );
};

interface PaginatorButtonProps {
  page: number;
  setPage: (page: number) => void;
  maxPage: number;
}

const PaginatorButtons: React.FC<PaginatorButtonProps> = ({
  page,
  setPage,
  maxPage,
}) => {
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

export default Paginator;
