import { Dispatch, SetStateAction, useState, useEffect, useRef } from "react";
import { Input, Button, Grid, Stack, Divider, Box } from "@mui/joy";
import Tooltip from "@mui/joy/Tooltip";
import { motion, AnimatePresence } from "framer-motion"; // Import necessary components from framer-motion
import { CommandIcon, SearchIcon, ChatIcon } from "@/components/Icons";

interface SearchBoxProps {
  handleSearch: () => Promise<void>;
  searchQuery: string;
  setSearchQuery: Dispatch<SetStateAction<string>>;
  inSearchSession: boolean;
  setChatVisible: Dispatch<SetStateAction<boolean>>;
}

interface extraProperties {
  match_name: string;
  match_source: string;
  match_doctype: string;
  match_docket_id: string;
  match_document_class: string;
  match_author: string;
}
const extraPropertiesInformation = {
  match_name: {
    displayName: "Name",
    description: "The name associated with the search item.",
    details: "Searches for items approximately matching the title",
  },
  match_source: {
    displayName: "Source",
    description: "The ",
    details: "Filters results matching the provided source exactly.",
  },
  match_doctype: {
    displayName: "Document Type",
    description: "The type or category of the document.",
    details: "Searches for items that match the specified document type.",
  },
  match_docket_id: {
    displayName: "Docket ID",
    description: "The unique identifier for the docket.",
    details: "Filters search results based on the docket ID.",
  },
  match_document_class: {
    displayName: "Document Class",
    description: "The classification or category of the document.",
    details: "Searches for documents that fall under the specified class.",
  },
  match_author: {
    displayName: "Author",
    description: "The author of the document.",
    details: "Searches for items created or written by the specified author.",
  },
};
const emptyExtraProperties: extraProperties = {
  match_name: "",
  match_source: "",
  match_doctype: "",
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
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setQueryOptions((prevOptions) => ({
      ...prevOptions,
      [name]: value,
    }));
  };
  const [showAdvancedSettings, setShowAdvancedSettings] = useState(false);

  return (
    <>
      <div className="flex items-center color-white justify-center">
        <Stack direction={{ xs: "column" }} spacing={{ xs: 1, sm: 2, md: 4 }}>
          <div className="flex items-center color-white justify-center">
            <span
              onClick={() => setShowAdvancedSettings(!showAdvancedSettings)}
              style={{ cursor: "pointer", textDecoration: "underline" }}
            >
              {showAdvancedSettings
                ? "Hide advanced settings"
                : "Show advanced settings"}
            </span>
          </div>
          <AnimatePresence initial={false}>
            {showAdvancedSettings && (
              <motion.div
                initial={{ height: 0, opacity: 0 }}
                animate={{ height: "auto", opacity: 1 }}
                exit={{ height: 0, opacity: 0 }}
                transition={{ duration: 0.2 }} // Duration of the animation in seconds
              >
                <Grid container spacing={1}>
                  {Object.keys(queryOptions).map((key, index) => {
                    const extraInfo =
                      extraPropertiesInformation[key as keyof extraProperties];
                    return (
                      <Grid item xs={12} sm={6} key={index}>
                        <Tooltip title={extraInfo.description} variant="solid">
                          <p>{extraInfo.displayName}</p>
                        </Tooltip>
                        <Input
                          type="text"
                          id={key}
                          name={key}
                          value={queryOptions[key as keyof extraProperties]}
                          onChange={handleChange}
                          title={extraInfo.displayName}
                        />
                      </Grid>
                    );
                  })}
                </Grid>
              </motion.div>
            )}
          </AnimatePresence>
        </Stack>
      </div>
    </>
  );
};

const SearchBox = ({
  handleSearch,
  searchQuery,
  setSearchQuery,
  inSearchSession,
}: SearchBoxProps) => {
  const [queryOptions, setQueryOptions] =
    useState<extraProperties>(emptyExtraProperties);

  const textRef = useRef<HTMLInputElement>(null);
  const handleEnter = (event: any) => {
    if (event.key === "Enter") {
      // Trigger function when "Enter" is pressed
      handleSearch();
    }
  };

  useEffect(() => {
    textRef.current?.focus();
  }, []);
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
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full"
            placeholder="Search"
            style={{ backgroundColor: "white" }}
            ref={textRef}
            onKeyDown={handleEnter}
          />
          <button
            style={{
              backgroundColor: "--brand-yellow-rgb",
              color: "black",
              borderRadius: "10px",
              border: "2px solid grey",
              padding: "2px",
            }}
            className="max-w-60"
            onClick={handleSearch}
          >
            <SearchIcon />
          </button>
        </Stack>
        <AdvancedSettings
          setQueryOptions={setQueryOptions}
          queryOptions={queryOptions}
        />
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
  const [isMacOS, setIsMacOS] = useState(false);

  useEffect(() => {
    if (navigator.platform.toUpperCase().indexOf("MAC") >= 0) {
      setIsMacOS(true);
    }
  }, []);

  const handleSearchClick = () => {
    setMinimized(false);
  };
  const handleChatClick = (e: React.MouseEvent<HTMLElement, MouseEvent>) => {
    e.stopPropagation(); // This will prevent the div's onClick from firing
    setChatVisible(true);
    console.log("chat element clicked");
  };

  return (
    <Stack direction="row" spacing={2} className="flex items-center">
      <Tooltip
        title={
          isMacOS ? (
            <>
              <CommandIcon /> K
            </>
          ) : (
            "Ctrl K"
          )
        }
      >
        <div onClick={handleSearchClick}>
          <SearchIcon />
        </div>
      </Tooltip>
      <Tooltip
        title={
          isMacOS ? (
            <>
              <CommandIcon /> J
            </>
          ) : (
            "Ctrl J"
          )
        }
      >
        <button onClick={handleChatClick}>
          <ChatIcon />
        </button>
      </Tooltip>
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
        <div>
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
