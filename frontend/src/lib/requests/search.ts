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

export const getFilingMetadata = async (id: string): Promise<Filing | null> => {
  const valid_id = z.string().uuid().parse(id);
  const response = await axios.get(
    `${apiURL}/v2/public/files/${valid_id}/metadata`,
  );
  const filing = await ParseFilingDataSingular(response.data);
  return filing;
};

export const completeFileSchemaGet = async (
  url: string,
): Promise<CompleteFileSchema> => {
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
    author: file_schema.mdata.author,
    authors_information: file_schema.authors,
    item_number: file_schema.mdata.item_number,
    file_class: file_schema.mdata.file_class,
    url: file_schema.mdata.url,
  };
};
export const ParseFilingDataSingular = async (
  f: any,
): Promise<Filing | null> => {
  try {
    const completeFileSchema: CompleteFileSchema =
      CompleteFileSchemaValidator.parse(f);
    const newFiling: Filing = generateFilingFromFileSchema(completeFileSchema);
    return newFiling;
  } catch (error) {}

  try {
    console.log("Parsing document ID", f);
    console.log("filing source id", f.sourceID);
    const docID = z.string().uuid().parse(f.sourceID);
    const metadata_url = `${apiURL}/v2/public/files/${docID}/metadata`;
    try {
      const completeFileSchema = await completeFileSchemaGet(metadata_url);
      const newFiling: Filing =
        generateFilingFromFileSchema(completeFileSchema);
      return newFiling;
    } catch (error) {
      console.log("Error getting complete file schema", f, "error:", error);
      return null;
    }
  } catch (error) {
    console.log("Invalid document ID", f, "error:", error);
    return null;
  }
};

export const ParseFilingData = async (filingData: any): Promise<Filing[]> => {
  const filings_promises: Promise<Filing | null>[] = filingData.map(
    ParseFilingDataSingular,
  );
  const filings_with_errors = await Promise.all(filings_promises);
  const filings_null = filings_with_errors.filter(
    (f: Filing | null) => f !== null && f !== undefined,
  );
  const filings = filings_null as Filing[];
  return filings;
};

export default getSearchResults;
