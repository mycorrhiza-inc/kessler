import { ChevronDownIcon, ChevronUpIcon } from '@/components/Icons';
import { data } from 'autoprefixer';
import { range } from 'lodash-es';
import React, { useState } from 'react';
import { DateRangePicker, RangeKeyDict } from 'react-date-range';
import 'react-date-range/dist/styles.css'; // main style file
import 'react-date-range/dist/theme/default.css'; // theme css file

declare global {
	namespace JSX {
		interface IntrinsicElements {
			'calendar-date': React.DetailedHTMLProps<React.HTMLAttributes<HTMLElement>, HTMLElement>;
			'calendar-month': React.DetailedHTMLProps<React.HTMLAttributes<HTMLElement>, HTMLElement>;
			'calendar-range': React.DetailedHTMLProps<React.HTMLAttributes<HTMLElement>, HTMLElement>;
		}
	}
}


type DateRangePickerProps = {
	startPick: Date | null;
	endPick: Date | null;
	setStartPick: (date: Date) => void;
	setEndPick: (date: Date) => void;
}
const RangePicker = ({ startPick, endPick, setStartPick, setEndPick }: DateRangePickerProps) => {
	const [open, setOpen] = useState(false);
	const flip = () => {
		setOpen((prev) => !prev)
	}

	const selectionRange = {
		startDate: new Date(),
		endDate: new Date(),
		key: 'selection',
	}

	const handleChange = (rangesByKey: RangeKeyDict) => {
		const range = rangesByKey[0];
		if (range.startDate) { setStartPick(range.startDate) }
		if (range.endDate) { setEndPick(range.endDate) }
	}
	return (
		<div className='flex flex-row'>
			<details className="dropdown">
				<summary className='' onClick={flip}>
					<u>{startPick !== null ? startPick.toLocaleDateString("en-US") : "?"}</u> to
					<u>{endPick !== null ? endPick.toLocaleDateString("en-US") : "?"}</u>
					{open ? <ChevronDownIcon /> : <ChevronUpIcon />}
				</summary>
				<div className='dropdown-content'>
					<DateRangePicker
						ranges={[selectionRange]}
						onChange={handleChange}
					/>
				</div>
			</details>
		</div>
	)
}