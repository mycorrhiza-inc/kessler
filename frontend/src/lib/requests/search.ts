import { Filing } from "@/lib/types/FilingTypes";
import axios from "axios";
import {
  CompleteFileSchema,
  CompleteFileSchemaValidator,
} from "../types/backend_schemas";
import { getRuntimeEnv } from "../env_variables_hydration_script";

export const hydratedSearchResultsToFilings = (
  hydratedSearchResults: any | null,
): Filing[] => {
  const verifified_nullable_files: Array<CompleteFileSchema | null> =
    hydratedSearchResults
      ? hydratedSearchResults.map(
          (hydrated_result: any): CompleteFileSchema | null => {
            const maybe_file = hydrated_result.file;
            try {
              const file = CompleteFileSchemaValidator.parse(maybe_file);
              return file;
            } catch (error) {
              console.log("Error parsing file", error);
              return null;
            }
          },
        )
      : [];
  const valid_files: CompleteFileSchema[] = verifified_nullable_files.filter(
    (file) => file !== null,
  ) as CompleteFileSchema[]; // filter out null _files.filter
  const filings = valid_files.map(
    (file: CompleteFileSchema): Filing => generateFilingFromFileSchema(file),
  );
  return filings;
};

export const getRecentFilings = async (
  page?: number,
  page_size?: number,
): Promise<Filing[]> => {
  if (!page) {
    page = 0;
  }
  const default_page_size = 40;
  const queryString = queryStringFromPageMaxHits(
    page,
    page_size || default_page_size,
  );
  // Incorrect Code:
  // const default_page_size = 40;
  // const limit = page_size || default_page_size;
  // const queryString = queryStringFromPageMaxHits(limit, page_size);
  const runtimeConfig = getRuntimeEnv();
  const response = await axios.get(
    `${runtimeConfig.public_api_url}/v2/search/file/recent_updates${queryString}`,
  );
  if (response.status >= 400) {
    throw new Error(`Request failed with status code ${response.status}`);
  }
  // console.log("recent data", response.data);
  if (response.data.length > 0) {
    return hydratedSearchResultsToFilings(response.data);
  }
  throw new Error(
    "No recent filings found, their should absolutely be some files in the DB to show.",
  );
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
    docket_id: file_schema.mdata.docket_id,
    url: file_schema.mdata.url,
    extension: file_schema.extension,
  };
};
