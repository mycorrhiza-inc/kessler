import { QueryDataFile, QueryFilterFields } from "@/lib/filters";
import { Filing } from "@/lib/types/FilingTypes";
import axios from "axios";
const getSearchResults = async (queryData: QueryDataFile) => {
  const searchQuery = queryData.query;
  const searchFilters = queryData.filters;
  console.log(`searchhing for ${searchQuery}`);
  try {
    const response = await axios
      // .post("https://api.kessler.xyz/v2/search", {
      .post("http://localhost/v2/search", {
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
        return response.data;
      })
      // coerce Filing typ
      .then((data) => {
        const filings = ParseFilingData(data);
        return filings;
      });
    console.log("getting data");
    console.log(response);
    return response.data;
  } catch (error) {
    console.log(error);
  }
};

export const getRecentFilings = async (page?: number) => {
  if (!page) {
    page = 0;
  }
  const response = await axios.post(
    "http://localhost/v2/recent_updates",
    // "http://api.kessler.xyz/v2/recent_updates",
    {
      page: page,
    }
  );
  console.log(response.data);
  if (response.data.length > 0) {
    return response.data;
  }
};

export const getFilingMetadata = async (id: string) => {
  const response = await axios.get(
    // `http://api.kessler.xyz/v2/public/files/${id}/metadata`
    `http://localhost/v2/public/files/${id}/metadata`
  );
  return ParseFilingData([response.data]);
};

export const ParseFilingData = (data: any) => {
  const filings = data.map((f: any) => {
    const newFiling: Filing = {
      id: data.sourceID,
      lang: data.metadata.lang,
      title: data.title,
      date: data.medata.date,
      author: data.medata.author,
      source: data.metadata.source,
      language: data.metadata.language,
      item_number: data.medata.item_number,
      author_organisation: data.medata.author_organizatino,
      url: data.medata.url,
    };
    return newFiling;
  });
  return filings;
};

export default getSearchResults;
