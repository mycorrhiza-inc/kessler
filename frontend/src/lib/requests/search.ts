import { QueryDataFile, QueryFilterFields } from "@/lib/filters";
import { Filing } from "@/lib/types/FilingTypes";
import axios from "axios";
import { apiURL } from "../env_variables";
import {
  CompleteFileSchema,
  CompleteFileSchemaValidator,
} from "../types/backend_schemas";

export const getSearchResults = async (
  queryData: QueryDataFile,
): Promise<Filing[]> => {
  const searchQuery = queryData.query;
  console.log("query data", queryData);
  const searchFilters = queryData.filters;
  console.log("searchhing for", searchFilters);
  try {
    const searchResults: Filing[] = await axios
      // .post("https://api.kessler.xyz/v2/search", {
      .post(`${apiURL}/v2/search`, {
        query: searchQuery,
        filters: {
          name: searchFilters.match_name,
          author: searchFilters.match_author,
          docket_id: searchFilters.match_docket_id,
          file_class: searchFilters.match_file_class,
          doctype: searchFilters.match_doctype,
          source: searchFilters.match_source,
          date_from: searchFilters.match_after_date,
          date_to: searchFilters.match_before_date,
        },
        start_offset: queryData.start_offset,
        max_hits: 20,
      })
      // check error conditions
      .then((response) => {
        if (response.data?.length === 0 || typeof response.data === "string") {
          return [];
        }
        const filings_promise: Promise<Filing[]> = ParseFilingData(
          response.data,
        );
        return filings_promise;
      });
    console.log("getting data");
    // console.log(searchResults);
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

export const completeFileSchemaGet = async (url: string) => {
  const response = await axios.get(url);
  if (response.status !== 200) {
    throw new Error(
      "Error fetching data with response code: " + response.status,
    );
  }
  if (response.data === undefined) {
    throw new Error("No data returned from server");
  }

  try {
    // Parse and validate the response data
    const validatedData: CompleteFileSchema = CompleteFileSchemaValidator.parse(
      response.data,
    );
    return validatedData;
  } catch (error) {
    if (error instanceof Error) {
      throw new Error(`Invalid response data structure: ${error.message}`);
    }
    throw error;
  }
};

export const ParseFilingData = async (filingData: any) => {
  const filings_promises: Promise<Filing>[] = filingData.map(async (f: any) => {
    const mdata_str = f.Mdata;
    if (!mdata_str) {
      console.log("no metadata string, fetching from source");
      const docID = f.sourceID;

      const response = await axios.get(
        // `http://api.kessler.xyz/v2/public/files/${id}/metadata`
        `${apiURL}/v2/public/files/${docID}/metadata`,
      );
      const metadata = response.data.mdata;
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
    const metadata = f.mdata;
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
