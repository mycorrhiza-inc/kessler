import React from "react";
import { Dispatch, SetStateAction, useMemo } from "react";
import { FilterField, QueryFilterFields } from "@/lib/filters";
import clsx from "clsx";
import {
  InputType,
  PropertyInformation,
  queryFiltersInformation,
} from "./FiltersInfo";
import { OrgMultiSelect, ConvoMultiSelect } from "./FilterMultiSelect";

export function BasicDocumentFiltersList({
  queryOptions,
  setQueryOptions,
  showQueries,
  disabledQueries,
}: {
  queryOptions: QueryFilterFields;
  setQueryOptions: Dispatch<SetStateAction<QueryFilterFields>>;
  showQueries: FilterField[];
  disabledQueries?: FilterField[];
}) {
  return (
    <BasicDocumentFilters
      className="grid grid-flow-row auto-rows-max gap-4"
      queryOptions={queryOptions}
      setQueryOptions={setQueryOptions}
      showQueries={showQueries}
      disabledQueries={disabledQueries}
    />
  );
}
export function BasicDocumentFiltersGrid({
  queryOptions,
  setQueryOptions,
  showQueries,
  disabledQueries,
}: {
  queryOptions: QueryFilterFields;
  setQueryOptions: Dispatch<SetStateAction<QueryFilterFields>>;
  showQueries: FilterField[];
  disabledQueries?: FilterField[];
}) {
  return (
    <BasicDocumentFilters
      className="grid grid-cols-4 gap-4"
      queryOptions={queryOptions}
      setQueryOptions={setQueryOptions}
      showQueries={showQueries}
      disabledQueries={disabledQueries}
    />
  );
}
export function BasicDocumentFilters({
  className,
  queryOptions,
  setQueryOptions,
  showQueries,
  disabledQueries,
  max_w_xs,
}: {
  className?: string;
  queryOptions: QueryFilterFields;
  setQueryOptions: Dispatch<SetStateAction<QueryFilterFields>>;
  showQueries: FilterField[];
  disabledQueries?: FilterField[];
  max_w_xs?: boolean;
}) {
  // const [docFilterValues, setDocFilterValues] = useState(emptyQueryOptions);
  const disabledQueriesDict = useMemo(() => {
    const dict: Partial<Record<FilterField, boolean>> = {};
    if (disabledQueries) {
      disabledQueries.forEach((query) => {
        dict[query] = true;
      });
    }
    return dict;
  }, [disabledQueries]);
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

  const max_w_xs_ClassString = max_w_xs ? "max-w-xs" : "";
  const generateDocumentFilter = (
    filterId: FilterField,
    filterData: PropertyInformation,
  ) => {
    const isDisabled = disabledQueriesDict[filterId] || false;
    if (isDisabled) {
      return (
        <div className="box-border">
          <div className="tooltip" data-tip={filterData.description}>
            <p>{filterData.displayName}</p>
          </div>
          <br />
          <input
            className={clsx(
              "input input-bordered w-full",
              max_w_xs_ClassString,
            )}
            type="text"
            disabled={isDisabled}
            value={docFilterValues[filterId]}
            onChange={(e) => handleChange(e, filterId)}
            title={filterData.displayName}
          />
        </div>
      );
    }
    switch (filterData.type) {
      case InputType.Text:
        return (
          <div className="box-border">
            <div className="tooltip" data-tip={filterData.description}>
              <p>{filterData.displayName}</p>
            </div>
            <br />
            <input
              className={clsx(
                "input input-bordered w-full",
                max_w_xs_ClassString,
              )}
              type="text"
              disabled={isDisabled}
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
            </div>
            <br />
            <select
              disabled={isDisabled}
              className={clsx(
                "select select-bordered w-full",
                max_w_xs_ClassString,
              )}
              value={docFilterValues[filterId]}
              onChange={(e) => handleChange(e, filterId)}
            >
              {filterData.options?.map((option, index) => (
                <option
                  key={option.value}
                  value={option.value}
                  selected={index === 0}
                >
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
            <br />
            <input
              className={clsx(
                "input input-bordered w-full",
                max_w_xs_ClassString,
              )}
              type="date"
              disabled={isDisabled}
              value={docFilterValues[filterId]}
              onChange={(e) => handleChange(e, filterId)}
              title={filterData.displayName}
            />
          </div>
        );
      case InputType.OrgMultiselect:
        return (
          <div className="box-border">
            <div className="tooltip" data-tip={filterData.description}>
              <p>{filterData.displayName}</p>
            </div>
            <br />
            <OrgMultiSelect />
          </div>
        );
      case InputType.ConvoMultiselect:
        return (
          <div className="box-border">
            <div className="tooltip" data-tip={filterData.description}>
              <p>{filterData.displayName}</p>
            </div>
            <br />
            <ConvoMultiSelect />
          </div>
        );
      case InputType.NotShown:
        return <></>;
    }
  };

  return (
    <div className={className}>
      {sortedFilters.map(({ filterId, filterData }) =>
        generateDocumentFilter(filterId, filterData),
      )}
    </div>
  );
}
export default BasicDocumentFilters;
