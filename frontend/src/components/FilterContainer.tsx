interface FilterProps {
	FilterName: string;
	FilterValue: string;
	FilterColor?: string; // Optional hex color of filter
}

const FilterContainer = ({ FilterName, FilterValue, FilterColor }: FilterProps) => {
	return (
		<div
			style={{
				display: "flex",
				flexDirection: "row",
				alignItems: "center",
				justifyContent: "space-between",
				padding: "5px",
				backgroundColor: FilterColor ? FilterColor : "#f0f0f0",
				borderRadius: "10px",
				margin: "5px",
			}}
		>
			<span>{FilterName}</span>
			<span>{FilterValue}</span>
		</div>
	);
};

export default FilterContainer;