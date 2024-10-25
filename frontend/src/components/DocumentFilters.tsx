import { Dispatch, SetStateAction, useState, useEffect, useRef } from "react";
import { extraProperties, emptyExtraProperties } from "@/utils/interfaces";

import React from "react";
import Select from "react-select";
import { Calendar } from "react-date-range";

type ExtraPropertiesInformation = {
  [key: string]: {
    displayName: string;
    description: string;
    details: string;
    options?: { label: string; value: string }[];
    isDate?: boolean;
  };
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

  const extraPropertiesInformation: ExtraPropertiesInformation = {
    match_name: {
      displayName: "Name",
      description: "The name associated with the search item.",
      details: "Searches for items approximately matching the title",
    },
    match_source: {
      displayName: "State",
      description: "What state do you want to limit your search to?",
      details: "",
      options: [
        { label: "New York", value: "usa-ny" },
        { label: "Colorado", value: "usa-co" },
        { label: "California", value: "usa-ca" },
        { label: "National", value: "usa-national" },
      ],
    },
    match_doctype: {
      displayName: "Document File Type",
      description: "The type or category of the document.",
      details: "Searches for items that match the specified document type.",
      options: [
        { label: "PDF", value: "pdf" },
        { label: "Docx", value: "docx" },
        { label: "Xls", value: "xls" },
      ],
    },
    match_docket_id: {
      displayName: "Docket ID",
      description: "The unique identifier for the docket.",
      details: "Filters search results based on the docket ID.",
    },
    match_document_class: {
      displayName: "Document Class",
      description: "The classification or category of the document.",
      details: "Searches for documents that fall under the specified class.",
      options: [
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
      displayName: "Author",
      description: "The author of the document.",
      details: "Searches for items created or written by the specified author.",
    },
    match_date: {
      displayName: "Date",
      description: "The date related to the document.",
      details: "Filters results by the specified date.",
      isDate: true,
    },
  };

  return (
    <>
      <div className="grid grid-cols-2 gap-4">
        {Object.keys(queryOptions)
          .slice(0, 7)
          .map((key, index) => {
            const extraInfo = extraPropertiesInformation[key];
            return (
              <div className="box-border" key={index}>
                <div className="tooltip" data-tip={extraInfo.description}>
                  <p>{extraInfo.displayName}</p>
                </div>
                {extraInfo.options ? (
                  <Select
                    className="select select-bordered w-full max-w-xs"
                    id={key}
                    name={key}
                    options={extraInfo.options}
                    value={extraInfo.options.find(
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
                    // title={extraInfo.displayName}
                  />
                ) : extraInfo.isDate ? (
                  <Calendar id={key} name={key} />
                ) : (
                  <input
                    className="input input-bordered w-full max-w-xs"
                    type="text"
                    id={key}
                    name={key}
                    value={queryOptions[key as keyof extraProperties]}
                    onChange={handleChange}
                    title={extraInfo.displayName}
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
