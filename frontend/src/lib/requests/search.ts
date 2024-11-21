import { QueryDataFile, QueryFilterFields } from "@/lib/filters";
import { Filing } from "@/lib/types/FilingTypes";
import axios from "axios";
import { apiURL } from "../env_variables";
import {
  CompleteFileSchema,
  CompleteFileSchemaValidator,
} from "../types/backend_schemas";
import { z } from "zod";
import { fi } from "date-fns/locale";

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
  const valid_id = z.string().uuid().parse(id);
  const response = await axios.get(
    // `http://api.kessler.xyz/v2/public/files/${id}/metadata`
    `${apiURL}/v2/public/files/${valid_id}/metadata`,
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
      console.error("Error parsing data:", error.message);
      console.error("Data:", response.data);
      throw new Error(
        `Invalid response data structure:${response.data} \n raising error: ${error.message}`,
      );
    }
    throw error;
  }
};

export const generateFilingFromFileSchema = (
  file_schema: CompleteFileSchema,
): Filing => {
  return {
    id: file_schema.id,
    title: file_schema.name,
    source: file_schema.mdata.docID,
    lang: file_schema.lang,
    date: file_schema.mdata.date,
    author: String(file_schema.authors),
    language: file_schema.lang, // This is redundant with lang
    item_number: file_schema.mdata.item_number,
    file_class: file_schema.mdata.file_class,
    url: file_schema.mdata.url,
  };
};

export const ParseFilingData = async (filingData: any) => {
  const generate_filing = async (f: any) => {
    var docID = "";
    try {
      docID = z.string().uuid().parse(f.sourceID);
    } catch (error) {
      console.log(error);
      console.log(f);
      return;
    }
    console.log("no metadata string, fetching from source");
    const metadata_url = `${apiURL}/v2/public/files/${docID}/metadata`;
    const completeFileSchema = await completeFileSchemaGet(metadata_url);
    const newFiling: Filing = generateFilingFromFileSchema(completeFileSchema);
    return newFiling;
  };

  const filings_promises: Promise<Filing | undefined>[] =
    filingData.map(generate_filing);
  const filings_with_errors = await Promise.all(filings_promises);
  const filings = filings_with_errors.filter((f): boolean => Boolean(f));
  return filings;
};

export default getSearchResults;
