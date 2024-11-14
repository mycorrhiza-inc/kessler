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
    const searchResults: Filing[] = await axios
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
        const filings_promise: Promise<Filing[]> = ParseFilingData(
          response.data,
        );
        return filings_promise;
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

export const getFilingMetadata = async (id: string): Promise<Filing> => {
  const response = await axios.get(
    // `http://api.kessler.xyz/v2/public/files/${id}/metadata`
    `${apiURL}/v2/public/files/${id}/metadata`,
  );
  const filings = await ParseFilingData([response.data]);
  return filings[0];
};

export const ParseFilingData = async (filingData: any) => {
  const filings_promises: Promise<Filing>[] = filingData.map(async (f: any) => {
    console.log("fdata", f);
    const mdata_str = f.Mdata;
    if (!mdata_str) {
      console.log("no metadata string, fetching from source");
      const docID = f.sourceID;

      const response = await axios.get(
        // `http://api.kessler.xyz/v2/public/files/${id}/metadata`
        `${apiURL}/v2/public/files/${docID}/metadata`,
      );
      const metadata = JSON.parse(atob(response.data.Mdata));
      console.log("metadata: ", metadata);
      const newFiling: Filing = {
        // These names are swaped in the backend, maybe change later
        id: f.sourceID,
        title: f.name,
        source: f.docID,
        lang: metadata.lang,
        date: metadata.date,
        author: metadata.author,
        language: metadata.language,
        item_number: metadata.item_number,
        author_organisation: metadata.author_organizatino,
        file_class: metadata.file_class,
        url: metadata.url,
      };
      return newFiling;
    }
    console.log("metadata string: ", mdata_str);
    const metadata = JSON.parse(atob(f.Mdata));
    console.log("metadata: ", metadata);
    const newFiling: Filing = {
      id: metadata.id,
      title: metadata.title,
      source: metadata.docket_id,
      lang: metadata.lang,
      date: metadata.date,
      author: metadata.author,
      language: metadata.language,
      item_number: metadata.item_number,
      author_organisation: metadata.author_organizatino,
      file_class: metadata.file_class,
      url: metadata.url,
    };
    return newFiling;
  });
  const filings: Filing[] = await Promise.all(filings_promises);
  return filings;
};

export default getSearchResults;
