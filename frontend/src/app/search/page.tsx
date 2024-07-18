"use client";

import { Container, Box, Input, Flex } from "@chakra-ui/react";
import { ChangeEvent, useEffect, useState } from "react";

const ResultComponent = ({ content }: { content: string }) => {
  return <Container className="searchResult"></Container>;
};
const SearchResultList = ({ ids }: { ids: string[] }) => {
  const [count, setCount] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => setCount((count + 1) % 10), 1000);
    // return () => clearInterval(interval);
  }, []);

  if (ids.length <= 0) {
    return <>Loading {".".repeat(count)}</>;
  } else {
    return (
      <>
        {ids.map((result) => (
          <ResultComponent content={result} />
        ))}
      </>
    );
  }
};

export default function SearchPage() {
  const [queryText, setQuery] = useState("");
  const [searchResults, setSearchResults] = useState<string[]>([]);

  const handleSearchChange = (event: ChangeEvent<HTMLInputElement>) =>
    setQuery(event.target.value);

  const evaluateSearch = () => {};

  return (
    <Box w="100vw" h="100vh" m="0" p="0">
      <Flex direction="column" minWidth="m" alignItems="center">
        <Container className="searchInput" minW="90vw" justifySelf="center">
          <Input value={queryText} onChange={handleSearchChange} />
        </Container>
        <SearchResultList ids={searchResults} />
      </Flex>
    </Box>
  );
}
