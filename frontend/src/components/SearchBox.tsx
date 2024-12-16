import { Dispatch, SetStateAction, useState, useEffect, useRef } from "react";
import Tooltip from "@mui/joy/Tooltip";
import { motion, AnimatePresence } from "framer-motion"; // Import necessary components from framer-motion
import { SearchIcon, ChatIcon, FileUploadIcon } from "@/components/Icons";

import { QueryFilterFields, allFilterFields } from "@/lib/filters";
import { BasicDocumentFiltersGrid } from "@/components/Filters/DocumentFilters";

const AdvancedFilters = ({
  queryOptions,
  setQueryOptions,
}: {
  queryOptions: QueryFilterFields;
  setQueryOptions: Dispatch<SetStateAction<QueryFilterFields>>;
}) => {
  const [showAdvancedFilters, setShowAdvancedFilters] = useState(true);

  return (
    <>
      <div className="flex items-center justify-center">
        <div className="flex flex-col space-y-1 sm:space-y-2 md:space-y-4">
          <div className="flex items-center  justify-center">
            <span
              onClick={() => setShowAdvancedFilters(!showAdvancedFilters)}
              style={{ cursor: "pointer", textDecoration: "underline" }}
            >
              {showAdvancedFilters
                ? "Hide Advanced Filters"
                : "Show Advanced Filters"}
            </span>
          </div>
          <AnimatePresence initial={false}>
            {showAdvancedFilters && (
              <motion.div
                initial={{ height: 0, width: 0, opacity: 0 }}
                animate={{ height: "auto", width: "auto", opacity: 1 }}
                exit={{ height: 0, width: 0, opacity: 0 }}
                transition={{ duration: 0.3 }} // Duration of the animation in seconds
              >
                <BasicDocumentFiltersGrid
                  queryOptions={queryOptions}
                  setQueryOptions={setQueryOptions}
                  showQueries={allFilterFields}
                />
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </>
  );
};

interface SearchBoxProps {
  handleSearch: () => Promise<void>;
  searchQuery: string;
  setSearchQuery: Dispatch<SetStateAction<string>>;
  inSearchSession: boolean;
  chatVisible: boolean;
  setChatVisible: Dispatch<SetStateAction<boolean>>;
  queryOptions: QueryFilterFields;
  setQueryOptions: Dispatch<SetStateAction<QueryFilterFields>>;
}
const SearchBox = ({
  handleSearch,
  searchQuery,
  setSearchQuery,
  inSearchSession,
  queryOptions,
  setQueryOptions,
}: SearchBoxProps) => {
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
      <AdvancedFilters
        setQueryOptions={setQueryOptions}
        queryOptions={queryOptions}
      />
    </div>
  );
};
interface UploadBoxProps {}
const UploadBox = ({}: UploadBoxProps) => {
  // const textRef = useRef<HTMLInputElement>(null);
  // const handleEnter = (event: any) => {
  //   if (event.key === "Enter") {
  //     // Trigger function when "Enter" is pressed
  //     handleSearch();
  //   }
  // };
  //
  // useEffect(() => {
  //   textRef.current?.focus();
  // }, []);
  return (
    <>
      <label className="form-control w-full max-w-xs">
        <div className="label">
          <span className="label-text">
            Pick a Folder to Upload all Files in Folder
          </span>
        </div>
        <input
          type="file"
          className="file-input file-input-bordered file-input-accent w-full max-w-xs"
        />
        <div className="label"></div>
      </label>
      <label className="form-control w-full max-w-xs">
        <div className="label">
          <span className="label-text">Pick an Individual File to Upload</span>
        </div>
        <input
          type="file"
          className="file-input file-input-bordered file-input-primary w-full max-w-xs"
        />
        <div className="label"></div>
      </label>
    </>
  );
};

const MinimizedSearchBox = ({
  setSearchMinimized,
  setUploadMinimized,
  setChatVisible,
}: {
  setSearchMinimized: Dispatch<SetStateAction<boolean>>;
  setUploadMinimized: Dispatch<SetStateAction<boolean>>;
  setChatVisible: Dispatch<SetStateAction<boolean>>;
}) => {
  const [isMacOS, setIsMacOS] = useState(false);

  useEffect(() => {
    if (navigator.platform.toUpperCase().indexOf("MAC") >= 0) {
      setIsMacOS(true);
    }
  }, []);

  const handleSearchClick = () => {
    setSearchMinimized(false);
    setUploadMinimized(true);
  };
  const handleUploadClick = () => {
    setUploadMinimized(false);
    setSearchMinimized(true);
    console.log("Uploaded File");
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
              <kbd className="kbd text-base-content">⌘</kbd>+
              <kbd className="kbd text-base-content">K</kbd>
            </>
          ) : (
            <>
              <kbd className="kbd text-base-content">ctrl</kbd>+
              <kbd className="kbd text-base-content">K</kbd>
            </>
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
              <kbd className="kbd text-base-content">⌘</kbd>+
              <kbd className="kbd text-base-content">J</kbd>
            </>
          ) : (
            <>
              <kbd className="kbd text-base-content">ctrl</kbd>+
              <kbd className="kbd text-base-content">J</kbd>
            </>
          )
        }
      >
        <button className="scale-150" onClick={handleChatClick}>
          <ChatIcon />
        </button>
      </Tooltip>
      <div className="w-3 h-6  mx-4"></div>
      <Tooltip title="Upload files here!">
        <button className="scale-150" onClick={handleUploadClick}>
          <FileUploadIcon />
        </button>
      </Tooltip>
    </div>
  );
};

