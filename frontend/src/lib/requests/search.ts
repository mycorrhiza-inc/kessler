import { QueryDataFile, QueryFilterFields } from "@/lib/filters";
import { Filing } from "@/lib/types/FilingTypes";
import axios from "axios";
import { apiURL } from "../env_variables";

export const getSearchResults = async (
  queryData: QueryDataFile,
): Promise<Filing[]> => {
  const searchQuery = queryData.query;
  const searchFilters = queryData.filters;
  console.log(`searchhing for ${searchQuery}`);
  try {
    const searchResults = await axios
      // .post("https://api.kessler.xyz/v2/search", {
      .post(`${apiURL}/v2/search`, {
        query: searchQuery,
        filters: {
          name: searchFilters.match_name,
          author: searchFilters.match_author,
          docket_id: searchFilters.match_docket_id,
          doctype: searchFilters.match_doctype,
          source: searchFilters.match_source,
        },
      })
      // check error conditions
      .then((response) => {
        if (response.data.length === 0 || typeof response.data === "string") {
          return [];
        }
        const filings: Filing[] = ParseFilingData(response.data);
        return filings;
      });
    console.log("getting data");
    console.log(searchResults);
    return searchResults;
  } catch (error) {
    console.log(error);
    throw error;
  }
};

export const getRecentFilings = async (page?: number) => {
  if (!page) {
    page = 0;
  }
  const response = await axios.post(
    `${apiURL}/v2/recent_updates`,
    // "http://api.kessler.xyz/v2/recent_updates",
    {
      page: page,
    },
  );
  console.log("recent data", response.data);
  if (response.data.length > 0) {
    return response.data;
  }
};

export const getFilingMetadata = async (id: string) => {
  const response = await axios.get(
    // `http://api.kessler.xyz/v2/public/files/${id}/metadata`
    `${apiURL}/v2/public/files/${id}/metadata`,
  );
  return ParseFilingData([response.data])[0];
};

export const ParseFilingData = (filingData: any) => {
  const filings = filingData.map((f: any) => {
    console.log("fdata", f);
    const mdata_str = f.Mdata;
    if (!mdata_str) {
      console.log("no metadata string");
      const newFiling: Filing = {
        id: f.docID,
        title: f.name,
        source: f.sourceID,
        // lang: f.metadata.lang,
        // date: f.metadata.date,
        // author: f.metadata.author,
        // source: f.metadata.source,
        // language: f.metadata.language,
        // item_number: f.metadata.item_number,
        // author_organisation: f.metadata.author_organizatino,
        // url: f.metadata.url,
      };
      return newFiling;
    }
    console.log("metadata string: ", mdata_str);
    const metadata = JSON.parse(atob(f.Mdata));
    console.log("metadata: ", metadata);
    const newFiling: Filing = {
      id: f.docID,
      title: f.name,
      source: f.sourceID,
      lang: metadata.lang,
      date: metadata.date,
      author: metadata.author,
      language: metadata.language,
      item_number: metadata.item_number,
      author_organisation: metadata.author_organizatino,
      url: metadata.url,
    };
    return newFiling;
  });
  return filings;
};

export default getSearchResults;
