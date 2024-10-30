"use client";
"use client";
import React, { Dispatch, SetStateAction, useState, useRef } from "react";
import { ConversationView } from '../ConversationVeiw';
import { type ClassValue } from 'clsx';
import {
  Box,
  TextField,
  Typography,
  Container,
  Grid,
  Paper,
} from "@mui/material";
import {
  article,
  exampleBattles,
  Organization,
  Battle,
  Faction,
  Action,
} from "@/utils/interfaces";
import test from "node:test";
import { BasicDocumentFiltersList } from '@/components/DocumentFilters';
import { emptyQueryOptions, QueryFilterFields, CaseFilterFields } from "@/lib/filters";



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
};

const testFiling: Filing = {
  "id": "0",
  "url": "https://documents.dps.ny.gov/public/Common/ViewDoc.aspx?DocRefId={7F4AA7FC-CF71-4C2B-8752-A1681D8F9F46}",
  "date": "05/12/2022",
  "lang": "en",
  "title": "Press Release - PSC Announces CLCPA Tracking Initiative",
  "author": "Public Service Commission",
  "source": "Public Service Commission",
  "language": "en",
  "extension": "pdf",
  "file_class": "Press Releases",
  "item_number": "3",
  "author_organisation": "Public Service Commission"
}

const filings: Filing[] = [testFiling];

const TableFilters = ({ searchFilters, setSearchFilters }: { searchFilters: QueryFilterFields, setSearchFilters: Dispatch<SetStateAction<QueryFilterFields>> }) => {
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
}

const TableRow = ({ filing }: { filing: Filing }) => {
  return (
    <tr className="border-b border-gray-200">
      <td style={{ borderRight: 'solid 1px #f00', borderLeft: 'solid 1px #f00' }}>{filing.date}</td>
      <td style={{ borderRight: 'solid 1px #f00', borderLeft: 'solid 1px #f00' }}>{filing.title}</td>
      <td style={{ borderRight: 'solid 1px #f00', borderLeft: 'solid 1px #f00' }}>{filing.author}</td>
      <td style={{ borderRight: 'solid 1px #f00', borderLeft: 'solid 1px #f00' }}>{filing.source}</td>
      <td style={{ borderRight: 'solid 1px #f00', borderLeft: 'solid 1px #f00' }}>{filing.item_number}</td>
      <td style={{ borderRight: 'solid 1px #f00', borderLeft: 'solid 1px #f00' }}><a href={filing.url}>View</a></td>
    </tr>
  )
}
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
            <tr className="border-b border-gray-200">
              <td>{filing.date}</td>
              <td>{filing.title}</td>
              <td>{filing.author}</td>
              <td>{filing.source}</td>
              <td>{filing.item_number}</td>
              <td><a href={filing.url} target="_blank" rel="noopener noreferrer">View</a></td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
const ConversationComponent: React.FC = ({
}) => {
  const [loading, setLoading] = useState(false);
  const [searchFilters, setSearchFilters] = useState<QueryFilterFields>(emptyQueryOptions);
  const [isFocused, setIsFocused] = useState(false);
  const divRef = useRef(null);

  const showFilters = () => {
    setIsFocused(!isFocused);
  };

  return (
    <div className="w-full h-full p-10 card grid grid-flow-col auto-cols-2 box-border border-2 border-black ">
      <div
        style={{
          display: isFocused ? 'block' : 'none',
          padding: '10px',
          transition: 'width 0.3s ease-in-out'
        }}
      >
        <button onClick={showFilters} className="btn "
          style={{
            display: isFocused ? 'flex' : 'none',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="32"
            height="32"
            viewBox="0 0 512 512">
            <polygon
              points="400 145.49 366.51 112 256 222.51 145.49 112 112 145.49 222.51 256 112 366.51 145.49 400 256 289.49 366.51 400 400 366.51 289.49 256 400 145.49" />
          </svg>
        </button>
        <BasicDocumentFiltersList
          queryOptions={searchFilters}
          setQueryOptions={setSearchFilters}
          showQueries={CaseFilterFields}
        />
      </div>
      <div className=" p-10">
        <div id="conversation-header p-10 justify-between">
        </div>
        <h1 className=" text-2xl font-bold">Conversation</h1>
        <button onClick={showFilters} className="btn btn-outline"
          style={{
            display: !isFocused ? 'inline-block' : 'none',
          }}
        >
          Filters
        </button>
        <div className="w-full overflow-x-scroll">
          {loading ? <div>Loading...</div> : <FilingTable filings={filings} />}
        </div>
      </div>
    </div>
  );
};


export default ConversationComponent;