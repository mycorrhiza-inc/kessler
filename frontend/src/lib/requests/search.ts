import { Filing } from "@/lib/types/FilingTypes";
import axios from "axios";
import {
  CompleteFileSchema,
  CompleteFileSchemaValidator,
} from "../types/backend_schemas";

export const hydratedSearchResultsToFilings = (
  hydratedSearchResults: any | null,
): Filing[] => {
  if (!hydratedSearchResults) {
    console.log("Search results from server was undefined");
    throw new Error("Got undefined search results from server");
    return [];
  }
  const verifified_nullable_files: Array<CompleteFileSchema | null> =
    hydratedSearchResults.map(
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
    );
  const valid_files: CompleteFileSchema[] = verifified_nullable_files.filter(
    (file) => file !== null,
  ) as CompleteFileSchema[]; // filter out null _files.filter
  const filings = valid_files.map(
    (file: CompleteFileSchema): Filing => generateFilingFromFileSchema(file),
  );
  return filings;
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

