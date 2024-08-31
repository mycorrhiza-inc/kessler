import { Dispatch, SetStateAction, useState, useEffect, useRef } from "react";
import { Input, Button, Grid, Stack, Divider, Box } from "@mui/joy";
import { motion } from "framer-motion";

interface SearchBoxProps {
  handleSearch: () => Promise<void>;
  searchQuery: string;
  setSearchQuery: Dispatch<SetStateAction<string>>;
  inSearchSession: boolean;
}

const SearchIcon = () => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      className="icon icon-tabler icons-tabler-outline icon-tabler-search"
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
      <path d="M10 10m-7 0a7 7 0 1 0 14 0a7 7 0 1 0 -14 0" />
      <path d="M21 21l-6 -6" />
    </svg>
  );
};

const CommandIcon = () => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      className="icon icon-tabler icons-tabler-outline icon-tabler-command"
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
      <path d="M7 9a2 2 0 1 1 2 -2v10a2 2 0 1 1 -2 -2h10a2 2 0 1 1 -2 2v-10a2 2 0 1 1 2 2h-10" />
    </svg>
  );
};

const ChatIcon = () => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      className="icon icon-tabler icons-tabler-outline icon-tabler-message-plus"
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
      <path d="M8 9h8" />
      <path d="M8 13h6" />
      <path d="M12.01 18.594l-4.01 2.406v-3h-2a3 3 0 0 1 -3 -3v-8a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v5.5" />
      <path d="M16 19h6" />
      <path d="M19 16v6" />
    </svg>
  );
};

const SearchBox = ({
  handleSearch,
  searchQuery,
  setSearchQuery,
  inSearchSession,
}: SearchBoxProps) => {
  return (
    <Box>
      <Stack>
        {/* Your UI elements here */}
        <Stack
          direction="row"
          spacing={2}
          className="flex items-center color-white justify-center"
        >
          <Input
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            variant="outlined"
            className="w-full"
            placeholder="Search"
            style={{ backgroundColor: "white" }}
          />
          <Button
            style={{
              backgroundColor: "--brand-yellow-rgb",
              color: "black",
              border: "2px solid grey",
            }}
            className="max-w-60"
            onClick={handleSearch}
          >
            <SearchIcon />
          </Button>
        </Stack>
        <div className="flex items-center color-white justify-center">
          advanced settings
          <Stack
            direction={{ xs: "column", sm: "row" }}
            spacing={{ xs: 1, sm: 2, md: 4 }}
          ></Stack>
        </div>
      </Stack>
    </Box>
  );
};

const MinimizedSearchBox = ({
  setMinimized,
}: {
  setMinimized: Dispatch<SetStateAction<boolean>>;
}) => {
  return (
    <Stack direction="row" spacing={2} className="flex items-center">
      <SearchIcon />
      <Divider orientation="vertical" />
      <CommandIcon /> K
      <Divider orientation="vertical" />
      <ChatIcon />
    </Stack>
  );
};

export const CenteredFloatingSearhBox = ({
  handleSearch,
  searchQuery,
  setSearchQuery,
  inSearchSession,
}: SearchBoxProps) => {
  const divRef = useRef<HTMLDivElement>(null);
  const [isMinimized, setIsMinimized] = useState(true);
  const [lastScrollY, setLastScrollY] = useState(0);
  const [searchVisible, setSearchVisible] = useState(true);

  const clickMinimized = () => {
    if (isMinimized) {
      setIsMinimized(false);
    }
  };

  const handleKeyDown = (event: KeyboardEvent) => {
    if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") {
      event.preventDefault();
      setIsMinimized((prevState) => !prevState);
    }
  };

  const clickOutFromSearch = (event: MouseEvent) => {
    if (divRef.current && !divRef.current.contains(event.target as Node)) {
      setIsMinimized(true);
    }
  };
  const handleScroll = () => {
    const currentScrollY = window.scrollY;
    if (currentScrollY > lastScrollY) {
      setSearchVisible(false);
    }
    {
      setSearchVisible(true);
    }
    setLastScrollY(currentScrollY);
  };

  useEffect(() => {
    window.addEventListener("keydown", handleKeyDown);
    window.addEventListener("mousedown", clickOutFromSearch);
    window.addEventListener("scroll", handleScroll);

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
      window.removeEventListener("mousedown", clickOutFromSearch);
      window.removeEventListener("scroll", handleScroll);
    };
  }, []);

  return (
    <motion.div
      ref={divRef}
      layout
      data-isOpen={!isMinimized}
      initial={{
        width: "20%",
      }}
      animate={{
        height: "auto",
        width: "auto",
        display: searchVisible ? "block" : "none",
      }}
      style={{
        position: "fixed",
        bottom: "30px",
        backgroundColor: "white",
        borderRadius: "10px",
        border: "2px solid grey",
        padding: "10px",
        zIndex: 1000,
      }}
      className="parent"
    >
      {isMinimized ? (
        <div onClick={clickMinimized}>
          <MinimizedSearchBox setMinimized={setIsMinimized} />
        </div>
      ) : (
        <SearchBox
          searchQuery={searchQuery}
          setSearchQuery={setSearchQuery}
          handleSearch={handleSearch}
          inSearchSession={inSearchSession}
        />
      )}
    </motion.div>
  );
};

export default SearchBox;
