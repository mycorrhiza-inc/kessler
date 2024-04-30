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
import React, { useState } from "react";
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

import SearchDialog from "./SearchDialog";

export default function Page({ children }: { children: React.ReactNode }) {
  const [isOpen, toggleOpen] = useState(false);
  const [searchModal, changeSearchModal] = useState(false);

  const toggleSearchModal = () => {
    changeSearchModal(!searchModal);
  };

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
      //   navbar={
      //     <Navbar borderBottomWidth="1px" position="sticky" top="0">
      //       <NavbarContent justifyContent="flex-end">
      //         <NavbarItem>
      //           <SearchInput size="sm" />
      //         </NavbarItem>
      //         <NavbarItem padding="10px">
      //           <SignedOut>
      //             <SignInButton />
      //           </SignedOut>
      //           <SignedIn>
      //             <UserButton />
      //           </SignedIn>
      //         </NavbarItem>
      //       </NavbarContent>
      //     </Navbar>
      //   }
      sidebar={
        <Sidebar
          toggleBreakpoint={false}
          variant={isOpen ? "default" : "compact"}
          transition="width"
          transitionDuration="normal"
          width={isOpen ? "280px" : "16"}
          minWidth="auto"
        >
          <SidebarSection>
            <NavItem
              padding="3px"
              display="flex"
              alignItems="center"
              justifyContent="center"
            >
              <SignedOut>
                <SignInButton />
              </SignedOut>
              <SignedIn>
                <UserButton />
              </SignedIn>
            </NavItem>
          </SidebarSection>
          <SidebarSection direction={isOpen ? "row" : "column"}>
            <IconButton
              onClick={() => {
                toggleOpen(!isOpen);
              }}
              variant="ghost"
              size="sm"
              icon={isOpen ? <FiChevronsLeft /> : <FiChevronsRight />}
              aria-label="Toggle Sidebar"
            />
          </SidebarSection>

          <SidebarSection flex="1" overflowY="auto" overflowX="hidden">
            <NavGroup>
              <NavItem href="/" icon={<FiHome />} isActive={pathIs("")}>
                Home
              </NavItem>
              <NavItem
                href="/chat"
                icon={<FiMessageCircle />}
                isActive={pathIs("chat")}
              >
                Chat
              </NavItem>
              <NavItem
                href="/projects"
                icon={<FiFeather />}
                isActive={pathIs("projects")}
              >
                Projects
              </NavItem>
              <NavItem icon={<FiBookmark />} isActive={pathIs("saved")}>
                Saved Documents
              </NavItem>
              <SearchDialog />
            </NavGroup>
          </SidebarSection>
          <SidebarOverlay zIndex="1" />
        </Sidebar>
      }
    >
      <Box as="main" flex="1" py="2" px="4">
        {children}
      </Box>
      <Modal isOpen={searchModal} onClose={toggleSearchModal}>
        <ModalOverlay />
        <ModalContent maxH="1500px" maxW="1500px" overflow="scroll">
          <ModalBody>
            <SearchInput placeholder="Search" />
          </ModalBody>
        </ModalContent>
      </Modal>
    </AppShell>
  );
}
