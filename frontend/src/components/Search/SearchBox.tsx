import React, { useEffect, useRef, useState } from 'react';

// Mock API call
type Suggestion = {
	id: string;
	type: string;
	label: string;
	value: string;
};

type Filter = {
	id: string;
	type: string;
	label: string;
}

const mockFetchSuggestions = async (query: string): Promise<Suggestion[]> => {
	// Simulate API delay
	await new Promise(resolve => setTimeout(resolve, 300));

	const suggestions: Suggestion[] = [
		{ id: '1', type: 'organization', label: 'Acme Corp', value: 'acme' },
		{ id: '2', type: 'organization', label: 'Apple Inc', value: 'apple' },
		{ id: '3', type: 'case', label: 'Bug Report #123', value: 'bug-123' },
		{ id: '4', type: 'case', label: 'Feature Request #456', value: 'feature-456' }
	].filter(s =>
		s.label.toLowerCase().includes(query.toLowerCase()) ||
		s.type.toLowerCase().includes(query.toLowerCase())
	);

	return suggestions;
};
type FiltersPoolProps = {
	selected: Filter[];
	handleFilterRemove: (filterId: string) => void;
};

const FiltersPool: React.FC<FiltersPoolProps> = ({ selected, handleFilterRemove }) => {
	return (
		selected.length > 0 && (
			<div>
				<div className="divider pl-5 pr-5"></div>
				<div className="flex flex-wrap gap-2 p-2">
					{selected.map(filter => (
						<div
							key={filter.id}
							className="flex items-center gap-1 px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm"
						>
							<span className="font-medium">
								{filter.type}: {filter.label}
							</span>
							<button
								onClick={() => handleFilterRemove(filter.id)}
								className="ml-1 text-blue-600 hover:text-blue-800 font-bold"
							>
								×
							</button>
						</div>
					))}
				</div>

			</div>
		)
	)
}

const SearchBox = () => {
	const [query, setQuery] = useState('');
	const [suggestions, setSuggestions] = useState<Suggestion[]>([]);
	const [selectedFilters, setSelectedFilters] = useState<Filter[]>([]);
	const [isLoading, setIsLoading] = useState(false);
	const searchContainerRef = useRef<HTMLDivElement>(null);

	// Handle clicks outside of the search container
	useEffect(() => {
		const handleClickOutside = (event: any) => {
			// Check if the click was outside and we have suggestions open
			if (searchContainerRef.current &&
				!searchContainerRef.current.contains(event.target) &&
				suggestions.length > 0) {
				console.log('Click outside detected, closing suggestions');
				setSuggestions([]);
				setQuery('');
			}
		};

		// Use mousedown and touchstart for better mobile support
		document.addEventListener('click', handleClickOutside, true);
		document.addEventListener('touchstart', handleClickOutside, true);

		return () => {
			document.removeEventListener('click', handleClickOutside, true);
			document.removeEventListener('touchstart', handleClickOutside, true);
		};
	}, [suggestions.length]); // Add suggestions.length as dependency

	const handleInputChange = async (e: any) => {
		const newQuery = e.target.value;
		setQuery(newQuery);

		if (newQuery.trim()) {
			setIsLoading(true);
			const results = await mockFetchSuggestions(newQuery);
			setSuggestions(results);
			setIsLoading(false);
		} else {
			setSuggestions([]);
		}
	};

	const handleSuggestionClick = (suggestion: Suggestion) => {
		if (!selectedFilters.some(f => f.id === suggestion.id)) {
			setSelectedFilters([...selectedFilters, suggestion]);
		}
		setQuery('');
		setSuggestions([]);
	};

	const handleFilterRemove = (filterId: string) => {
		setSelectedFilters(selectedFilters.filter(f => f.id !== filterId));
	};


	return (
		<div className="p-4 max-w-xl mx-auto">
			<div className="flex flex-col gap-2">
				{/* Search container */}
				<div className="relative">
					{/* Search input */}
					<div className="relative">
						<input
							type="text"
							value={query}
							onChange={handleInputChange}
							placeholder="Search organizations or cases..."
							className="w-full p-3 border rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
						/>

						{isLoading && (
							<div className="absolute right-3 top-1/2 transform -translate-y-1/2">
								<div className="animate-spin h-5 w-5 border-2 border-blue-500 border-t-transparent rounded-full" />
							</div>
						)}
					</div>

					{/* Suggestions dropdown - Now positioned relative to search container */}
					{suggestions.length > 0 && (
						<div className="absolute left-0 right-0 top-full mt-1 z-50">
							<ul className="bg-white border rounded-lg shadow-lg max-h-60 overflow-auto">
								{suggestions.map(suggestion => (
									<li key={suggestion.id}>
										<button
											onClick={() => handleSuggestionClick(suggestion)}
											className="w-full px-4 py-3 text-left hover:bg-gray-50 focus:bg-gray-50 focus:outline-none transition-colors"
										>
											<span className="text-gray-500 text-sm font-medium">
												{suggestion.type}:
											</span>{' '}
											<span className="text-gray-900">
												{suggestion.label}
											</span>
										</button>
									</li>
								))}

								<li className='p-2'>
									<div></div>
									<FiltersPool selected={selectedFilters} handleFilterRemove={handleFilterRemove} /></li>
							</ul>
						</div>
					)}
				</div>

				<FiltersPool selected={selectedFilters} handleFilterRemove={handleFilterRemove} />
			</div>
		</div>
	);
};

export default SearchBox;