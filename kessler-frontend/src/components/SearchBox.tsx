import { Dispatch, SetStateAction, useState } from "react";

interface SearchBoxProps {
  handleSearch: () => Promise<void>;
  searchQuery: string;
  setSearchQuery: Dispatch<SetStateAction<string>>;
}


const SearchBox = ({ handleSearch, searchQuery, setSearchQuery}: SearchBoxProps) => {

  return (
    <div className="z-10 w-full items-center justify-between font-mono text-sm lg:flex searchbox m-20%">
      {/* Your UI elements here */}
      <div className="flex flex-col justify-self-center mx-auto w-full">
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="border border-gray-300 rounded-lg px-4 py-2 bg-gray-200 text-gray-800 mb-2 w-30%"
        />
        <button
          onClick={handleSearch}
          className="border border-gray-300 rounded-lg px-4 py-2"
        >
          Search
        </button>
      </div>
    </div>
  );
};

export default SearchBox;