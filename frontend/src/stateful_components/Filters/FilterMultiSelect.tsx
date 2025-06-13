"use client"
import { TextPill } from '@/style_components/Pills/TextPills';
import React, { useState, useCallback, useMemo, useEffect } from 'react';

// Mock FilterFieldDefinition for demo
interface FilterFieldDefinition {
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

export interface DynamicMultiSelectProps {
  fieldDefinition: FilterFieldDefinition;
  value: string;
  onChange: (value: string) => void;
  onFocus?: () => void;
  onBlur?: () => void;
  disabled?: boolean;
  className?: string;
}

export function DynamicMultiSelect({
  fieldDefinition,
  value,
  onChange,
  onFocus,
  onBlur,
  disabled = false,
  className = "",
}: DynamicMultiSelectProps) {
  const [selectedValues, setSelectedValues] = useState<string[]>(
    value ? value.split(',').filter(Boolean) : []
  );
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  const handleSelectionChange = useCallback((selectedOptions: string[]) => {
    setSelectedValues(selectedOptions);
    onChange(selectedOptions.join(','));
  }, [onChange]);

  const toggleOption = useCallback((optionValue: string) => {
    const newSelection = selectedValues.includes(optionValue)
      ? selectedValues.filter(v => v !== optionValue)
      : [...selectedValues, optionValue];
    handleSelectionChange(newSelection);
  }, [selectedValues, handleSelectionChange]);

  const removeSelectedItem = useCallback((optionValue: string) => {
    const newSelection = selectedValues.filter(v => v !== optionValue);
    handleSelectionChange(newSelection);
  }, [selectedValues, handleSelectionChange]);

  // Filter options based on search term
  const filteredOptions = useMemo(() => {
    if (!searchTerm) return fieldDefinition.options || [];
    return (fieldDefinition.options || []).filter(option =>
      option.label.toLowerCase().includes(searchTerm.toLowerCase()) ||
      option.value.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [fieldDefinition.options, searchTerm]);

  // Get selected option labels for display
  const selectedLabels = useMemo(() => {
    return selectedValues.map(value => {
      const option = fieldDefinition.options?.find(opt => opt.value === value);
      return option ? option.label : value;
    });
  }, [selectedValues, fieldDefinition.options]);

  const handleDropdownToggle = useCallback(() => {
    if (!disabled) {
      setIsOpen(!isOpen);
      if (!isOpen) {
        onFocus?.();
      } else {
        onBlur?.();
      }
    }
  }, [disabled, isOpen, onFocus, onBlur]);

  const handleSearchChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  }, []);

  const clearSearch = useCallback(() => {
    setSearchTerm('');
  }, []);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as Element;
      if (isOpen && !target.closest('.multiselect-container')) {
        setIsOpen(false);
        onBlur?.();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isOpen, onBlur]);

  return (
    <div className={`relative w-full multiselect-container ${className}`}>
      {/* Main Button */}
      <button
        type="button"
        className={`
        w-full min-h-[3rem] p-3 text-left border-2 rounded-lg transition-all duration-200
        ${isOpen ? 'border-blue-500 ring-2 ring-blue-200' : 'border-gray-300'}
        ${disabled ? 'bg-gray-100 cursor-not-allowed' : 'bg-base-100 hover:border-gray-400 cursor-pointer'}
        ${selectedValues.length > 0 ? 'bg-base-100' : ''}
      `}
        onClick={handleDropdownToggle}
        disabled={disabled}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            handleDropdownToggle();
          }
        }}
      >
        <div className="flex items-center justify-between">
          <div className="flex-1 min-w-0">
            {selectedValues.length === 0 ? (
              <span className="text-gray-500">{fieldDefinition.placeholder || 'Select options...'}</span>
            ) : (
              <div className="space-y-2">
                <div className="font-medium text-sm text-gray-700">
                  {fieldDefinition.displayName}: {selectedValues.length} selected
                </div>
                <div className="flex flex-wrap gap-1">
                  {selectedLabels.map((label, index) => (
                    <TextPill
                      key={selectedValues[index]}
                      text={label}
                      seed={selectedValues[index]}
                      removable
                      onRemove={() => removeSelectedItem(selectedValues[index])}
                      size="xs"
                      className="m-0 mt-0 mb-0"
                    />
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Dropdown Arrow */}
          <svg
            className={`w-5 h-5 text-gray-400 transition-transform flex-shrink-0 ml-2 ${isOpen ? 'rotate-180' : ''}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      </button>

      {/* Dropdown Content */}
      {isOpen && (
        <div className="absolute top-full left-0 right-0 z-50 mt-1 bg-base-100 border-2 border-gray-200 rounded-lg shadow-lg max-h-80 overflow-hidden">
          {/* Search Box */}
          <div className="p-3 border-b border-gray-200 bg-base-100">
            <div className="relative">
              <input
                type="text"
                className="w-full pl-9 pr-9 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder={`Search ${fieldDefinition.displayName.toLowerCase()}...`}
                value={searchTerm}
                onChange={handleSearchChange}
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
                  onClick={clearSearch}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter' || e.key === ' ') {
                      e.preventDefault();
                      clearSearch();
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

          {/* Options List */}
          <div className="max-h-60 overflow-y-auto">
            {filteredOptions.length === 0 ? (
              <div className="p-4 text-center text-gray-500 text-sm">
                {searchTerm ? 'No options match your search' : 'No options available'}
              </div>
            ) : (
              <div className="p-1">
                {filteredOptions.map((option) => {
                  const isSelected = selectedValues.includes(option.value);
                  const optionLabel = fieldDefinition.options?.find(opt => opt.value === option.value)?.label || option.value;

                  return (
                    <div key={option.value}>
                      <label
                        className={`
                        cursor-pointer flex items-center justify-between p-3 rounded-md transition-colors
                        ${option.disabled ? 'opacity-50 cursor-not-allowed' : 'hover:bg-base-100'}
                        ${isSelected ? 'bg-base-200' : 'bg-base-100'}
                      `}
                      >
                        <div className="flex items-center gap-3 flex-1 min-w-0">
                          <input
                            type="checkbox"
                            checked={isSelected}
                            onChange={() => !option.disabled && toggleOption(option.value)}
                            disabled={disabled || option.disabled}
                            className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                            onClick={(e) => e.stopPropagation()}
                          />
                          <span className="text-sm text-gray-900 truncate flex-1">{option.label}</span>
                        </div>

                        {/* Option Pill */}
                        {isSelected && (
                          <TextPill
                            text={optionLabel}
                            seed={option.value}
                            removable
                            onRemove={() => toggleOption(option.value)}
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

          {/* Footer with selection summary */}
          {selectedValues.length > 0 && (
            <div className="p-3 border-t border-gray-200 bg-base-100 text-xs text-gray-600 flex justify-between items-center">
              <span>
                {selectedValues.length} item{selectedValues.length !== 1 ? 's' : ''} selected
              </span>
              <span
                className="text-blue-600 hover:text-blue-800 underline font-medium cursor-pointer"
                onClick={() => handleSelectionChange([])}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault();
                    handleSelectionChange([]);
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
