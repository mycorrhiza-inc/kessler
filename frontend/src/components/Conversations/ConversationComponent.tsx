"use client";
"use client";
import React, {
  Dispatch,
  SetStateAction,
  useState,
  useRef,
  useMemo,
  Suspense,
} from "react";
import { BasicDocumentFiltersList } from "@/components/DocumentFilters";
import {
  emptyQueryOptions,
  QueryFilterFields,
  CaseFilterFields,
  InheritedFilterValues,
} from "@/lib/filters";
import Modal from "../styled-components/Modal";
import DocumentModalBody from "../DocumentModalBody";
import { AnimatePresence, motion } from "framer-motion";
import axios from "axios";

type Filing = {
  id: string;
  lang: string;
  title: string;
  date: string;
  author: string;
  source: string;
  language: string;
  extension: string;
  file_class: string;
  item_number: string;
  author_organisation: string;
  url: string;
  uuid: string;
};

const testFiling: Filing = {
  id: "0",
  url: "https://documents.dps.ny.gov/public/Common/ViewDoc.aspx?DocRefId={7F4AA7FC-CF71-4C2B-8752-A1681D8F9F46}",
  date: "05/12/2022",
  lang: "en",
  title: "Press Release - PSC Announces CLCPA Tracking Initiative",
  author: "Public Service Commission",
  source: "Public Service Commission",
  language: "en",
  extension: "pdf",
  file_class: "Press Releases",
  item_number: "3",
  author_organisation: "Public Service Commission",
  uuid: "3c4ba5f3-febc-41f2-aa86-2820db2b459a",
};

const filings: Filing[] = [testFiling];

const TableFilters = ({
  searchFilters,
  setSearchFilters,
}: {
  searchFilters: QueryFilterFields;
  setSearchFilters: Dispatch<SetStateAction<QueryFilterFields>>;
}) => {
  return (
    <div className="collapse w-auto ">
      <input type="checkbox" />
      <div className="collapse-title font-medium">Filters</div>
      <div className="collapse-content flex flex-col space-y-1 sm:space-y-2 md:space-y-3 items-center  justify-center">
        <BasicDocumentFiltersList
          queryOptions={searchFilters}
          setQueryOptions={setSearchFilters}
          showQueries={CaseFilterFields}
        />
      </div>
    </div>
  );
};

const TableRow = ({ filing }: { filing: Filing }) => {
  const [open, setOpen] = useState(false);
  return (
    <>
      <tr
        className="border-b border-gray-200"
        onClick={() => {
          setOpen((previous) => !previous);
        }}
      >
        <td>{filing.date}</td>
        <td>{filing.title}</td>
        <td>{filing.author}</td>
        <td>{filing.source}</td>
        <td>{filing.item_number}</td>
        <td>
          <a href={filing.url}>View</a>
        </td>
      </tr>
      <Modal open={open} setOpen={setOpen}>
        <DocumentModalBody
          open={open}
          objectId={filing.uuid}
          overridePDFUrl={filing.url}
        />
      </Modal>
    </>
  );
};
const FilingTable = ({ filings }: { filings: Filing[] }) => {
  return (
    <div className="overflow-x-scroll max-h-[500px] overflow-y-auto">
      <table className="w-full divide-y divide-gray-200 table">
        <tbody>
          <tr className="border-b border-gray-200">
            <th className="text-left p-2 sticky top-0 bg-white">Date Filed</th>
            <th className="text-left p-2 sticky top-0 bg-white">Title</th>
            <th className="text-left p-2 sticky top-0 bg-white">Author</th>
            <th className="text-left p-2 sticky top-0 bg-white">Source</th>
            <th className="text-left p-2 sticky top-0 bg-white">Item Number</th>
          </tr>
          {filings.map((filing) => (
            <TableRow filing={filing} />
          ))}
        </tbody>
      </table>
    </div>
  );
};
const ConversationComponent = ({
  inheritedFilters,
}: {
  inheritedFilters: InheritedFilterValues;
}) => {
  const disabledFilters = useMemo(() => {
    return inheritedFilters.map((val) => {
      return val.filter;
    });
  }, [inheritedFilters]);

  const initialFilterState = useMemo(() => {
    var initialFilters = emptyQueryOptions;
    inheritedFilters.map((val) => {
      initialFilters[val.filter] = val.value;
    });
    return initialFilters;
  }, [inheritedFilters]);
  const [searchFilters, setSearchFilters] =
    useState<QueryFilterFields>(initialFilterState);
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");

  const searchResultsGet = async () => {
    console.log(`searchhing for ${searchQuery}`);
    try {
      const response = await axios.post("https://api.kessler.xyz/v2/search", {
        query: searchQuery,
        filters: {
          name: searchFilters.match_name,
          author: searchFilters.match_author,
          docket_id: searchFilters.match_docket_id,
          doctype: searchFilters.match_doctype,
          source: searchFilters.match_source,
        },
      });
      if (response.data.length === 0) {
        return;
      }
      if (typeof response.data === "string") {
        setSearchResults([]);
        return;
      }
      console.log("getting data");
      console.log(response.data);
      return response.data;
    } catch (error) {
      console.log(error);
    }
  };
  // This should work with the async promises for
  const searchResultsHandler = async () => {
    setSearchResults(await searchResultsGet());
  };

  const [isFocused, setIsFocused] = useState(false);
  const showFilters = () => {
    setIsFocused(!isFocused);
  };

  return (
    <div className="w-full h-full p-10 card grid grid-flow-col auto-cols-2 box-border border-2 border-black ">
      <AnimatePresence>
        {isFocused && (
          <motion.div
            style={{
              padding: "10px",
              transition: "width 0.3s ease-in-out",
            }}
            initial={{ x: "-50%" }}
            animate={{ x: "0" }}
            exit={{ x: "-50%", opacity: 0 }}
          >
            <button
              onClick={showFilters}
              className="btn "
              style={{
                display: isFocused ? "flex" : "none",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="32"
                height="32"
                viewBox="0 0 512 512"
              >
                <polygon points="400 145.49 366.51 112 256 222.51 145.49 112 112 145.49 222.51 256 112 366.51 145.49 400 256 289.49 366.51 400 400 366.51 289.49 256 400 145.49" />
              </svg>
            </button>
            <BasicDocumentFiltersList
              queryOptions={searchFilters}
              setQueryOptions={setSearchFilters}
              showQueries={CaseFilterFields}
              disabledQueries={disabledFilters}
            />
          </motion.div>
        )}
      </AnimatePresence>
      <div className=" p-10">
        <div id="conversation-header p-10 justify-between"></div>
        <h1 className=" text-2xl font-bold">Conversation</h1>
        <button
          onClick={showFilters}
          className="btn btn-outline"
          style={{
            display: !isFocused ? "inline-block" : "none",
          }}
        >
          Filters
        </button>
        <div className="w-full overflow-x-scroll">
          <Suspense fallback={<div>Loading...</div>}>
            <FilingTable filings={filings} />
          </Suspense>
        </div>
      </div>
    </div>
  );
};

export default ConversationComponent;
