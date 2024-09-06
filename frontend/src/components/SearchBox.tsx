import { Dispatch, SetStateAction, useState, useEffect, useRef } from "react";
import { Input, Button, Grid, Stack, Divider, Box } from "@mui/joy";
import { motion } from "framer-motion";
import { CommandIcon, SearchIcon, ChatIcon } from "@/components/Icons";

interface SearchBoxProps {
  handleSearch: () => Promise<void>;
  searchQuery: string;
  setSearchQuery: Dispatch<SetStateAction<string>>;
  inSearchSession: boolean;
}

interface extraProperties {
  match_name: string;
  match_source: string;
  match_doctype: string;
  match_stage: string;
  match_docket_id: string;
  match_document_class: string;
  match_author: string;
}
const emptyExtraProperties: extraProperties = {
  match_name: "",
  match_source: "",
  match_doctype: "",
  match_stage: "",
  match_docket_id: "",
  match_document_class: "",
  match_author: "",
};

const AdvancedSettings = ({
  queryOptions,
  setQueryOptions,
}: {
  queryOptions: extraProperties;
  setQueryOptions: Dispatch<SetStateAction<extraProperties>>;
}) => {
  return <p>Bob Lob Law</p>;
};

const SearchBox = ({
  handleSearch,
  searchQuery,
  setSearchQuery,
  inSearchSession,
}: SearchBoxProps) => {
  const [showAdvancedSettings, setShowAdvancedSettings] = useState(false);
  const [queryOptions, setQueryOptions] =
    useState<extraProperties>(emptyExtraProperties);

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
        <div className="flex items-center color-white justify-center">
          <Button
            onClick={() => setShowAdvancedSettings(!showAdvancedSettings)}
          >
            {showAdvancedSettings
              ? "Hide advanced settings"
              : "Show advanced settings"}
          </Button>
        </div>
        {showAdvancedSettings && (
          <Stack
            direction={{ xs: "column", sm: "row" }}
            spacing={{ xs: 1, sm: 2, md: 4 }}
          >
            <AdvancedSettings
              queryOptions={queryOptions}
              setQueryOptions={setQueryOptions}
            />
          </Stack>
        )}
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
        color: "black",
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
