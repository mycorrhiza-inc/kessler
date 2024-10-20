import { Dispatch, SetStateAction, useState, useEffect, useRef } from "react";
import {
  extraProperties,
  extraPropertiesInformation,
  emptyExtraProperties,
} from "@/utils/interfaces";

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
  return (
    <>
      <div className="grid grid-cols-2 gap-4">
        {Object.keys(queryOptions)
          .slice(0, 6)
          .map((key, index) => {
            const extraInfo =
              extraPropertiesInformation[key as keyof extraProperties];
            return (
              <div className="box-border" key={index}>
                <div className="tooltip" data-tip={extraInfo.description}>
                  <p>{extraInfo.displayName}</p>
                </div>
                <input
                  className="input input-bordered w-full max-w-xs"
                  type="text"
                  id={key}
                  name={key}
                  value={queryOptions[key as keyof extraProperties]}
                  onChange={handleChange}
                  title={extraInfo.displayName}
                />
              </div>
            );
          })}
      </div>
    </>
  );
}
export default BasicDocumentFilters;
