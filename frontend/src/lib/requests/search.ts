import { QueryDataFile, QueryFilterFields } from "@/lib/filters";
import { Filing } from "@/lib/types/FilingTypes";
import axios from "axios";

export const getSearchResults = async (queryData: QueryDataFile) => {
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
    `http://localhost/v2/public/files/${id}/metadata`,
  );
  return ParseFilingData([response.data])[0];
};

export const ParseFilingData = (filingData: any) => {
  const filings = filingData.map((f: any) => {
    console.log("fdata", f);
    const metadata = JSON.parse(atob(f.Mdata));
    console.log("metadata: ", metadata);
    f.metadata = metadata;
    const newFiling: Filing = {
      id: f.ID,
      lang: f.metadata.lang,
      title: f.Name,
      date: f.metadata.date,
      author: f.metadata.author,
      source: f.metadata.source,
      language: f.metadata.language,
      item_number: f.metadata.item_number,
      author_organisation: f.metadata.author_organizatino,
      url: f.metadata.url,
    };
    return newFiling;
  });
  return filings;
};

export default getSearchResults;
