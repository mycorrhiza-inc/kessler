import { Button, Box, useStatStyles } from "@chakra-ui/react";
import { Command, useHotkeys, NavItem } from "@saas-ui/react";
import {
  CommandBar,
  CommandBarDialog,
  CommandBarContent,
  CommandBarInput,
  CommandBarList,
  CommandBarGroup,
  CommandBarItem,
  CommandBarSeparator,
  CommandBarLoading,
  CommandBarEmpty,
} from "@saas-ui/command-bar";
import {
  FiUserCheck,
  FiUser,
  FiCircle,
  FiBarChart,
  FiTag,
  FiSearch,
  FiCalendar,
} from "react-icons/fi";
import { useState } from "react";

interface SearchItem {
  icon: React.ReactNode;
  label: String;
  shortcut: String;
}

export default function SearchDialog() {
  const [searchResults, setResults] = useState<SearchItem[]>([]);
  let q = "";

  const [searchDialogOpen, changeSeachDialog] = useState(false);
  const toggleSearchDialog = () => {
    changeSeachDialog(!searchDialogOpen);
  };
  const [isLoading, setLoading] = useState(false);

  function getResults(value: String) {
    console.log(`searchQuery: ${value}`);
    setResults([]);
  }

  useHotkeys("ctrl+\\", () => {
    toggleSearchDialog();
  });
  useHotkeys("cmd+\\", () => {
    toggleSearchDialog();
  });

  return (
    <>
      <NavItem onClick={toggleSearchDialog} icon={<FiSearch />}>
        Search (ctrl+ \)
      </NavItem>

      <CommandBar
        // TODO: update this in real time
        value={q}
        // onChange={(value) => getResults(value)}
        isOpen={searchDialogOpen}
        onClose={toggleSearchDialog}
        closeOnSelect
      >
        <CommandBarDialog>
          <CommandBarContent>
            <CommandBarInput
              placeholder="search for documents..."
              autoFocus
            />

            <CommandBarList>
              {isLoading && <CommandBarLoading>Hang on…</CommandBarLoading>}

              <CommandBarEmpty>No results found.</CommandBarEmpty>

              {searchResults.map(({ icon, label, shortcut }) => {
                return (
                  <CommandBarItem key={label} value={label}>
                    {icon}
                    {label}
                    <Command ms="auto">{shortcut}</Command>
                  </CommandBarItem>
                );
              })}
            </CommandBarList>
          </CommandBarContent>
        </CommandBarDialog>
      </CommandBar>
    </>
  );
}
