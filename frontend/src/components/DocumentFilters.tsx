import { Dispatch, SetStateAction, useState, useEffect, useRef } from "react";
import { extraProperties, emptyExtraProperties } from "@/utils/interfaces";

import React from "react";
import Select from "react-select";
import { Calendar } from "react-date-range";

enum InputType {
  Text = "text",
  Select = "select",
  Date = "date",
}
type PropertyInformation = {
  type: InputType;
  displayName: string;
  description: string;
  details: string;
  options?: { label: string; value: string }[];
  isDate?: boolean;
};
const extraPropertiesInformation: ExtraPropertiesInformation = {
  match_name: {
    type: InputType.Text,
    displayName: "Name",
    description: "The name associated with the search item.",
    details: "Searches for items approximately matching the title",
  },
  match_source: {
    type: InputType.Select,
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
  match_docket_id: {
    type: InputType.Text,
    displayName: "Docket ID",
    description: "The unique identifier for the docket.",
    details: "Filters search results based on the docket ID.",
  },
  match_document_class: {
    type: InputType.Select,
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
  match_author: {
    type: InputType.Text,
    displayName: "Author",
    description: "The author of the document.",
    details: "Searches for items created or written by the specified author.",
  },
  match_date: {
    type: InputType.Date,
    displayName: "Date",
    description: "The date related to the document.",
    details: "Filters results by the specified date.",
  },
};
type ExtraPropertiesInformation = {
  [key: string]: PropertyInformation;
};
function BasicDocumentFilters({
  queryOptions,
  setQueryOptions,
}: {
  queryOptions: extraProperties;
  setQueryOptions: Dispatch<SetStateAction<extraProperties>>;
}) {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setQueryOptions((prevOptions) => ({
      ...prevOptions,
      [name]: value,
    }));
  };

  const DocumentFilter = ({
    filterData,
  }: {
    filterData: PropertyInformation;
  }) => {
    switch (filterData.type) {
      case InputType.Text:
        return 3;
      case InputType.Select:
        return 4;
      case InputType.Date:
        return 5;
    }
  };

  return (
    <>
      <div className="grid grid-cols-2 gap-4">
        {Object.keys(queryOptions)
          .slice(0, 7)
          .map((key, index) => {
            const filterData = extraPropertiesInformation[key];
            return (
              <div className="box-border" key={index}>
                <div className="tooltip" data-tip={filterData.description}>
                  <p>{filterData.displayName}</p>
                </div>
                {filterData.options ? (
                  <Select
                    className="select select-bordered w-full max-w-xs"
                    id={key}
                    name={key}
                    options={filterData.options}
                    value={filterData.options.find(
                      (option) =>
                        option.value ===
                        queryOptions[key as keyof extraProperties],
                    )}
                    onChange={(selectedOption) =>
                      handleChange({
                        // @ts-ignore
                        target: { name: key, value: selectedOption.value },
                      })
                    }
                    // title={filterData.displayName}
                  />
                ) : filterData.isDate ? (
                  <Calendar id={key} name={key} />
                ) : (
                  <input
                    className="input input-bordered w-full max-w-xs"
                    type="text"
                    id={key}
                    name={key}
                    value={queryOptions[key as keyof extraProperties]}
                    onChange={handleChange}
                    title={filterData.displayName}
                  />
                )}
              </div>
            );
          })}
      </div>
    </>
  );
}
export default BasicDocumentFilters;
