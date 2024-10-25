import React from "react";
import { Dispatch, SetStateAction, useState, useEffect, useRef } from "react";

export enum FilterField {
  MatchName = "match_name",
  MatchSource = "match_source",
  MatchDoctype = "match_doctype",
  MatchDocketId = "match_docket_id",
  MatchDocumentClass = "match_document_class",
  MatchAuthor = "match_author",
  MatchBeforeDate = "match_before_date",
  MatchAfterDate = "match_after_date",
}
export type QueryFilterFields = {
  [key in FilterField]: string;
};

enum InputType {
  Text = "text",
  Select = "select",
  Datalist = "datalist",
  Date = "date",
}
export const emptyQueryOptions: QueryFilterFields = {
  [FilterField.MatchName]: "",
  [FilterField.MatchSource]: "",
  [FilterField.MatchDoctype]: "",
  [FilterField.MatchDocketId]: "",
  [FilterField.MatchDocumentClass]: "",
  [FilterField.MatchAuthor]: "",
  [FilterField.MatchBeforeDate]: "",
  [FilterField.MatchAfterDate]: "",
};
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
function BasicDocumentFilters({
  queryOptions,
  setQueryOptions,
}: {
  queryOptions: QueryFilterFields;
  setQueryOptions: Dispatch<SetStateAction<QueryFilterFields>>;
}) {
  const DocumentFilter = ({
    filterData,
    filterID,
  }: {
    filterData: PropertyInformation;
    filterID: FilterField;
  }) => {
    const handleChange = (
      e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
    ) => {
      const value = e.target.value;
      setQueryOptions((prevOptions) => ({
        ...prevOptions,
        [filterID]: value,
      }));
    };
    switch (filterData.type) {
      case InputType.Text:
        return (
          <div className="box-border">
            <div className="tooltip" data-tip={filterData.description}>
              <p>{filterData.displayName}</p>
            </div>
            <input
              className="input input-bordered w-full max-w-xs"
              type="text"
              onChange={handleChange}
              title={filterData.displayName}
            />
          </div>
        );
      case InputType.Select:
        return (
          <div className="box-border">
            <div className="tooltip" data-tip={filterData.description}>
              <p>{filterData.displayName}</p>
            </div>
            <select
              className="select select-bordered w-full max-w-xs"
              onChange={handleChange}
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
            </div>
            <input
              className="input input-bordered w-full max-w-xs"
              type="date"
              onChange={handleChange}
              title={filterData.displayName}
            />
          </div>
        );
    }
  };

  const sortedFilters = Object.keys(FilterField)
    .map((key) => {
      const filterId = FilterField[key as keyof typeof FilterField];
      const filterData = queryFiltersInformation[filterId];
      const placementIndex = filterData.index;
      return { filterId, filterData, placementIndex };
    })
    .sort((a, b) => a.placementIndex - b.placementIndex);

  return (
    <>
      <div className="grid grid-cols-4 gap-4">
        {sortedFilters.map(({ filterId, filterData }) => (
          <DocumentFilter filterData={filterData} filterID={filterId} />
        ))}
      </div>
    </>
  );
}
export default BasicDocumentFilters;
