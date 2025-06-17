import React from 'react';
import { TextPill } from '@/components/style/Pills/TextPills';

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

export interface FilterSingleSelectRawProps {
  fieldDefinition: FilterFieldDefinition;
  value: string;
  selectedOption?: Option;
  filteredOptions: Option[];
  isOpen: boolean;
  searchTerm: string;
  disabled?: boolean;
  className?: string;
  allowClear?: boolean;
  containerStyle?: React.CSSProperties;
  onToggleDropdown: () => void;
  onSearchChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onClearSearch: () => void;
  onSelectOption: (value: string) => void;
  onClearSelection: () => void;
}

export function FilterSingleSelectRaw({
  fieldDefinition,
  value,
  selectedOption,
  filteredOptions,
  isOpen,
  searchTerm,
  disabled = false,
  className = '',
  allowClear = true,
  containerStyle,
  onToggleDropdown,
  onSearchChange,
  onClearSearch,
  onSelectOption,
  onClearSelection,
}: FilterSingleSelectRawProps) {
  return (
    <div
      className={`relative z-50 singleselect-container ${className} ${!containerStyle ? 'w-full' : ''}`}
      style={containerStyle}
    >
      <button
        type="button"
        className={
          `w-full min-h-[3rem] p-3 text-left border-2 rounded-lg transition-all duration-200 flex items-center justify-between
          ${isOpen ? 'border-primary-500 ring-2 ring-primary-200' : 'border-gray-300'}
          ${disabled ? 'bg-gray-100 cursor-not-allowed' : 'bg-base-100 hover:border-gray-400 cursor-pointer'}
          ${selectedOption ? 'bg-base-100' : ''}`
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
        <div className="flex-1 min-w-0">
          {!selectedOption ? (
            <span className="text-gray-500">
              {fieldDefinition.placeholder || 'Select an option...'}
            </span>
          ) : (
            <div className="flex items-center justify-between">
              <TextPill
                text={selectedOption.label}
                seed={selectedOption.value}
                removable={allowClear}
                onRemove={onClearSelection}
                size="xs"
                className="m-0 mt-0 mb-0 flex-shrink-0 ml-2"
              />
            </div>
          )}
        </div>
        <svg
          className={`w-5 h-5 text-gray-400 transition-transform flex-shrink-0 ml-2 ${isOpen ? 'rotate-180' : ''}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {isOpen && (
        <div className="absolute top-full left-0 right-0 z-50 mt-1 bg-white border-2 border-gray-200 rounded-lg shadow-lg max-h-80 overflow-hidden">
          <div className="p-3 border-b border-gray-200 bg-gray-50">
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
                <button
                  type="button"
                  onClick={onClearSearch}
                  className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              )}
            </div>
          </div>

          <div className="max-h-60 overflow-y-auto">
            {allowClear && selectedOption && (
              <div className="border-b border-gray-100">
                <button
                  type="button"
                  onClick={onClearSelection}
                  className="w-full text-left p-3 hover:bg-gray-100 transition-colors text-gray-600 italic"
                >
                  <div className="flex items-center gap-3">
                    <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                    <span>Clear selection</span>
                  </div>
                </button>
              </div>
            )}

            {filteredOptions.length === 0 ? (
              <div className="p-4 text-center text-gray-500 text-sm">
                {searchTerm ? 'No options match your search' : 'No options available'}
              </div>
            ) : (
              <div className="p-1">
                {filteredOptions.map((option) => {
                  const isSelected = value === option.value;
                  return (
                    <button
                      key={option.value}
                      type="button"
                      onClick={() => !option.disabled && onSelectOption(option.value)}
                      disabled={option.disabled}
                      className={
                        `w-full text-left p-3 rounded-md transition-colors flex items-center justify-between
                        ${option.disabled ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-100 cursor-pointer'}
                        ${isSelected ? 'bg-blue-50 border-l-4 border-blue-500' : ''}`
                      }
                    >
                      <div className="flex items-center gap-3 flex-1 min-w-0">
                        <div className={
                          `w-4 h-4 rounded-full border-2 flex items-center justify-center flex-shrink-0
                          ${isSelected ? 'border-blue-500 bg-blue-500' : 'border-gray-300'}`
                        }>
                          {isSelected && <div className="w-2 h-2 bg-white rounded-full" />}
                        </div>
                        <span className="text-sm text-gray-900 truncate flex-1">{option.label}</span>
                      </div>
                      {isSelected && (
                        <span className="inline-flex items-center px-2 py-1 bg-blue-100 text-blue-800 rounded-full text-xs flex-shrink-0 ml-2">
                          Selected
                        </span>
                      )}
                    </button>
                  );
                })}
              </div>
            )}
          </div>

          <div className="p-3 border-t border-gray-200 bg-gray-50 text-xs text-gray-600">
            {selectedOption ? (
              <div className="flex justify-between items-center">
                <span>Selected: {selectedOption.label}</span>
                {allowClear && (
                  <button
                    type="button"
                    onClick={onClearSelection}
                    className="text-primary-600 hover:text-primary-800 underline font-medium"
                  >
                    Clear selection
                  </button>
                )}
              </div>
            ) : (
              <span>No selection made</span>
            )}
          </div>
        </div>
      )}
    </div>
  );
}