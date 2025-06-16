import { FaSearch } from "react-icons/fa";


// take in an extra param searchExecute, and run that closure if enter is pressed (but not if the user is also holding down shift)
export function SearchBar({
  value,
  setQuery: onChange,
  placeholder,
  searchExecute,
}: {
  value: string;
  setQuery: (value: string) => void;
  placeholder?: string;
  searchExecute: () => void;
}) {
  return (
    <div className="relative">
      <input
        type="text"
        placeholder={placeholder}
        className="w-full px-4 py-3 pr-10 border border-gray-300 rounded-lg focus:outline-hidden focus:ring-2 focus:ring-primary focus:border-transparent"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === 'Enter' && !e.shiftKey) {
            searchExecute();
          }
        }}
      />
      <FaSearch className="h-6 w-6 text-gray-400 absolute right-3 top-1/2 transform -translate-y-1/2" />
    </div>
  );
}
