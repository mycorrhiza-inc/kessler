"use client";
import {
  ClerkProvider,
  SignInButton,
  SignedIn,
  SignedOut,
  UserButton,
} from "@clerk/nextjs";
import { Box, Image, IconButton, Spacer, Menu } from "@chakra-ui/react";
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
} from "react-icons/fi";

import { usePathname } from "next/navigation";

export default function Page({ children }: { children: React.ReactNode }) {
  const [isOpen, toggleOpen] = useState(false);
  const pathname = usePathname();

  function pathIs(name: string) {
    let leafName = pathname.split("/")[pathname.length - 1];
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
          <SidebarSection direction={isOpen ? "row" : "column"}>
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
            <Spacer />
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
              <NavItem icon={<FiHome />} isActive={pathIs("home")}>
                Home
              </NavItem>
              <NavItem icon={<FiBookmark />} isActive={pathIs("saved")}>
                Saved Documents
              </NavItem>
              <NavItem icon={<FiSearch />} isActive={pathIs("search")}>
                Search
              </NavItem>
            </NavGroup>
          </SidebarSection>
          <SidebarOverlay zIndex="1" />
        </Sidebar>
      }
    >
      <Box as="main" flex="1" py="2" px="4">
        {children}
      </Box>
    </AppShell>
  );
}
