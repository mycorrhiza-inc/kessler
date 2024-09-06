import { Dispatch, SetStateAction, useState, useEffect, useRef } from "react";
import { Input, Button, Grid, Stack, Divider, Box } from "@mui/joy";
import { motion } from "framer-motion";
import { CommandIcon, SearchIcon, ChatIcon } from "@/components/Icons";

interface SearchBoxProps {
  handleSearch: () => Promise<void>;
  searchQuery: string;
  setSearchQuery: Dispatch<SetStateAction<string>>;
  inSearchSession: boolean;
  setChatVisible: Dispatch<SetStateAction<boolean>>;
}

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
  setChatVisible,
}: {
  setMinimized: Dispatch<SetStateAction<boolean>>;
  setChatVisible: Dispatch<SetStateAction<boolean>>;
}) => {
  const handleSearchClick = () => {
    setMinimized(false);
  };
  const handleChatClick = (
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>
  ) => {
    e.stopPropagation(); // This will prevent the div's onClick from firing
    setChatVisible(true);
    console.log("chat element clicked");
  };
  return (
    <Stack direction="row" spacing={2} className="flex items-center">
      <div onClick={handleSearchClick}>
        <SearchIcon />
      </div>
      <Divider orientation="vertical" />
      <CommandIcon /> K
      <Divider orientation="vertical" />
      <CommandIcon /> J
      <Divider orientation="vertical" />
      <button onClick={handleChatClick}>
        <ChatIcon />
      </button>
    </Stack>
  );
};

export const CenteredFloatingSearhBox = ({
  handleSearch,
  searchQuery,
  setSearchQuery,
  setChatVisible,
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
    if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "j") {
      event.preventDefault();
      setChatVisible((prevState) => !prevState);
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
      layout
      ref={divRef}
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
        color: "black",
      }}
      className="parent"
    >
      {isMinimized ? (
        <div onClick={clickMinimized}>
          <MinimizedSearchBox
            setMinimized={setIsMinimized}
            setChatVisible={setChatVisible}
          />
        </div>
      ) : (
        <SearchBox
          searchQuery={searchQuery}
          setSearchQuery={setSearchQuery}
          handleSearch={handleSearch}
          inSearchSession={inSearchSession}
          setChatVisible={setChatVisible}
        />
      )}
    </motion.div>
  );
};

export default SearchBox;
