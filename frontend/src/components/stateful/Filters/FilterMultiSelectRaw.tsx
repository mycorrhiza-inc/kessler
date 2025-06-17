import React from 'react';
import { TextPill } from '@/components/style/Pills/TextPills';

// Shared definitions
export interface FilterFieldDefinition {
  id: string;
  displayName: string;
  description: string;
  placeholder?: string;
  options?: Array<{
    value: string;
    label: string;
    disabled?: boolean;
  }>;
}

export interface Option {
  value: string;
  label: string;
  disabled?: boolean;
}

export interface FilterMultiSelectRawProps {
  fieldDefinition: FilterFieldDefinition;
  selectedValues: string[];
  selectedLabels: string[];
  filteredOptions: Option[];
  isOpen: boolean;
  searchTerm: string;
  disabled?: boolean;
  className?: string;
  onToggleDropdown: () => void;
  onSearchChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onClearSearch: () => void;
  onToggleOption: (value: string) => void;
  onRemoveSelectedItem: (value: string) => void;
  onClearAll: () => void;
  onFocus?: () => void;
  onBlur?: () => void;
}

export function FilterMultiSelectRaw({
  fieldDefinition,
  selectedValues,
  selectedLabels,
  filteredOptions,
  isOpen,
  searchTerm,
  disabled = false,
  className = '',
  onToggleDropdown,
  onSearchChange,
  onClearSearch,
  onToggleOption,
  onRemoveSelectedItem,
  onClearAll,
}: FilterMultiSelectRawProps) {
  return (
    <div className={`relative w-full multiselect-container ${className}`}>      
      <button
        type="button"
        className={
          `w-full min-h-[3rem] p-3 text-left border-2 rounded-lg transition-all duration-200
        ${isOpen ? 'border-blue-500 ring-2 ring-blue-200' : 'border-gray-300'}
        ${disabled ? 'bg-gray-100 cursor-not-allowed' : 'bg-base-100 hover:border-gray-400 cursor-pointer'}
        ${selectedValues.length > 0 ? 'bg-base-100' : ''}`
        }
        onClick={onToggleDropdown}
        disabled={disabled}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            onToggleDropdown();
          }
        }}
      >
        <div className="flex items-center justify-between">
          <div className="flex-1 min-w-0">
            {selectedValues.length === 0 ? (
              <span className="text-gray-500">
                {fieldDefinition.placeholder || 'Select options...'}
              </span>
            ) : (
              <div className="space-y-2">
                <div className="font-medium text-sm text-gray-700">
                  {fieldDefinition.displayName}: {selectedValues.length} selected
                </div>
                <div className="flex flex-wrap gap-1">
                  {selectedLabels.map((label, idx) => (
                    <TextPill
                      key={selectedValues[idx]}
                      text={label}
                      seed={selectedValues[idx]}
                      removable
                      onRemove={() => onRemoveSelectedItem(selectedValues[idx])}
                      size="xs"
                      className="m-0 mt-0 mb-0"
                    />
                  ))}
                </div>
              </div>
            )}
          </div>
          <svg
            className={`w-5 h-5 text-gray-400 transition-transform flex-shrink-0 ml-2 ${
              isOpen ? 'rotate-180' : ''
            }`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      </button>

      {isOpen && (
        <div className="absolute top-full left-0 right-0 z-50 mt-1 bg-base-100 border-2 border-gray-200 rounded-lg shadow-lg max-h-80 overflow-hidden">
          <div className="p-3 border-b border-gray-200 bg-base-100">
            <div className="relative">
              <input
                type="text"
                className="w-full pl-9 pr-9 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder={`Search ${fieldDefinition.displayName.toLowerCase()}...`}
                value={searchTerm}
                onChange={onSearchChange}
                onClick={(e) => e.stopPropagation()}
              />
              <svg
                className="w-4 h-4 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              {searchTerm && (
                <span
                  className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600 cursor-pointer"
                  onClick={onClearSearch}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter' || e.key === ' ') {
                      e.preventDefault();
                      onClearSearch();
                    }
                  }}
                  role="button"
                  tabIndex={0}
                  aria-label="Clear search"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </span>
              )}
            </div>
          </div>

          <div className="max-h-60 overflow-y-auto">
            {filteredOptions.length === 0 ? (
              <div className="p-4 text-center text-gray-500 text-sm">
                {searchTerm ? 'No options match your search' : 'No options available'}
              </div>
            ) : (
              <div className="p-1">
                {filteredOptions.map((option) => {
                  const isSelected = selectedValues.includes(option.value);
                  return (
                    <div key={option.value}>
                      <label
                        className={`cursor-pointer flex items-center justify-between p-3 rounded-md transition-colors
                        ${option.disabled ? 'opacity-50 cursor-not-allowed' : 'hover:bg-base-100'}
                        ${isSelected ? 'bg-base-200' : 'bg-base-100'}`}
                      >
                        <div className="flex items-center gap-3 flex-1 min-w-0">
                          <input
                            type="checkbox"
                            checked={isSelected}
                            onChange={() => !option.disabled && onToggleOption(option.value)}
                            disabled={disabled || option.disabled}
                            className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                            onClick={(e) => e.stopPropagation()}
                          />
                          <span className="text-sm text-gray-900 truncate flex-1">{option.label}</span>
                        </div>
                        {isSelected && (
                          <TextPill
                            text={option.label}
                            seed={option.value}
                            removable
                            onRemove={() => onToggleOption(option.value)}
                            size="xs"
                            className="m-0 mt-0 mb-0 flex-shrink-0 ml-2"
                          />
                        )}
                      </label>
                    </div>
                  );
                })}
              </div>
            )}
          </div>

          {selectedValues.length > 0 && (
            <div className="p-3 border-t border-gray-200 bg-base-100 text-xs text-gray-600 flex justify-between items-center">
              <span>
                {selectedValues.length} item{selectedValues.length !== 1 ? 's' : ''} selected
              </span>
              <span
                className="text-blue-600 hover:text-blue-800 underline font-medium cursor-pointer"
                onClick={onClearAll}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault();
                    onClearAll();
                  }
                }}
                role="button"
                tabIndex={0}
                aria-label="Clear all selections"
              >
                Clear all
              </span>
            </div>
          )}
        </div>
      )}
    </div>
  );
}