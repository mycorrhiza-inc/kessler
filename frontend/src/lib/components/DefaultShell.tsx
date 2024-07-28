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
import SearchBox from "./SearchBox";

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
            <SearchBox></SearchBox>
          </ModalBody>
        </ModalContent>
      </Modal>
    </AppShell>
  );
}
