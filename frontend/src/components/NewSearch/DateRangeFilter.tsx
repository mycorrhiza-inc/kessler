import React from "react";

interface DateRangeFilterProps {
  dateFrom: string;
  dateTo: string;
  onDateFromChange: (date: string) => void;
  onDateToChange: (date: string) => void;
}

const DateRangeFilter: React.FC<DateRangeFilterProps> = ({
  dateFrom,
  dateTo,
  onDateFromChange,
  onDateToChange,
}) => {
  return (
    <div className="date-range-filter space-y-2">
      <h3 className="text-sm font-semibold">Date Range</h3>
      <div className="grid grid-cols-2 gap-2">
        <div>
          <label htmlFor="date-from" className="block text-xs mb-1">
            From
          </label>
          <input
            type="date"
            id="date-from"
            value={dateFrom}
            onChange={(e) => onDateFromChange(e.target.value)}
            className="w-full px-2 py-1 border rounded-md text-xs"
          />
        </div>
        <div>
          <label htmlFor="date-to" className="block text-xs mb-1">
            To
          </label>
          <input
            type="date"
            id="date-to"
            value={dateTo}
            onChange={(e) => onDateToChange(e.target.value)}
            className="w-full px-2 py-1 border rounded-md text-xs"
          />
        </div>
      </div>
    </div>
  );
};

export default DateRangeFilter;