const ActionBoxController = ({
  isSearchMinimized,
  chatVisible,
}: {
  isSearchMinimized: boolean;
  chatVisible: boolean;
}) => {};

export const CenteredFloatingSearhBox = ({
  handleSearch,
  searchQuery,
  setSearchQuery,
  chatVisible,
  setChatVisible,
  inSearchSession,
  queryOptions,
  setQueryOptions,
}: SearchBoxProps) => {
  const divRef = useRef<HTMLDivElement>(null);
  const [isSearchMinimized, setSearchMinimized] = useState(true);
  const [isUploadMinimized, setUploadMinimized] = useState(true);
  const [lastScrollY, setLastScrollY] = useState(0);
  const [searchVisible, setSearchVisible] = useState(true);
  const isEverythingMinimized = isSearchMinimized && isUploadMinimized;

  const handleKeyDown = (event: KeyboardEvent) => {
    if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") {
      event.preventDefault();
      setSearchMinimized((prevState) => !prevState);
    }
    if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "j") {
      event.preventDefault();
      setChatVisible((prevState) => !prevState);
    }
  };

  const clickOutFromSearch = (event: MouseEvent) => {
    if (divRef.current && !divRef.current.contains(event.target as Node)) {
      setSearchMinimized(true);
      setUploadMinimized(true);
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
      className="flex justify-center"
      initial={{ width: "100%" }}
      animate={chatVisible ? { width: "calc(100% - 35%)" } : { width: "100%" }}
      transition={{ type: "tween", stiffness: 200 }}
      style={{
        position: "fixed",
        bottom: "30px",
        left: "0px",
        zIndex: 1500,
      }}
    >
      <motion.div
        layout
        ref={divRef}
        data-isopen={!isSearchMinimized}
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
        {isEverythingMinimized && (
          <div>
            <MinimizedSearchBox
              setUploadMinimized={setUploadMinimized}
              setSearchMinimized={setSearchMinimized}
              setChatVisible={setChatVisible}
            />
          </div>
        )}
        {!isSearchMinimized && (
          <SearchBox
            searchQuery={searchQuery}
            setSearchQuery={setSearchQuery}
            handleSearch={handleSearch}
            inSearchSession={inSearchSession}
            chatVisible={chatVisible}
            setChatVisible={setChatVisible}
            setQueryOptions={setQueryOptions}
            queryOptions={queryOptions}
          />
        )}
        {!isUploadMinimized && <UploadBox />}
      </motion.div>
    </motion.div>
  );
};

export default SearchBox;
