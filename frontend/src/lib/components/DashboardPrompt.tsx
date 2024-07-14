"use client";
import {
  Box,
  Container,
  Center,
  Grid,
  GridItem,
  Image,
  Button,
  IconButton,
  Spacer,
  Menu,
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
} from "react-icons/fi";

export function DashboardPrompt() {
  return (
    <Center width="100%" height="100%">
      <Box
        border="solid"
        borderColor="oklch(92.83% 0.001 286.37)"
        width="80%"
        height="80%"
        borderRadius="10px"
        borderWidth="1px"
        margin="40px"
        padding="50px"
        justifySelf="center"
      >
        <Grid
          h="100%"
          templateRows="repeat(4, 1fr)"
          templateColumns="repeat(3, 1fr)"
          gap={3}
        >
          <GridItem colSpan={3} rowSpan={3} bg="tomato"></GridItem>
          <GridItem rowSpan={1} bg="papayawhip">
            Add Data Source
          </GridItem>
          <GridItem rowSpan={1} bg="antiquewhite">
            Get Chatting
          </GridItem>
          <GridItem rowSpan={1} bg="azure">
            Browse Documents
          </GridItem>
        </Grid>
      </Box>
    </Center>
  );
}

export default function DashboardFileBrowser() {
  return (
    <>
      <Center width="100%" height="100%">
        <Box
          border="solid"
          borderColor="oklch(92.83% 0.001 286.37)"
          width="90%"
          borderRadius="10px"
          borderWidth="1px"
          margin="20px"
          justifySelf="center"
        >
          { 
            // NOTE: Refactored file explorer so it broke this implementation, 
            // it seems like this shouldnt be here. so I removed it
            // <FileExplorer />
          }
        </Box>
      </Center>
    </>
  );
}
