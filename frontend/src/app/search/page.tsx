'use client';

import { Container, Box, Input, Flex, CircularProgress, Spacer, Center } from '@chakra-ui/react';
import { ChangeEvent, useEffect, useState } from 'react';


const ResultComponent = ({ item, index }: { item: any, index: number }) => {
  return <Container className='searchResult'>
    <div key={index} className="data-item">
      <p>ID: {item.id}</p>
      <p>Score: {item.score}</p>
      <p>Author: {item.metadata.author}</p>
      <p>Title: {item.metadata.title}</p>
    </div>
  </Container>
}

const SearchResultList = ({ results, searching }: { results: any[], searching: boolean }) => {
  return <>
    {searching && (
      <CircularProgress isIndeterminate color="green.300" />
    )}
    {results.length > 0 &&
      results.map((item: any, index) => (
        <>
          <ResultComponent item={item} index={index} />
        </>
      ))}

  </>
}

export default function SearchPage() {
  const [searching, setSearching] = useState(false);
  const [searchQuery, setQuery] = useState("");
  const [searchResults, setResults] = useState([]);
  const handleQueryChange = async (event: ChangeEvent<HTMLInputElement>) => {
    setResults([]);
    setQuery(event.target.value);
  };

  // search stuff
  const getSearchResults = async () => {
    let results = await fetch("/api/search", {
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

  const handleEnter = async (event) => {
    if (event.key === 'Enter') {
      await getSearchResults()
    }
  }

  return <Box w='100vw' h='100vh' m='0' p='0'>
    <Flex direction="column" alignItems='center'>
      <Center>
        <h1>Kessler</h1>
      </Center>
      <Spacer />
      <Center
        className='searchInput'
        minW='90vw'
        justifySelf='center'
        justifyContent='center'
      >
        <Input
          maxW="50vw"
          minW="500px"
          value={searchQuery}
          onChange={handleQueryChange}
          onKeyDown={handleEnter}
        />
      </Center>
      <Spacer />
      <Center>
        <SearchResultList results={searchResults} searching={searching} />
      </Center>
    </Flex>
  </Box>
}




