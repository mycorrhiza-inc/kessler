import { ChevronDownIcon, ChevronUpIcon } from '@/components/Icons';
import { data } from 'autoprefixer';
import { endOfDay } from 'date-fns';
import { range } from 'lodash-es';
import React, { useState } from 'react';
import { DateRangePicker, RangeKeyDict, Range as DateRange } from 'react-date-range';
import 'react-date-range/dist/styles.css'; // main style file
import 'react-date-range/dist/theme/default.css'; // theme css file

type DateRangePickerProps = {
	baseRange: DateRange
	updateRange: (newRange: DateRange) => void
}
export const RangePicker = ({ baseRange, updateRange }: DateRangePickerProps) => {
	const [open, setOpen] = useState(false);
	const flip = () => {
		setOpen((prev) => !prev)
	}

	const selectionRange = {
		startDate: baseRange.startDate !== null ? baseRange.startDate : new Date(),
		endDate: baseRange.endDate !== null ? baseRange.endDate : new Date(),
		key: 'selection',
	}

	const handleChange = (rangesByKey: RangeKeyDict) => {
		const newRange = Object.values(rangesByKey)[0];
		console.log("!@#!@#!@#", rangesByKey)
		 { updateRange({ ...baseRange, startDate: newRange.startDate, endDate: newRange.endDate }) }
	}
	return (
		<div className='flex flex-row'>
			<details className="dropdown">
				<summary className='' onClick={flip}>
					Date from <u>{baseRange.startDate !== undefined ? baseRange.startDate.toLocaleDateString("en-US") : "?"}</u> to
					<u>{baseRange.endDate !== undefined ? baseRange.endDate.toLocaleDateString("en-US") : "?"}</u>
				</summary>
				<div className='dropdown-content  z-9910'>
					<div className=" bg-base-100 rounded-lg">
					<DateRangePicker
						ranges={[selectionRange]}
						onChange={handleChange}
					/>
					</div>
				</div>
			</details>
		</div>
	)
}