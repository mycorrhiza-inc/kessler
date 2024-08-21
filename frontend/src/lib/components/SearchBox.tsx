"use client";
import {
  CircularProgress,
  Input,
  VStack,
  StackDivider,
  Grid,
  GridItem,
} from "@chakra-ui/react";
import React, { ChangeEvent, useEffect, useState } from "react";
import { FiSearch } from "react-icons/fi";

import { defaultLayout } from "./FileTable";
import CustomizeFileTableButton from "./CustomizeFileTableButton";
import FileTable from "./FileTable";
// import SearchDialog from "./SearchDialog";
import { Center } from "@chakra-ui/react";
function SearchBox() {
  const [searching, setSearching] = useState(false);
  const [searchQuery, setQuery] = useState("");
  const [searchResults, setResults] = useState([]);
  const handleQueryChange = async (event: ChangeEvent<HTMLInputElement>) => {
    setResults([]);
    setQuery(event.target.value);
    setSearching(true);
    await getSearchResults();
  };
  const [layout, setLayout] = useState(defaultLayout);
  const getSearchResults = async () => {
    let results = await fetch("/api/search?only_fileobj=true", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ query: searchQuery }),
    })
      .then(async (response) => {
        if (!response.ok) {
          return response.json().then((err) => {
            throw new Error(err.message);
          });
        }
        return await response.json();
      })
      .then((data) => {
        console.log("Success:", data);
        return data;
      })
      .catch((error) => {
        console.error("Error:", error);
        return Symbol("error");
      });

    setSearching(false);
    if (results == Symbol("Error")) {
      console.log("no search results");
      return;
    }
    setResults(results);
  };
  return (
    <VStack
      divider={<StackDivider borderColor="gray.200" />}
      spacing={4}
      align="stretch"
      overflowY="scroll"
      minH="10vh"
    >
      <Grid templateColumns="repeat(20, 5%)">
        <GridItem colSpan={1}>
          <Center justifySelf="center">
            <FiSearch />
          </Center>
        </GridItem>
        <GridItem colStart={2} colEnd={20} h="10">
          <Input
            value={searchQuery}
            placeholder="search"
            size="lg"
            onChange={handleQueryChange}
          />
        </GridItem>
      </Grid>
      {searching && <CircularProgress isIndeterminate color="green.300" />}
      {searchResults.length > 0 && (
        <>
          <FileTable files={searchResults} layout={layout}></FileTable>
          <CustomizeFileTableButton layout={layout} setLayout={setLayout} />
        </>
      )}
    </VStack>
  );
}

export default SearchBox;
