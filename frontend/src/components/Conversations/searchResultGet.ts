import { QueryDataFile, QueryFilterFields } from "@/lib/filters";
import axios from "axios";

const searchResultsGet = async (queryData: QueryDataFile) => {
  const searchQuery = queryData.query;
  const searchFilters = queryData.filters;
  console.log(`searchhing for ${searchQuery}`);
  try {
    const response = await axios.post("https://api.kessler.xyz/v2/search", {
      query: searchQuery,
      filters: {
        name: searchFilters.match_name,
        author: searchFilters.match_author,
        docket_id: searchFilters.match_docket_id,
        doctype: searchFilters.match_doctype,
        source: searchFilters.match_source,
      },
    });
    if (response.data.length === 0) {
      return [];
    }
    if (typeof response.data === "string") {
      return [];
    }
    console.log("getting data");
    console.log(response.data);
    return response.data;
  } catch (error) {
    console.log(error);
  }
};
export default searchResultsGet;
