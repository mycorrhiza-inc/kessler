import {
  ClerkProvider,
  SignInButton,
  SignedIn,
  SignedOut,
  UserButton,
} from "@clerk/nextjs";
import { Box } from "@chakra-ui/react";
import {
  AppShell,
  Sidebar,
  SidebarSection,
  NavItem,
  Navbar,
  NavbarContent,
  NavbarItem,
  SearchInput,
} from "@saas-ui/react";

export default function Page() {
  return (
    <AppShell
      variant="static"
      minH="100vh"
      navbar={
        <Navbar borderBottomWidth="1px" position="sticky" top="0">
          <NavbarContent justifyContent="flex-end">
            <NavbarItem>
              <SearchInput size="sm" />
            </NavbarItem>
            <NavbarItem padding="10px">
              <SignedOut>
                <SignInButton />
              </SignedOut>
              <SignedIn>
                <UserButton />
              </SignedIn>
            </NavbarItem>
          </NavbarContent>
        </Navbar>
      }
      sidebar={
        <Sidebar position="sticky" top="56px" toggleBreakpoint="sm">
          <SidebarSection>
            <NavItem>Home</NavItem>
            <NavItem>Settings</NavItem>
          </SidebarSection>
        </Sidebar>
      }
    >
      <Box as="main" flex="1" py="2" px="4">
        Your application content
      </Box>
    </AppShell>
  );
}
