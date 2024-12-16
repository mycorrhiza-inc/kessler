import {
  BackendFilterObject,
  QueryDataFile,
  backendFilterGenerate,
} from "@/lib/filters";
import { Filing } from "@/lib/types/FilingTypes";
import axios from "axios";
import { publicAPIURL } from "../env_variables";
import {
  CompleteFileSchema,
  CompleteFileSchemaValidator,
} from "../types/backend_schemas";
import { queryStringFromPageMaxHits } from "../pagination";

export const getSearchResults = async (
  queryData: QueryDataFile,
  page: number,
  maxHits: number,
): Promise<Filing[]> => {
  const searchQuery = queryData.query;
  console.log("query data", queryData);
  const searchFilters = queryData.filters;
  console.log("searchhing for", searchFilters);
  const filterObj: BackendFilterObject = backendFilterGenerate(searchFilters);
  try {
    const paginationQueryString = queryStringFromPageMaxHits(page, maxHits);
    const searchResults: Filing[] = await axios
      // .post("https://api.kessler.xyz/v2/search", {
      .post(`${publicAPIURL}/v2/search${paginationQueryString}`, {
        query: searchQuery,
        filters: filterObj,
      })
      // check error conditions
      .then((response) => {
        if (response.data?.length === 0 || typeof response.data === "string") {
          return [];
        }
        const filings = hydratedSearchResultsToFilings(response.data);
        console.log("response data:::::", response.data);
        return filings;
      });
    console.log("getting data");
    // console.log(searchResults);
    return searchResults;
  } catch (error) {
    console.log(error);
    throw error;
  }
};

export const hydratedSearchResultsToFilings = (
  hydratedSearchResults: any,
): Filing[] => {
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

export const getRecentFilings = async (
  page?: number,
  page_size?: number,
): Promise<Filing[]> => {
  if (!page) {
    page = 0;
  }
  const default_page_size = 40;
  const limit = page_size || default_page_size;
  const queryString = queryStringFromPageMaxHits(limit, page_size);
  const response = await axios.get(
    `${publicAPIURL}/v2/recent_updates${queryString}`,
  );
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
  };
};

// export const getFilingMetadata = async (id: string): Promise<Filing | null> => {
//   const valid_id = z.string().uuid().parse(id);
//   const response = await axios.get(
//     `${publicAPIURL}/v2/public/files/${valid_id}`,
//   );
//   const filing = await ParseFilingDataSingular(response.data);
//   return filing;
// };
// export const ParseFilingDataSingular = async (
//   f: any,
// ): Promise<Filing | null> => {
//   try {
//     const completeFileSchema: CompleteFileSchema =
//       CompleteFileSchemaValidator.parse(f);
//     const newFiling: Filing = generateFilingFromFileSchema(completeFileSchema);
//     return newFiling;
//   } catch (error) {}
//
//   try {
//     console.log("Parsing document ID", f);
//     console.log("filing source id", f.sourceID);
//     const docID = z.string().uuid().parse(f.sourceID);
//     const metadata_url = `${publicAPIURL}/v2/public/files/${docID}`;
//     try {
//       const completeFileSchema = await completeFileSchemaGet(metadata_url);
//       const newFiling: Filing =
//         generateFilingFromFileSchema(completeFileSchema);
//       return newFiling;
//     } catch (error) {
//       console.log("Error getting complete file schema", f, "error:", error);
//       return null;
//     }
//   } catch (error) {
//     console.log("Invalid document ID", f, "error:", error);
//     return null;
//   }
// };
//
// export const ParseFilingData = async (filingData: any): Promise<Filing[]> => {
//   const out: Filing[] = [];
//   if (filingData === null) {
//     return out;
//   }
//   const filings_promises: Promise<Filing | null>[] = filingData.map(
//     ParseFilingDataSingular,
//   );
//   const filings_with_errors = await Promise.all(filings_promises);
//   const filings_null = filings_with_errors.filter(
//     (f: Filing | null) => f !== null && f !== undefined,
//   );
//   const filings = filings_null as Filing[];
//   return filings;
// };

export default getSearchResults;
