import { Dispatch, SetStateAction, useState, useEffect, useRef } from "react";
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
      <div className="flex items-center justify-center">
        <div className="flex flex-col space-y-1 sm:space-y-2 md:space-y-4">
          <div className="flex items-center  justify-center">
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
                initial={{ height: 0, width: 0, opacity: 0 }}
                animate={{ height: "auto", width: "auto", opacity: 1 }}
                exit={{ height: 0, width: 0, opacity: 0 }}
                transition={{ duration: 0.3 }} // Duration of the animation in seconds
              >
                <div className="grid grid-cols-2 gap-4">
                  {Object.keys(queryOptions)
                    .slice(0, 6)
                    .map((key, index) => {
                      const extraInfo =
                        extraPropertiesInformation[
                          key as keyof extraProperties
                        ];
                      return (
                        <div className="box-border" key={index}>
                          <div
                            className="tooltip"
                            data-tip={extraInfo.description}
                          >
                            <p>{extraInfo.displayName}</p>
                          </div>
                          <input
                            className="input input-bordered w-full max-w-xs"
                            type="text"
                            id={key}
                            name={key}
                            value={queryOptions[key as keyof extraProperties]}
                            onChange={handleChange}
                            title={extraInfo.displayName}
                          />
                        </div>
                      );
                    })}
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
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
    <div>
      <div className="flex flex-row space-x-2 items-center  justify-center">
        <input
          className="input input-bordered w-full"
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Search"
          ref={textRef}
          onKeyDown={handleEnter}
        />
        <button
          className="max-w-60 bg-brand-yellow-rgb text-base-content rounded-lg border-2 border-gray-500 p-1"
          onClick={handleSearch}
        >
          <SearchIcon />
        </button>
      </div>
      <AdvancedSettings
        setQueryOptions={setQueryOptions}
        queryOptions={queryOptions}
      />
    </div>
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
    setChatVisible((prevState) => !prevState);
    console.log("chat element clicked");
  };

  return (
    <div className="flex flex-row space-x-2 items-center">
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
        <div className="scale-150" onClick={handleSearchClick}>
          <SearchIcon />
        </div>
      </Tooltip>
      <div className="w-3 h-6  mx-4"></div>
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
        <button className="scale-150" onClick={handleChatClick}>
          <ChatIcon />
        </button>
      </Tooltip>
    </div>
  );
};

const ActionBoxController = ({
  isMinimized,
  chatVisible,
}: {
  isMinimized: boolean;
  chatVisible: boolean;
}) => {};

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
    <div
      className="w-full flex justify-center"
      style={{
        position: "fixed",
        bottom: "30px",
        zIndex: 1500,
      }}
    >
      <motion.div
        layout
        ref={divRef}
        data-isOpen={!isMinimized}
        initial={{}}
        animate={{
          height: "auto",
          width: "auto",
          display: searchVisible ? "block" : "none",
        }}
        style={{
          borderRadius: "10px",
          border: "2px solid grey",
          padding: "10px",
        }}
        className="parent fixed bottom-7 rounded-lg border-2 border-gray-500 bg-base-200 text-base-content"
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
    </div>
  );
};

export default SearchBox;
