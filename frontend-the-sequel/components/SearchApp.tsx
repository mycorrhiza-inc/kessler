"use client";
import axios from "axios";
import SearchResult from "@/components/SearchResult";
import { useState } from "react";

const ResultContainer = ({ children }: { children: React.ReactNode }) => {
  return <div className="min-h-50vh">{children}</div>;
};

export default function SearchApp() {
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<any[]>([]);
  const [isSearching, setIsSearching] = useState(false);

  // ...

  const handleSearch = async () => {
    // try {
    // } catch (error) {
    //   // Handle the error here
    // }
    console.log("Searching for:", searchQuery);
    setSearchResults([]);
    setIsSearching(true);

    setTimeout(async () => {
      try {
        const response = await axios.post("http://localhost:4041/search", {
          query: searchQuery,
        });
        setSearchResults(response.data);

        // // Move the searchbox class to the top of the page
        // const searchbox = document.querySelector(".searchbox");
        // if (searchbox) {
        //   searchbox.classList.add("animate-searchbox");
        //   setTimeout(() => {
        //     searchbox.classList.remove("searchbox");
        //   }, 2000);
        // }

        // setSearchResults([
        //   {
        //     name: "Name 1",
        //     text: "Text 1",
        //     docketId: "Docket ID 1",
        //   },
        //   {
        //     name: "Name 2",
        //     text: "Text 2",
        //     docketId: "Docket ID 2",
        //   },
        //   {
        //     name: "Name 3",
        //     text: "Text 3",
        //     docketId: "Docket ID 3",
        //   },
        // ]);
      } catch (error) {
        // Handle the error here
        console.log(error);
      } finally {
        setIsSearching(false);
      }
    }, 300);
  };

  // ...

  return (
    <main className="flex flex-col min-w-screen min-h-screen items-center justify-center">
      <div
        className="viewport"
        style={{ width: "80vw", height: "80vh", margin: 0, padding: "50px" }}
      >
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
        <div className="searchResults">
          <div
            className="searchResultsContent"
            style={{ justifyContent: "center" }}
          >
            {isSearching ? (
              <ResultContainer>
                <div className="flex items-center justify-center mt-4">
                  <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-white"></div>
                </div>
              </ResultContainer>
            ) : (
              <ResultContainer>
                {searchResults.map((result, index) => (
                  <SearchResult key={index} data={result} />
                ))}
              </ResultContainer>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}
