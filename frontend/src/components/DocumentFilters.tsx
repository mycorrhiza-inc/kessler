import React from "react";
import {
  Dispatch,
  SetStateAction,
  useState,
  useEffect,
  useRef,
  useMemo,
} from "react";
import {
  FilterField,
  QueryFilterFields,
  emptyQueryOptions,
} from "@/lib/filters";

// export enum FilterField {
//   MatchName = "match_name",
//   MatchSource = "match_source",
//   MatchDoctype = "match_doctype",
//   MatchDocketId = "match_docket_id",
//   MatchDocumentClass = "match_document_class",
//   MatchAuthor = "match_author",
//   MatchBeforeDate = "match_before_date",
//   MatchAfterDate = "match_after_date",
// }
// export type QueryFilterFields = {
//   [key in FilterField]: string;
// };
//
// export const emptyQueryOptions: QueryFilterFields = {
//   [FilterField.MatchName]: "",
//   [FilterField.MatchSource]: "",
//   [FilterField.MatchDoctype]: "",
//   [FilterField.MatchDocketId]: "",
//   [FilterField.MatchDocumentClass]: "",
//   [FilterField.MatchAuthor]: "",
//   [FilterField.MatchBeforeDate]: "",
//   [FilterField.MatchAfterDate]: "",
// };
enum InputType {
  Text = "text",
  Select = "select",
  Datalist = "datalist",
  Date = "date",
}
type PropertyInformation = {
  type: InputType;
  index: number;
  displayName: string;
  description: string;
  details: string;
  options?: { label: string; value: string }[];
  isDate?: boolean;
};
type QueryFiltersInformation = {
  [key in FilterField]: PropertyInformation;
};
const queryFiltersInformation: QueryFiltersInformation = {
  match_name: {
    type: InputType.Text,
    index: 1,
    displayName: "Name",
    description: "The name associated with the search item.",
    details: "Searches for items approximately matching the title",
  },
  match_docket_id: {
    type: InputType.Text,
    index: 2,
    displayName: "Docket ID",
    description: "The unique identifier for the docket.",
    details: "Filters search results based on the docket ID.",
  },
  match_author: {
    type: InputType.Text,
    index: 3,
    displayName: "Author",
    description: "The author of the document.",
    details: "Searches for items created or written by the specified author.",
  },
  match_source: {
    type: InputType.Select,
    index: 4,
    displayName: "State",
    description: "What state do you want to limit your search to?",
    details: "",
    options: [
      { label: "All", value: "" },
      { label: "New York", value: "usa-ny" },
      { label: "Colorado", value: "usa-co" },
      { label: "California", value: "usa-ca" },
      { label: "National", value: "usa-national" },
    ],
  },
  match_doctype: {
    type: InputType.Select,
    index: 5,
    displayName: "Document File Type",
    description: "The type or category of the document.",
    details: "Searches for items that match the specified document type.",
    options: [
      { label: "All", value: "" },
      { label: "PDF", value: "pdf" },
      { label: "Docx", value: "docx" },
      { label: "Xls", value: "xls" },
    ],
  },
  match_document_class: {
    type: InputType.Select,
    index: 6,
    displayName: "Document Class",
    description: "The classification or category of the document.",
    details: "Searches for documents that fall under the specified class.",
    options: [
      { label: "All", value: "" },
      { label: "Correspondence", value: "Correspondence" },
      { label: "Comments", value: "Comments" },
      { label: "Reports", value: "Reports" },
      { label: "Plans and Proposals", value: "Plans and Proposals" },
      { label: "Motions", value: "Motions" },
      { label: "Letters", value: "Letters" },
      { label: "Orders", value: "Orders" },
      { label: "Notices", value: "Notices" },
    ],
  },
  match_before_date: {
    type: InputType.Date,
    index: 7,
    displayName: "Before Date",
    description: "The date related to the document.",
    details: "Filters results by the specified date.",
  },
  match_after_date: {
    type: InputType.Date,
    index: 8,
    displayName: "After Date",
    description: "The date related to the document.",
    details: "Filters results by the specified date.",
  },
};
export function BasicDocumentFiltersList({
  queryOptions,
  setQueryOptions,
  showQueries,
}: {
  queryOptions: QueryFilterFields;
  setQueryOptions: Dispatch<SetStateAction<QueryFilterFields>>;
  showQueries: FilterField[];
}) {
  return <>
      <div className="grid grid-flow-row auto-rows-max gap-4">
      <BasicDocumentFilters
        queryOptions={queryOptions}
        setQueryOptions={setQueryOptions}
        showQueries={showQueries}
      />
      </div>
  </>
}
export function BasicDocumentFiltersGrid({
  queryOptions,
  setQueryOptions,
  showQueries,
}: {
  queryOptions: QueryFilterFields;
  setQueryOptions: Dispatch<SetStateAction<QueryFilterFields>>;
  showQueries: FilterField[];
}) {
  return <>
      <div className="grid grid-cols-4 gap-4">
      <BasicDocumentFilters
        queryOptions={queryOptions}
        setQueryOptions={setQueryOptions}
        showQueries={showQueries}
      />
      </div>
  </>
}
export function BasicDocumentFilters({
  queryOptions,
  setQueryOptions,
  showQueries,
}: {
  queryOptions: QueryFilterFields;
  setQueryOptions: Dispatch<SetStateAction<QueryFilterFields>>;
  showQueries: FilterField[];
}) {
  // const [docFilterValues, setDocFilterValues] = useState(emptyQueryOptions);
  const docFilterValues = queryOptions;
  const setDocFilterValues = setQueryOptions;
  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
    filterID: FilterField,
  ) => {
    setDocFilterValues((prevOptions) => ({
      ...prevOptions,
      [filterID]: e.target.value,
    }));
  };

  const sortedFilters = useMemo(() => {
    return showQueries
      .map((filterId) => {
        const filterData = queryFiltersInformation[filterId];
        const placementIndex = filterData.index;
        return { filterId, filterData, placementIndex };
      })
      .sort((a, b) => a.placementIndex - b.placementIndex);
  }, [showQueries, queryFiltersInformation]);

  return (
    <>
        {sortedFilters.map(({ filterId, filterData }) => {
          switch (filterData.type) {
            case InputType.Text:
              return (
                <div className="box-border">
                  <div className="tooltip" data-tip={filterData.description}>
                    <p>{filterData.displayName}</p>
                  </div><br/>
                  <input
                    className="input input-bordered w-full max-w-xs"
                    type="text"
                    value={docFilterValues[filterId]}
                    onChange={(e) => handleChange(e, filterId)}
                    title={filterData.displayName}
                  />
                </div>
              );
            case InputType.Select:
              return (
                <div className="box-border">
                  <div className="tooltip" data-tip={filterData.description}>
                    <p>{filterData.displayName}</p>
                  </div><br/>
                  <select
                    className="select select-bordered w-full max-w-xs"
                    value={docFilterValues[filterId]}
                    onChange={(e) => handleChange(e, filterId)}
                  >
                    {filterData.options?.map((option, index) => (
                      <option key={option.value} selected={index === 0}>
                        {option.label}
                      </option>
                    ))}
                  </select>
                </div>
              );
            case InputType.Date:
              return (
                <div className="box-border">
                  <div className="tooltip" data-tip={filterData.description}>
                    <p>{filterData.displayName}</p>
                  </div><br/>
                  <input
                    className="input input-bordered w-full max-w-xs"
                    type="date"
                    value={docFilterValues[filterId]}
                    onChange={(e) => handleChange(e, filterId)}
                    title={filterData.displayName}
                  />
                </div>
              );
          }
        })}
    </>
  );
}
export default BasicDocumentFilters;
