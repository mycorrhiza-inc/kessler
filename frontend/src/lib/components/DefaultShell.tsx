"use client";
import {
  ClerkProvider,
  SignInButton,
  SignedIn,
  SignedOut,
  UserButton,
} from "@clerk/nextjs";
import {
  Box,
  Image,
  Button,
  IconButton,
  Spacer,
  Menu,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  useToast,
  useDisclosure,
  CircularProgress,
  Input,
  VStack,
  StackDivider,
  Grid,
  GridItem,
} from "@chakra-ui/react";
import {
  AppShell,
  Sidebar,
  SidebarOverlay,
  SidebarSection,
  SidebarToggleButton,
  NavItem,
  Navbar,
  NavbarContent,
  NavGroup,
  NavbarItem,
  SearchInput,
} from "@saas-ui/react";
import React, { ChangeEvent, useEffect, useState } from "react";
import { Node } from "reactflow";
import {
  FiHome,
  FiUsers,
  FiSettings,
  FiBookmark,
  FiStar,
  FiSearch,
  FiChevronsLeft,
  FiChevronsRight,
  FiMessageCircle,
  FiLayers,
  FiFeather,
} from "react-icons/fi";

import { usePathname } from "next/navigation";
import ColorModeToggle from "./ColorModeToggle";

// import SearchDialog from "./SearchDialog";
import { Center } from "@chakra-ui/react";

export default function Page({ children }: { children: React.ReactNode }) {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [sidebarIsOpen, toggleSidebar] = useState(true);
  const pathname = usePathname();

  function pathIs(name: string) {
    let leafName = pathname.split("/")[-1];
    if (leafName == name) {
      return true;
    }
    return false;
  }

  const toast = useToast();
  const [searching, setSearching] = useState(false);
  const [searchQuery, setQuery] = useState("");
  const [searchResults, setResults] = useState([]);
  const handleQueryChange = async (event: ChangeEvent<HTMLInputElement>) => {
    setResults([]);
    setQuery(event.target.value);
    setSearching(true);
    await getSearchResults();
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

  return (
    <AppShell
      variant="static"
      minH="100vh"
      sidebar={
        <Sidebar
          toggleBreakpoint={false}
          variant={sidebarIsOpen ? "default" : "compact"}
          transition="width"
          transitionDuration="normal"
          width={sidebarIsOpen ? "280px" : "16"}
          minWidth="auto"
        >
          <SidebarSection>
            <NavItem
              padding="3px"
              display="flex"
              alignItems="center"
              justifyContent="center"
            >
              {/* TODO: enable this when users work */}
              {/* <SignedOut>
                <SignInButton />
              </SignedOut>
              <SignedIn>
                <UserButton />
              </SignedIn> */}
            </NavItem>
          </SidebarSection>
          <SidebarSection direction={sidebarIsOpen ? "row" : "column"}>
            <IconButton
              onClick={() => {
                toggleSidebar(!sidebarIsOpen);
              }}
              variant="ghost"
              size="sm"
              icon={sidebarIsOpen ? <FiChevronsLeft /> : <FiChevronsRight />}
              aria-label="Toggle Sidebar"
            />
          </SidebarSection>

          <SidebarSection flex="1" overflowY="auto" overflowX="hidden">
            <NavGroup>
              <NavItem
                href="/"
                icon={<FiMessageCircle />}
                isActive={pathIs("")}
              >
                Chat with Documents
              </NavItem>
              <NavItem
                href="/documents"
                icon={<FiFeather />}
                isActive={pathIs("documents")}
              >
                Modify Document Database
              </NavItem>
              <NavItem
                href="/basic-chat"
                icon={<FiMessageCircle />}
                isActive={pathIs("basic-chat")}
              >
                Basic LLM Chat
              </NavItem>
              <NavItem icon={<FiSearch />} onClick={onOpen}>
                Search
              </NavItem>
            </NavGroup>
          </SidebarSection>
          <SidebarOverlay zIndex="1" />
          <ColorModeToggle />
        </Sidebar>
      }
    >
      <Box as="main" flex="1" py="2" px="4">
        {children}
      </Box>
      <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent maxH="1500px" maxW="1500px" overflow="scroll">
          <ModalBody>
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
              {searching && (
                <CircularProgress isIndeterminate color="green.300" />
              )}
              {(searchResults.length > 0) && searchResults.map((item, index) => (
                <div key={index} className="data-item">
                    <p>ID: {item.id}</p>
                    <p>Score: {item.score}</p>
                    {/* Render other fields as needed */}
                    <p>Author: {item.metadata.author}</p>
                    <p>Title: {item.metadata.title}</p>
                </div>
            ))}
            </VStack>
          </ModalBody>
        </ModalContent>
      </Modal>
    </AppShell>
  );
}
